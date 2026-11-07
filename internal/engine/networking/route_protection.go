package networking

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// EnsureTraefikAccessRoute writes the dynamic configuration that routes
// /_codedock/access to the Codedock backend with high priority.
func EnsureTraefikAccessRoute(dynamicDir string, backendPort int) error {
	if dynamicDir == "" {
		dynamicDir = DynamicDir
	}
	if backendPort <= 0 {
		backendPort = 8080
	}
	if err := os.MkdirAll(dynamicDir, 0755); err != nil {
		return fmt.Errorf("create traefik dynamic dir: %w", err)
	}

	content := fmt.Sprintf(`http:
  routers:
    codedock-access-router:
      rule: "PathPrefix(`+"`"+`/_codedock/access`+"`"+`)"
      priority: 100000
      service: "codedock-access-svc"
      entryPoints:
        - "web"
        - "websecure"
  services:
    codedock-access-svc:
      loadBalancer:
        servers:
          - url: "http://host.docker.internal:%d"
`, backendPort)

	destPath := filepath.Join(dynamicDir, "codedock-access.yaml")
	existing, err := os.ReadFile(destPath)
	if err == nil && string(existing) == content {
		return nil
	}
	return os.WriteFile(destPath, []byte(content), 0644)
}

// BuildRouteProtectionMiddlewareLabels generates Traefik ForwardAuth labels
// for a protected router.
func BuildRouteProtectionMiddlewareLabels(routerName string, backendPort int) (map[string]string, string) {
	if backendPort <= 0 {
		backendPort = 8080
	}
	mwName := sanitizeLabelKey(routerName + "-protect")
	authURL := fmt.Sprintf("http://host.docker.internal:%d/api/v1/edge/auth", backendPort)

	labels := map[string]string{
		fmt.Sprintf("traefik.http.middlewares.%s.forwardauth.address", mwName):            authURL,
		fmt.Sprintf("traefik.http.middlewares.%s.forwardauth.trustForwardHeader", mwName): "true",
		fmt.Sprintf("traefik.http.middlewares.%s.forwardauth.authResponseHeaders", mwName): strings.Join([]string{
			"Cache-Control",
			"Set-Cookie",
			"X-Codedock-Authorized",
		}, ","),
	}
	return labels, mwName + "@docker"
}
