package hetzner

import (
	"fmt"
	"strings"
)

func ProvisionUserData(hostname string) string {
	safe := make([]rune, 0, len(hostname))
	for _, symbol := range strings.ToLower(hostname) {
		if (symbol >= 'a' && symbol <= 'z') || (symbol >= '0' && symbol <= '9') || symbol == '-' {
			safe = append(safe, symbol)
		}
	}
	name := strings.Trim(string(safe), "-")
	if name == "" {
		name = "codedock-managed"
	}
	if len(name) > 63 {
		name = name[:63]
	}
	return fmt.Sprintf(`#cloud-config
hostname: %s
manage_etc_hosts: true
packages: [curl]
runcmd:
  - [sh, -c, "curl -fsSL https://get.docker.com | sh"]
  - [systemctl, enable, --now, docker]
`, name)
}
