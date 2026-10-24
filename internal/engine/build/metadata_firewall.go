package build

import (
	"fmt"
	"net"
	"strings"
)

const BuildMetadataPolicyVersion = "1"
const BuildMetadataPolicyLabel = "codedock.build-metadata-policy"

var BuildMetadataIPv4 = []string{
	"169.254.169.254/32", // EC2, Azure, GCE, Hetzner, DigitalOcean
	"169.254.170.2/32",   // ECS task credentials
	"169.254.170.23/32",  // EKS Pod Identity credentials
	"100.100.100.200/32", // Alibaba ECS
	"169.254.42.42/32",   // Scaleway
}

var BuildMetadataIPv6 = []string{
	"fd00:ec2::254/128", // EC2 IMDS
	"fd00:ec2::23/128",  // EKS Pod Identity credentials
	"fd20:ce::254/128",  // GCE IPv6-only metadata
}

// BuildMetadataFirewallScript generates the iptables/ip6tables script to isolate
// workers and BuildKit containers from hypervisor cloud metadata endpoints.
func BuildMetadataFirewallScript() string {
	chain := "CODEDOCK_METADATA_V1"
	var sb strings.Builder
	sb.WriteString("set -eu\n")

	type toolConfig struct {
		cmd string
		ips []string
		dns string
	}

	tools := []toolConfig{
		{"iptables", BuildMetadataIPv4, "169.254.169.254/32"},
		{"ip6tables", BuildMetadataIPv6, "fd20:ce::254/128"},
	}

	for _, tc := range tools {
		fmt.Fprintf(&sb, "command -v %s >/dev/null\n", tc.cmd)
		fmt.Fprintf(&sb, "%s -w 5 -t raw -N %s 2>/dev/null || %s -w 5 -t raw -S %s >/dev/null\n", tc.cmd, chain, tc.cmd, chain)
		fmt.Fprintf(&sb, "%s -w 5 -t raw -F %s\n", tc.cmd, chain)
		// Allow DNS resolution on port 53 (GCE uses its metadata IP for DNS)
		fmt.Fprintf(&sb, "%s -w 5 -t raw -A %s -d %s -p udp --dport 53 -j RETURN\n", tc.cmd, chain, tc.dns)
		fmt.Fprintf(&sb, "%s -w 5 -t raw -A %s -d %s -p tcp --dport 53 -j RETURN\n", tc.cmd, chain, tc.dns)

		for _, ip := range tc.ips {
			fmt.Fprintf(&sb, "%s -w 5 -t raw -A %s -d %s -j DROP\n", tc.cmd, chain, ip)
		}

		for _, hook := range []string{"PREROUTING", "OUTPUT"} {
			fmt.Fprintf(&sb, "%s -w 5 -t raw -I %s 1 -j %s\n", tc.cmd, hook, chain)
			fmt.Fprintf(&sb, "test \"$(%s -w 5 -t raw -S %s | sed -n '/^-A /{p;q;}')\" = '-A %s -j %s'\n", tc.cmd, hook, hook, chain)
		}

		for _, ip := range tc.ips {
			fmt.Fprintf(&sb, "%s -w 5 -t raw -C %s -d %s -j DROP\n", tc.cmd, chain, ip)
		}
	}

	return sb.String()
}

// BuildKitGuardEntrypoint produces the shell entrypoint script installed in build containers.
func BuildKitGuardEntrypoint() string {
	return "#!/bin/sh\n" +
		BuildMetadataFirewallScript() +
		fmt.Sprintf("printf '%%s\\n' 'CODEDOCK_BUILD_METADATA_GUARD_READY:%s' >&2\n", BuildMetadataPolicyVersion) +
		"exec /usr/bin/buildkitd \"$@\"\n"
}

// IsProtectedBuildKitConfig checks whether a container configuration has the firewall policy applied.
func IsProtectedBuildKitConfig(labels map[string]string, entrypoint []string) bool {
	if labels == nil || labels[BuildMetadataPolicyLabel] != BuildMetadataPolicyVersion {
		return false
	}
	if len(entrypoint) == 0 {
		return false
	}
	ep := strings.Join(entrypoint, " ")
	return strings.Contains(ep, "buildkit") || strings.Contains(ep, "codedock-buildkitd")
}

// ValidateOutboundTarget verifies that an outbound destination (e.g. for webhooks or backup storage)
// does not point to sensitive cloud metadata hypervisor endpoints.
func ValidateOutboundTarget(hostOrIP string) error {
	clean := strings.TrimSpace(hostOrIP)
	if clean == "" {
		return nil
	}
	// Strip port if present
	if h, _, err := net.SplitHostPort(clean); err == nil {
		clean = h
	}

	ip := net.ParseIP(clean)
	if ip == nil {
		// Resolve hostname if possible or check known literal hostnames
		if ips, err := net.LookupIP(clean); err == nil && len(ips) > 0 {
			ip = ips[0]
		}
	}

	if ip != nil {
		for _, cidr := range append(BuildMetadataIPv4, BuildMetadataIPv6...) {
			_, network, err := net.ParseCIDR(cidr)
			if err == nil && network.Contains(ip) {
				return fmt.Errorf("outbound target %s resolves to restricted cloud metadata destination %s", hostOrIP, cidr)
			}
		}
	}

	return nil
}
