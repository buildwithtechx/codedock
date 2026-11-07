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
ssh_pwauth: false
disable_root: false
package_update: true
packages: [curl, ca-certificates, unattended-upgrades]
write_files:
  - path: /etc/ssh/sshd_config.d/10-codedock.conf
    permissions: "0644"
    content: |
      PasswordAuthentication no
      KbdInteractiveAuthentication no
      PermitRootLogin prohibit-password
runcmd:
  - [systemctl, reload, ssh]
  - [sh, -c, "curl -fsSL https://get.docker.com | sh"]
  - [systemctl, enable, --now, docker]
  - [systemctl, enable, --now, unattended-upgrades]
  - [sh, -c, "mkdir -p /var/lib/codedock && date -u +%%FT%%TZ > /var/lib/codedock/bootstrap-complete"]
`, name)
}
