package docker

import (
	"fmt"
	"strconv"
	"strings"
)

func catalogMetadata(tmpl ComposeTemplate) *CodedockMetadata {
	if tmpl.XCodedock != nil {
		return tmpl.XCodedock
	}
	for _, service := range tmpl.Services {
		if service.XCodedock != nil {
			return service.XCodedock
		}
	}
	return nil
}

func validateCatalogDependencies(services map[string]ComposeService) []string {
	failures := []string{}
	active, visited := map[string]bool{}, map[string]bool{}
	var visit func(string) error
	visit = func(name string) error {
		if _, ok := services[name]; !ok {
			return fmt.Errorf("unknown dependency: %s", name)
		}
		if active[name] {
			return fmt.Errorf("dependency cycle at %s", name)
		}
		if visited[name] {
			return nil
		}
		active[name] = true
		for _, dependency := range services[name].DependsOn {
			if err := visit(dependency); err != nil {
				return err
			}
		}
		active[name], visited[name] = false, true
		return nil
	}
	for name := range services {
		if err := visit(name); err != nil {
			failures = append(failures, err.Error())
		}
	}
	return failures
}

func ParseCatalogPort(entry string) (CatalogPortMapping, error) {
	mapping := CatalogPortMapping{Protocol: "tcp"}
	value, protocol, _ := strings.Cut(entry, "/")
	if protocol != "" {
		if protocol != "tcp" && protocol != "udp" && protocol != "sctp" {
			return mapping, fmt.Errorf("invalid port protocol in %q", entry)
		}
		mapping.Protocol = protocol
	}
	parts := strings.Split(value, ":")
	if len(parts) > 2 {
		return mapping, fmt.Errorf("invalid port mapping %q", entry)
	}
	container, err := parseCatalogPortNumber(parts[len(parts)-1])
	if err != nil {
		return mapping, fmt.Errorf("invalid port mapping %q", entry)
	}
	mapping.Container = container
	if len(parts) == 2 && parts[0] != "" {
		host, err := parseCatalogPortNumber(parts[0])
		if err != nil {
			return mapping, fmt.Errorf("invalid port mapping %q", entry)
		}
		mapping.Host = host
	}
	return mapping, nil
}

func parseCatalogPortNumber(value string) (int, error) {
	number, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || number < 1 || number > 65535 {
		return 0, fmt.Errorf("port out of range")
	}
	return number, nil
}

func ParseCatalogVolume(entry string) (CatalogVolumeMount, error) {
	mount := CatalogVolumeMount{}
	parts := strings.Split(entry, ":")
	if len(parts) < 2 || len(parts) > 3 {
		return mount, fmt.Errorf("invalid volume mount %q", entry)
	}
	source := parts[0]
	if source == "" || strings.ContainsAny(source, "/\\.") || strings.HasPrefix(source, "~") || strings.Contains(source, "$") {
		return mount, fmt.Errorf("volume %q must be a declared named volume", entry)
	}
	if !strings.HasPrefix(parts[1], "/") {
		return mount, fmt.Errorf("volume %q needs an absolute container path", entry)
	}
	mount.Source, mount.Target = source, parts[1]
	return mount, nil
}
