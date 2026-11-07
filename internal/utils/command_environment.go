package utils

import (
	"regexp"
	"strings"
	"unicode"
)

type CommandEnvironment struct {
	Variables  map[string]string `json:"variables"`
	Keys       []string          `json:"keys"`
	Command    string            `json:"command"`
	Migratable bool              `json:"migratable"`
}

var (
	prefixRegex = regexp.MustCompile(`^\s*(env|cross-env|export)\s+`)
	keyRegex    = regexp.MustCompile(`^\s*([A-Za-z_][A-Za-z0-9_]*)=`)
)

func inspectLeadingEnvironment(input string) CommandEnvironment {
	variables := make(map[string]string)
	var keys []string
	safe := true
	offset := 0

	prefixMatch := prefixRegex.FindStringSubmatch(input)
	var prefixKind string
	if len(prefixMatch) > 0 {
		offset = len(prefixMatch[0])
		prefixKind = prefixMatch[1]
	}

	for offset < len(input) {
		sub := input[offset:]
		match := keyRegex.FindStringSubmatchIndex(sub)
		if len(match) == 0 {
			break
		}
		key := sub[match[2]:match[3]]
		if key == "__proto__" {
			safe = false
		}
		offset += match[1]

		var valBuilder strings.Builder
		quote := rune(0)

		for offset < len(input) {
			r := rune(input[offset])
			if quote == 0 && (unicode.IsSpace(r) || strings.ContainsRune(";&|<>", r)) {
				break
			}
			offset++
			if r == '\'' && quote != '"' {
				if quote == '\'' {
					quote = 0
				} else {
					quote = '\''
				}
				continue
			}
			if r == '"' && quote != '\'' {
				if quote == '"' {
					quote = 0
				} else {
					quote = '"'
				}
				continue
			}
			if r == '\\' && quote != '\'' {
				if offset < len(input) {
					next := rune(input[offset])
					offset++
					if quote == '"' && !strings.ContainsRune("$`\"\\", next) {
						valBuilder.WriteRune('\\')
					}
					valBuilder.WriteRune(next)
					continue
				} else {
					safe = false
					break
				}
			}
			if quote != '\'' && (r == '$' || r == '`') {
				safe = false
			}
			if quote == 0 && strings.ContainsRune("~*?[]{}()", r) {
				safe = false
			}
			valBuilder.WriteRune(r)
		}

		if quote != 0 {
			safe = false
		}
		if _, exists := variables[key]; exists {
			safe = false
		}
		variables[key] = valBuilder.String()
		keys = append(keys, key)
	}

	command := strings.TrimLeftFunc(input[offset:], unicode.IsSpace)
	if prefixKind == "export" {
		if strings.HasPrefix(command, "&&") {
			command = strings.TrimLeftFunc(command[2:], unicode.IsSpace)
		} else {
			safe = false
		}
	}
	if command == "" || (len(command) > 0 && strings.ContainsAny(string(command[0]), ";&|<>")) {
		safe = false
	}

	finalCommand := input
	if len(keys) > 0 {
		finalCommand = command
	}

	return CommandEnvironment{
		Variables:  variables,
		Keys:       keys,
		Command:    finalCommand,
		Migratable: len(keys) > 0 && safe,
	}
}

func InspectCommandEnvironment(input string) CommandEnvironment {
	leading := inspectLeadingEnvironment(input)
	keySet := make(map[string]bool)
	for _, k := range leading.Keys {
		keySet[k] = true
	}

	quote := rune(0)
	chained := false

	for i := 0; i < len(input); i++ {
		r := rune(input[i])
		if r == '\\' && quote != '\'' {
			i++
			continue
		}
		if r == '\'' && quote != '"' {
			if quote == '\'' {
				quote = 0
			} else {
				quote = '\''
			}
			continue
		}
		if r == '"' && quote != '\'' {
			if quote == '"' {
				quote = 0
			} else {
				quote = '"'
			}
			continue
		}
		if quote != 0 || !strings.ContainsRune(";&|(\n", r) {
			continue
		}
		if r == '\n' && strings.TrimSpace(input[:i]) == "" {
			continue
		}

		next := inspectLeadingEnvironment(input[i+1:])
		if len(next.Keys) > 0 {
			chained = true
		}
		for _, k := range next.Keys {
			keySet[k] = true
		}
	}

	allKeys := make([]string, 0, len(keySet))
	for k := range keySet {
		allKeys = append(allKeys, k)
	}

	return CommandEnvironment{
		Variables:  leading.Variables,
		Keys:       allKeys,
		Command:    leading.Command,
		Migratable: leading.Migratable && !chained,
	}
}

func CommandEnvironmentWarning(command, label string) string {
	res := InspectCommandEnvironment(command)
	if len(res.Keys) == 0 {
		return ""
	}
	msg := label + " embeds environment variables: " + strings.Join(res.Keys, ", ") + ". Manage these in Environment; saved variables are available during builds and at runtime."
	if !res.Migratable {
		msg += " This command needs a manual edit before migration."
	}
	return msg
}
