package commands

import (
	"context"
	"fmt"
	"os"
	"strings"

	"codedock.run/codedock/internal/repositories"
	"codedock.run/codedock/internal/utils"
)

func printDeploymentURL(settingsRepo *repositories.SettingsRepo, appName string) {
	wildcard := ""
	settings, err := settingsRepo.GetServerSettings(context.Background())
	if err == nil {
		wildcard = settings.DefaultWildcardDomain
	}
	if wildcard == "" {
		wildcard = os.Getenv("CODEDOCK_DOMAIN")
	}

	if wildcard != "" {
		cleanName := utils.SanitizeDomainName(appName)
		base := strings.TrimPrefix(wildcard, "*.")
		if strings.HasPrefix(base, "http") {
			base = strings.TrimPrefix(base, "https://")
			base = strings.TrimPrefix(base, "http://")
		}
		fmt.Printf("   URL: https://%s.%s\n", cleanName, base)
	} else {
		hostIP := os.Getenv("CODEDOCK_HOST_IP")
		if hostIP == "" {
			hostIP = "127.0.0.1"
		}
		cleanName := utils.SanitizeDomainName(appName)
		cleanIP := strings.ReplaceAll(hostIP, ".", "-")
		magicDomain := os.Getenv("CODEDOCK_MAGIC_DOMAIN")
		if magicDomain == "" {
			magicDomain = "sslip.io"
		}
		fmt.Printf("   URL: http://%s.%s.%s\n", cleanName, cleanIP, magicDomain)
	}
}
