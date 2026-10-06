package build

import (
	"fmt"
	"strings"
)

func InterpolateEnvVarsStrict(values map[string]string, registry map[string]map[string]string) (map[string]string, error) {
	result := make(map[string]string, len(values))
	var resolve func(string, map[string]bool, int) (string, error)
	resolve = func(value string, active map[string]bool, depth int) (string, error) {
		if depth > 50 {
			return "", fmt.Errorf("variable binding depth exceeded")
		}
		var resolutionErr error
		resolved := interpolationPattern.ReplaceAllStringFunc(value, func(match string) string {
			parts := interpolationPattern.FindStringSubmatch(match)
			name, key := parts[1], parts[2]
			if active[name+"."+key] {
				resolutionErr = fmt.Errorf("variable binding cycle at %s.%s", name, key)
				return ""
			}
			source := registry[name]
			if source == nil {
				for candidate, variables := range registry {
					if strings.EqualFold(strings.ReplaceAll(candidate, "_", "-"), strings.ReplaceAll(name, "_", "-")) {
						source = variables
						break
					}
				}
			}
			value, ok := source[key]
			if !ok {
				resolutionErr = fmt.Errorf("variable source %s.%s is unavailable in this environment", name, key)
				return ""
			}
			active[name+"."+key] = true
			expanded, err := resolve(value, active, depth+1)
			delete(active, name+"."+key)
			if err != nil {
				resolutionErr = err
			}
			return expanded
		})
		return resolved, resolutionErr
	}
	for key, value := range values {
		resolved, err := resolve(value, map[string]bool{}, 0)
		if err != nil {
			return nil, err
		}
		result[key] = resolved
	}
	return result, nil
}
