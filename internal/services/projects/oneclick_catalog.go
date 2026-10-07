package projects

import (
	"sort"
	"strings"

	"codedock.run/codedock/internal/engine/compose"
	"codedock.run/codedock/internal/models"
)

func extractOneClickApp(id string, tmpl *compose.ComposeTemplate) *models.OneClickApp {
	meta := findOneClickMetadata(tmpl)
	if meta == nil {
		return nil
	}
	return &models.OneClickApp{
		ID:           id,
		Name:         meta.Name,
		Description:  meta.Description,
		Icon:         meta.Icon,
		Category:     catalogCategory(meta.Category),
		DockerImage:  catalogPrimaryImage(tmpl),
		DefaultPort:  extractPort(tmpl),
		Services:     catalogServiceNames(tmpl),
		Volumes:      catalogVolumeNames(tmpl),
		EnvVariables: catalogEnvVariables(tmpl, meta),
	}
}

func findOneClickMetadata(tmpl *compose.ComposeTemplate) *compose.CodedockMetadata {
	if tmpl.XCodedock != nil && tmpl.XCodedock.IsOneClick {
		return tmpl.XCodedock
	}
	for _, name := range catalogServiceNames(tmpl) {
		if meta := tmpl.Services[name].XCodedock; meta != nil && meta.IsOneClick {
			return meta
		}
	}
	return nil
}

func catalogCategory(value string) string {
	if strings.TrimSpace(value) == "" {
		return "apps"
	}
	return value
}

func catalogPrimaryImage(tmpl *compose.ComposeTemplate) string {
	for _, name := range catalogServiceNames(tmpl) {
		if len(tmpl.Services[name].Ports) > 0 {
			return tmpl.Services[name].Image
		}
	}
	names := catalogServiceNames(tmpl)
	if len(names) == 0 {
		return ""
	}
	return tmpl.Services[names[0]].Image
}

func catalogServiceNames(tmpl *compose.ComposeTemplate) []string {
	names := make([]string, 0, len(tmpl.Services))
	for name := range tmpl.Services {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func catalogVolumeNames(tmpl *compose.ComposeTemplate) []string {
	seen := map[string]bool{}
	names := []string{}
	for _, serviceName := range catalogServiceNames(tmpl) {
		for _, entry := range tmpl.Services[serviceName].Volumes {
			mount, err := compose.ParseCatalogVolume(entry)
			if err != nil || seen[mount.Source] {
				continue
			}
			seen[mount.Source] = true
			names = append(names, mount.Source)
		}
	}
	return names
}

func catalogEnvVariables(tmpl *compose.ComposeTemplate, meta *compose.CodedockMetadata) []models.OneClickEnvVar {
	secrets := map[string]compose.CodedockSecretSpec{}
	for _, spec := range meta.Secrets {
		secrets[spec.Var] = spec
	}
	variables := []models.OneClickEnvVar{}
	seen := map[string]bool{}
	for _, serviceName := range catalogServiceNames(tmpl) {
		for _, entry := range tmpl.Services[serviceName].Environment {
			key, value, _ := strings.Cut(entry, "=")
			if seen[key] {
				continue
			}
			seen[key] = true
			variable := models.OneClickEnvVar{Key: key, Label: key}
			if ref, ok := catalogSecretReference(value); ok {
				variable.Secret = true
				if spec, declared := secrets[ref]; declared && spec.Label != "" {
					variable.Label = spec.Label
				}
			} else if !strings.Contains(value, "${") {
				variable.DefaultValue = value
			}
			variables = append(variables, variable)
		}
	}
	return variables
}

func catalogSecretReference(value string) (string, bool) {
	start := strings.Index(value, "${secret:")
	if start < 0 {
		return "", false
	}
	rest := value[start+len("${secret:"):]
	end := strings.Index(rest, "}")
	if end <= 0 {
		return "", false
	}
	return rest[:end], true
}

func extractPort(tmpl *compose.ComposeTemplate) int {
	for _, name := range catalogServiceNames(tmpl) {
		service := tmpl.Services[name]
		if len(service.Ports) == 0 {
			continue
		}
		if mapping, err := compose.ParseCatalogPort(service.Ports[0]); err == nil {
			if mapping.Host > 0 {
				return mapping.Host
			}
			return mapping.Container
		}
	}
	return 3000
}
