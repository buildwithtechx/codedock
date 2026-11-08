package docker

import (
	"fmt"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

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
		if source.Healthcheck != nil {
			entry["healthcheck"] = renderInstallHealthcheck(source.Healthcheck)
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

func renderInstallHealthcheck(check *CatalogHealthcheck) map[string]any {
	rendered := map[string]any{"test": check.Test}
	if check.Interval != "" {
		rendered["interval"] = check.Interval
	}
	if check.Timeout != "" {
		rendered["timeout"] = check.Timeout
	}
	if check.Retries > 0 {
		rendered["retries"] = check.Retries
	}
	if check.StartPeriod != "" {
		rendered["start_period"] = check.StartPeriod
	}
	return rendered
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
