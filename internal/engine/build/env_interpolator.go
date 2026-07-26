package build

import (
	"regexp"
	"strings"
)

var interpolationPattern = regexp.MustCompile(`\$\{([a-zA-Z0-9_-]+)\.([a-zA-Z0-9_]+)\}`)

func InterpolateEnvVars(envMap map[string]string, registry map[string]map[string]string) map[string]string {
	if len(registry) == 0 {
		return envMap
	}
	result := make(map[string]string, len(envMap))
	for k, v := range envMap {
		result[k] = resolveValue(v, registry)
	}
	return result
}

func resolveValue(value string, registry map[string]map[string]string) string {
	return interpolationPattern.ReplaceAllStringFunc(value, func(match string) string {
		parts := interpolationPattern.FindStringSubmatch(match)
		if len(parts) != 3 {
			return match
		}
		serviceName := parts[1]
		varKey := parts[2]

		if serviceVars, ok := registry[serviceName]; ok {
			if resolved, ok := serviceVars[varKey]; ok {
				return resolved
			}
		}

		normalizedName := strings.ToLower(strings.ReplaceAll(serviceName, "_", "-"))
		for name, vars := range registry {
			if strings.ToLower(strings.ReplaceAll(name, "_", "-")) == normalizedName {
				if resolved, ok := vars[varKey]; ok {
					return resolved
				}
			}
		}

		return match
	})
}
