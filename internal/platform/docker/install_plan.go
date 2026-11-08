package docker

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"sort"
	"strings"
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
	inputs, err := resolveInstallInputs(catalogInputSpecs(tmpl), req.Environment)
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
		env, err := resolveInstallEnv(serviceName, service, secrets, inputs, req.Environment, plan.secretKeys)
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
	if meta := catalogMetadata(tmpl); meta != nil {
		return meta.Secrets
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

func catalogInputSpecs(tmpl ComposeTemplate) []CodedockInputSpec {
	if meta := catalogMetadata(tmpl); meta != nil {
		return meta.Inputs
	}
	return nil
}

func resolveInstallInputs(specs []CodedockInputSpec, provided map[string]string) (map[string]string, error) {
	inputs := map[string]string{}
	for _, spec := range specs {
		if value := strings.TrimSpace(provided[spec.Var]); value != "" {
			inputs[spec.Var] = value
			continue
		}
		if spec.Default != "" {
			inputs[spec.Var] = spec.Default
			continue
		}
		if spec.Required {
			return nil, fmt.Errorf("input %s is required", spec.Var)
		}
		inputs[spec.Var] = ""
	}
	return inputs, nil
}

func resolveInstallEnv(serviceName string, service ComposeService, secrets, inputs, overrides map[string]string, secretKeys map[string]bool) ([]string, error) {
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
		resolved = catalogInputRef.ReplaceAllStringFunc(resolved, func(match string) string {
			return inputs[catalogInputRef.FindStringSubmatch(match)[1]]
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
