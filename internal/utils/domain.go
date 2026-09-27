package utils

import (
	"fmt"
	"hash/fnv"
	"strings"

	"codedock.run/codedock/internal/config"
)

func GenerateAppDomain(projectNameOrID string, hostIP string, wildcardDomain string) string {
	if wildcardDomain == "" {
		wildcardDomain = config.Get().Domains.WildcardDomain
	}

	if wildcardDomain != "" {
		cleanName := SanitizeDomainName(projectNameOrID)
		base := strings.TrimPrefix(wildcardDomain, "*.")
		if !strings.HasPrefix(base, "http") {
			return fmt.Sprintf("https://%s.%s", cleanName, base)
		}
		return fmt.Sprintf("%s://%s.%s", parseScheme(base), cleanName, parseHost(base))
	}

	if hostIP == "" {
		hostIP = config.Get().Server.HostIP
	}

	if hostIP == "" {
		hostIP = "127.0.0.1"
	}

	cleanIP := strings.ReplaceAll(strings.TrimSpace(hostIP), ".", "-")
	cleanName := SanitizeDomainName(projectNameOrID)
	magicDomain := config.Get().Domains.MagicDomain
	if magicDomain == "" {
		magicDomain = "sslip.io"
	}

	return fmt.Sprintf("http://%s.%s.%s", cleanName, cleanIP, magicDomain)
}

func parseScheme(url string) string {
	if strings.HasPrefix(url, "https://") {
		return "https"
	}
	return "http"
}

func parseHost(url string) string {
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	return url
}

func SanitizeDomainName(name string) string {
	raw := strings.ToLower(strings.TrimSpace(name))
	var builder strings.Builder
	altered := false
	for _, character := range raw {
		if (character >= 'a' && character <= 'z') || (character >= '0' && character <= '9') || character == '-' {
			builder.WriteRune(character)
		} else if character == ' ' || character == '_' {
			builder.WriteByte('-')
			altered = true
		} else {
			altered = true
		}
	}
	clean := strings.Trim(builder.String(), "-")
	for strings.Contains(clean, "--") {
		clean = strings.ReplaceAll(clean, "--", "-")
	}

	if altered || clean != raw {
		h := fnv.New32a()
		h.Write([]byte(raw))
		hashSuffix := fmt.Sprintf("%08x", h.Sum32())[:6]
		if len(clean) > 25 {
			clean = clean[:25]
		}
		clean = strings.Trim(clean, "-")
		if clean == "" {
			clean = fmt.Sprintf("app-%s", hashSuffix)
		} else {
			clean = fmt.Sprintf("%s-%s", clean, hashSuffix)
		}
	} else {
		if len(clean) > 32 {
			clean = clean[:32]
		}
		clean = strings.Trim(clean, "-")
		if clean == "" {
			return "app"
		}
	}
	return clean
}
