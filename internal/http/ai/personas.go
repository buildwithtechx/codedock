package ai

import "strings"

func IsReadOnlyPersona(persona string) bool {
	switch strings.ToLower(strings.TrimSpace(persona)) {
	case "logs", "security":
		return true
	default:
		return false
	}
}

func PersonaDisplayName(persona string) string {
	switch strings.ToLower(strings.TrimSpace(persona)) {
	case "logs":
		return "Logs Investigator"
	case "security":
		return "Security Reviewer"
	case "devops":
		return "DevOps Co-pilot"
	default:
		return "General Assistant"
	}
}
