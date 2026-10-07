package compose

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type CatalogPortMapping struct {
	Host      int
	Container int
	Protocol  string
}

type CatalogVolumeMount struct {
	Source string
	Target string
}

var catalogEnvName = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
var catalogSecretRef = regexp.MustCompile(`\$\{secret:([A-Za-z_][A-Za-z0-9_]*)\}`)

func ValidateTemplate(id string, tmpl ComposeTemplate) error {
	failures := []string{}
	if len(tmpl.Services) == 0 {
		failures = append(failures, "services are required")
	}
	for name, service := range tmpl.Services {
		failures = append(failures, validateCatalogService(name, service, tmpl.Volumes)...)
	}
	failures = append(failures, validateCatalogMetadata(tmpl)...)
	failures = append(failures, validateCatalogSecrets(tmpl)...)
	failures = append(failures, validateCatalogDependencies(tmpl.Services)...)
	if len(failures) == 0 {
		return nil
	}
	return fmt.Errorf("catalogue template %s is invalid: %s", id, strings.Join(failures, "; "))
}

func validateCatalogService(name string, service ComposeService, volumes map[string]any) []string {
	failures := []string{}
	if strings.TrimSpace(service.Image) == "" {
		failures = append(failures, fmt.Sprintf("service %s needs an image", name))
	}
	for _, entry := range service.Ports {
		if _, err := ParseCatalogPort(entry); err != nil {
			failures = append(failures, fmt.Sprintf("service %s: %v", name, err))
		}
	}
	for _, entry := range service.Volumes {
		mount, err := ParseCatalogVolume(entry)
		if err != nil {
			failures = append(failures, fmt.Sprintf("service %s: %v", name, err))
			continue
		}
		if _, declared := volumes[mount.Source]; !declared {
			failures = append(failures, fmt.Sprintf("service %s: volume %s is not declared", name, mount.Source))
		}
	}
	for _, entry := range service.Environment {
		if key, _, _ := strings.Cut(entry, "="); !catalogEnvName.MatchString(key) {
			failures = append(failures, fmt.Sprintf("service %s: invalid environment entry %q", name, entry))
		}
	}
	return failures
}

func validateCatalogMetadata(tmpl ComposeTemplate) []string {
	meta := tmpl.XCodedock
	if meta == nil {
		for _, service := range tmpl.Services {
			if service.XCodedock != nil {
				meta = service.XCodedock
				break
			}
		}
	}
	if meta == nil {
		return []string{"x-codedock metadata is required"}
	}
	if meta.IsOneClick {
		failures := []string{}
		if strings.TrimSpace(meta.Name) == "" {
			failures = append(failures, "one-click name is required")
		}
		if strings.TrimSpace(meta.Description) == "" {
			failures = append(failures, "one-click description is required")
		}
		return failures
	}
	if !meta.IsDatabase {
		return []string{"template must be marked as one-click or database"}
	}
	return nil
}

func validateCatalogSecrets(tmpl ComposeTemplate) []string {
	failures := []string{}
	declared := map[string]bool{}
	meta := tmpl.XCodedock
	if meta == nil {
		for _, service := range tmpl.Services {
			if service.XCodedock != nil {
				meta = service.XCodedock
				break
			}
		}
	}
	if meta != nil {
		for _, spec := range meta.Secrets {
			if !catalogEnvName.MatchString(spec.Var) {
				failures = append(failures, fmt.Sprintf("invalid secret variable %q", spec.Var))
				continue
			}
			if spec.Length < 0 || spec.Length > 128 {
				failures = append(failures, fmt.Sprintf("secret %s length must be between 0 and 128", spec.Var))
			}
			declared[spec.Var] = true
		}
	}
	for name, service := range tmpl.Services {
		for _, entry := range service.Environment {
			for _, ref := range catalogSecretRef.FindAllStringSubmatch(entry, -1) {
				if !declared[ref[1]] {
					failures = append(failures, fmt.Sprintf("service %s: undeclared secret %s", name, ref[1]))
				}
			}
		}
	}
	if meta != nil {
		used := map[string]bool{}
		for _, service := range tmpl.Services {
			for _, entry := range service.Environment {
				for _, ref := range catalogSecretRef.FindAllStringSubmatch(entry, -1) {
					used[ref[1]] = true
				}
			}
		}
		for _, spec := range meta.Secrets {
			if catalogEnvName.MatchString(spec.Var) && !used[spec.Var] {
				failures = append(failures, fmt.Sprintf("secret %s is never referenced", spec.Var))
			}
		}
	}
	return failures
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
