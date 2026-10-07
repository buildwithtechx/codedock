package compose

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type InstallRequest struct {
	Name         string
	Secrets      map[string]string
	Environment  map[string]string
	HostPort     int
	Domain       string
	VolumePrefix string
}

type InstallServicePlan struct {
	Service string   `yaml:"-"`
	Image   string   `yaml:"image"`
	Env     []string `yaml:"environment,omitempty"`
	Ports   []string `yaml:"ports,omitempty"`
}

type InstallVolumePlan struct {
	Service string
	Name    string
	Target  string
}

type InstallPlan struct {
	AppID            string
	Name             string
	Services         []InstallServicePlan
	Volumes          []InstallVolumePlan
	GeneratedSecrets []string
	ComposeYAML      string
	Digest           string
	secretValues     map[string]string
	secretKeys       map[string]bool
}

type SecretGenerator func(length int) (string, error)

var installNameChars = regexp.MustCompile(`[^a-z0-9-]`)
var installDomainChars = regexp.MustCompile(`^[a-z0-9]([a-z0-9.-]*[a-z0-9])?$`)

func DefaultSecretGenerator(length int) (string, error) {
	if length <= 0 {
		length = 24
	}
	if length > 128 {
		length = 128
	}
	raw := make([]byte, length)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("generate secret: %w", err)
	}
	return hex.EncodeToString(raw), nil
}

func ResolveInstallPlan(appID string, tmpl ComposeTemplate, req InstallRequest, generate SecretGenerator) (*InstallPlan, error) {
	if generate == nil {
		generate = DefaultSecretGenerator
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, fmt.Errorf("install name is required")
	}
	specs := catalogSecretSpecs(tmpl)
	secrets, generated, err := resolveInstallSecrets(specs, req.Secrets, generate)
	if err != nil {
		return nil, err
	}
	primary := catalogPrimaryService(tmpl)
	plan := &InstallPlan{
		AppID:            appID,
		Name:             name,
		GeneratedSecrets: generated,
		secretValues:     secrets,
		secretKeys:       map[string]bool{},
	}
	seenHosts := map[string]string{}
	for _, serviceName := range catalogServiceNames(tmpl) {
		service := tmpl.Services[serviceName]
		env, err := resolveInstallEnv(serviceName, service, secrets, req.Environment, plan.secretKeys)
		if err != nil {
			return nil, err
		}
		ports, err := resolveInstallPorts(serviceName, service, serviceName == primary, req.HostPort, seenHosts)
		if err != nil {
			return nil, err
		}
		plan.Services = append(plan.Services, InstallServicePlan{Service: serviceName, Image: service.Image, Env: env, Ports: ports})
		plan.Volumes = append(plan.Volumes, resolveInstallVolumes(serviceName, service, installVolumePrefix(req.VolumePrefix, name))...)
	}
	document, err := renderInstallDocument(tmpl, plan, req.Domain)
	if err != nil {
		return nil, err
	}
	plan.ComposeYAML = document
	masked := document
	for _, secret := range secrets {
		if secret != "" {
			masked = strings.ReplaceAll(masked, secret, "***")
		}
	}
	digest := sha256.Sum256([]byte(appID + "|" + name + "|" + masked))
	plan.Digest = hex.EncodeToString(digest[:])
	return plan, nil
}

func (p *InstallPlan) Masked() *InstallPlan {
	masked := *p
	masked.Services = make([]InstallServicePlan, len(p.Services))
	for i, service := range p.Services {
		entry := service
		entry.Env = make([]string, len(service.Env))
		for j, value := range service.Env {
			key, _, _ := strings.Cut(value, "=")
			if p.secretKeys[key] {
				entry.Env[j] = key + "=***"
			} else {
				entry.Env[j] = value
			}
		}
		masked.Services[i] = entry
	}
	masked.ComposeYAML = p.ComposeYAML
	for _, secret := range p.secretValues {
		if secret != "" {
			masked.ComposeYAML = strings.ReplaceAll(masked.ComposeYAML, secret, "***")
		}
	}
	return &masked
}

func catalogSecretSpecs(tmpl ComposeTemplate) []CodedockSecretSpec {
	if tmpl.XCodedock != nil {
		return tmpl.XCodedock.Secrets
	}
	for _, service := range tmpl.Services {
		if service.XCodedock != nil {
			return service.XCodedock.Secrets
		}
	}
	return nil
}

func resolveInstallSecrets(specs []CodedockSecretSpec, provided map[string]string, generate SecretGenerator) (map[string]string, []string, error) {
	secrets := map[string]string{}
	generated := []string{}
	for _, spec := range specs {
		if value := strings.TrimSpace(provided[spec.Var]); value != "" {
			if len(value) < 8 {
				return nil, nil, fmt.Errorf("secret %s must be at least 8 characters", spec.Var)
			}
			secrets[spec.Var] = value
			continue
		}
		value, err := generate(spec.Length)
		if err != nil {
			return nil, nil, err
		}
		secrets[spec.Var] = value
		generated = append(generated, spec.Var)
	}
	return secrets, generated, nil
}

func resolveInstallEnv(serviceName string, service ComposeService, secrets map[string]string, overrides map[string]string, secretKeys map[string]bool) ([]string, error) {
	env := make([]string, 0, len(service.Environment)+len(overrides))
	seen := map[string]bool{}
	for _, entry := range service.Environment {
		key, value, _ := strings.Cut(entry, "=")
		if strings.Contains(value, "${db.") {
			return nil, fmt.Errorf("service %s: database placeholders require the database installer", serviceName)
		}
		resolved := catalogSecretRef.ReplaceAllStringFunc(value, func(match string) string {
			name := catalogSecretRef.FindStringSubmatch(match)[1]
			secretKeys[key] = true
			return secrets[name]
		})
		if override, ok := overrides[key]; ok {
			resolved = override
		}
		env = append(env, key+"="+resolved)
		seen[key] = true
	}
	extra := []string{}
	for key := range overrides {
		if !seen[key] {
			extra = append(extra, key)
		}
	}
	sort.Strings(extra)
	for _, key := range extra {
		if !catalogEnvName.MatchString(key) {
			return nil, fmt.Errorf("service %s: invalid environment variable %q", serviceName, key)
		}
		env = append(env, key+"="+overrides[key])
	}
	return env, nil
}

func resolveInstallPorts(serviceName string, service ComposeService, primary bool, hostPort int, seen map[string]string) ([]string, error) {
	ports := make([]string, 0, len(service.Ports))
	for index, entry := range service.Ports {
		mapping, err := ParseCatalogPort(entry)
		if err != nil {
			return nil, fmt.Errorf("service %s: %w", serviceName, err)
		}
		if primary && index == 0 && hostPort > 0 {
			if hostPort < 1 || hostPort > 65535 {
				return nil, fmt.Errorf("host port must be between 1 and 65535")
			}
			mapping.Host = hostPort
		}
		rendered := fmt.Sprintf("%d", mapping.Container)
		if mapping.Host > 0 {
			rendered = fmt.Sprintf("%d:%d", mapping.Host, mapping.Container)
		}
		if mapping.Protocol != "tcp" {
			rendered += "/" + mapping.Protocol
		}
		if mapping.Host > 0 {
			key := fmt.Sprintf("%d/%s", mapping.Host, mapping.Protocol)
			if previous := seen[key]; previous != "" {
				return nil, fmt.Errorf("host port %s is used by %s and %s", key, previous, serviceName)
			}
			seen[key] = serviceName
		}
		ports = append(ports, rendered)
	}
	return ports, nil
}

func resolveInstallVolumes(serviceName string, service ComposeService, prefix string) []InstallVolumePlan {
	plans := []InstallVolumePlan{}
	for _, entry := range service.Volumes {
		mount, err := ParseCatalogVolume(entry)
		if err != nil {
			continue
		}
		plans = append(plans, InstallVolumePlan{Service: serviceName, Name: prefix + mount.Source, Target: mount.Target})
	}
	return plans
}

func installVolumePrefix(prefix, name string) string {
	normalized := strings.ToLower(strings.TrimSpace(prefix))
	if normalized == "" {
		normalized = "app-" + installNameChars.ReplaceAllString(strings.ToLower(name), "-")
	}
	normalized = installNameChars.ReplaceAllString(normalized, "-")
	normalized = strings.Trim(normalized, "-")
	if normalized == "" {
		normalized = "app"
	}
	return normalized + "-"
}

func catalogServiceNames(tmpl ComposeTemplate) []string {
	names := make([]string, 0, len(tmpl.Services))
	for name := range tmpl.Services {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func catalogPrimaryService(tmpl ComposeTemplate) string {
	if tmpl.XCodedock != nil && tmpl.XCodedock.IsOneClick {
		for _, name := range catalogServiceNames(tmpl) {
			if len(tmpl.Services[name].Ports) > 0 {
				return name
			}
		}
	}
	for _, name := range catalogServiceNames(tmpl) {
		if meta := tmpl.Services[name].XCodedock; meta != nil && meta.IsOneClick {
			return name
		}
	}
	names := catalogServiceNames(tmpl)
	if len(names) == 0 {
		return ""
	}
	return names[0]
}

func renderInstallDocument(tmpl ComposeTemplate, plan *InstallPlan, domain string) (string, error) {
	trimmedDomain := strings.TrimSpace(domain)
	if trimmedDomain != "" && !installDomainChars.MatchString(strings.ToLower(trimmedDomain)) {
		return "", fmt.Errorf("invalid domain %q", domain)
	}
	services := map[string]any{}
	for _, resolved := range plan.Services {
		source := tmpl.Services[resolved.Service]
		entry := map[string]any{"image": resolved.Image}
		if len(resolved.Env) > 0 {
			entry["environment"] = resolved.Env
		}
		if len(resolved.Ports) > 0 {
			entry["ports"] = resolved.Ports
		}
		volumes := []string{}
		for _, mount := range plan.Volumes {
			if mount.Service == resolved.Service {
				volumes = append(volumes, mount.Name+":"+mount.Target)
			}
		}
		if len(volumes) > 0 {
			entry["volumes"] = volumes
		}
		if len(source.Command) > 0 {
			entry["command"] = source.Command
		}
		if len(source.DependsOn) > 0 {
			entry["depends_on"] = source.DependsOn
		}
		if trimmedDomain != "" && resolved.Service == catalogPrimaryService(tmpl) {
			entry["labels"] = installDomainLabels(plan.Name, trimmedDomain, resolved.Ports)
		}
		services[resolved.Service] = entry
	}
	document := map[string]any{"services": services}
	if len(plan.Volumes) > 0 {
		volumes := map[string]any{}
		for _, mount := range plan.Volumes {
			volumes[mount.Name] = nil
		}
		document["volumes"] = volumes
	}
	rendered, err := yaml.Marshal(document)
	if err != nil {
		return "", fmt.Errorf("render install document: %w", err)
	}
	return string(rendered), nil
}

func installDomainLabels(name, domain string, ports []string) []string {
	router := installNameChars.ReplaceAllString(strings.ToLower(name), "-")
	containerPort := ""
	for _, entry := range ports {
		if mapping, err := ParseCatalogPort(entry); err == nil {
			containerPort = fmt.Sprintf("%d", mapping.Container)
			break
		}
	}
	labels := []string{
		"traefik.enable=true",
		fmt.Sprintf("traefik.http.routers.%s.rule=Host(`%s`)", router, domain),
		fmt.Sprintf("traefik.http.routers.%s.entrypoints=websecure", router),
		fmt.Sprintf("traefik.http.routers.%s.tls.certresolver=letsencrypt", router),
	}
	if containerPort != "" {
		labels = append(labels, fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port=%s", router, containerPort))
	}
	return labels
}
