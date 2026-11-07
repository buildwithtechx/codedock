package repositories

import "strings"

func EscapeLikeTerm(s string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(s)
}

func ContainsLikePattern(term string) string {
	return "%" + EscapeLikeTerm(strings.TrimSpace(term)) + "%"
}
