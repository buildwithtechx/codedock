package networking

import (
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"codedock/internal/config"
)

var (
	CloudflareIPv4 = []string{
		"173.245.48.0/20",
		"103.21.244.0/22",
		"103.22.200.0/22",
		"103.31.4.0/22",
		"141.101.64.0/18",
		"108.162.192.0/18",
		"190.93.240.0/20",
		"188.114.96.0/20",
		"197.234.240.0/22",
		"198.41.128.0/17",
		"162.158.0.0/15",
		"104.16.0.0/13",
		"104.24.0.0/14",
		"172.64.0.0/13",
		"131.0.72.0/22",
	}

	CloudflareIPv6 = []string{
		"2400:cb00::/32",
		"2606:4700::/32",
		"2803:f800::/32",
		"2405:b500::/32",
		"2405:8100::/32",
		"2a06:98c0::/29",
		"2c0f:f248::/32",
	}

	headerRegex = regexp.MustCompile(`^[A-Za-z0-9-]{1,64}$`)
)

func ValidateRealIPHeader(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw != "" && headerRegex.MatchString(raw) {
		return raw
	}
	return "CF-Connecting-IP"
}

func ParseTrustedProxies(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	tokens := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || unicode.IsSpace(r)
	})

	var result []string
	seen := make(map[string]bool)

	for _, token := range tokens {
		s := strings.TrimSpace(token)
		if s == "" || len(s) > 64 || strings.Contains(s, "%") {
			continue
		}

		parts := strings.Split(s, "/")
		if len(parts) > 2 {
			continue
		}

		addr := parts[0]
		ip := net.ParseIP(addr)
		if ip == nil {
			continue
		}

		isV4 := ip.To4() != nil

		if len(parts) == 1 {
			cidr := s
			if isV4 {
				cidr = s + "/32"
			} else {
				cidr = s + "/128"
			}
			if !seen[cidr] {
				seen[cidr] = true
				result = append(result, cidr)
			}
			continue
		}

		prefixStr := parts[1]
		prefix, err := strconv.Atoi(prefixStr)
		if err != nil || prefix < 0 {
			continue
		}
		if isV4 && prefix > 32 {
			continue
		}
		if !isV4 && prefix > 128 {
			continue
		}

		_, ipNet, err := net.ParseCIDR(s)
		if err != nil || ipNet == nil {
			continue
		}

		normalized := ipNet.String()
		if !seen[normalized] {
			seen[normalized] = true
			result = append(result, normalized)
		}
	}

	return result
}

func GetEdgeTrustedProxies() []string {
	var result []string
	seen := make(map[string]bool)

	add := func(cidr string) {
		if !seen[cidr] {
			seen[cidr] = true
			result = append(result, cidr)
		}
	}

	for _, c := range CloudflareIPv4 {
		add(c)
	}
	for _, c := range CloudflareIPv6 {
		add(c)
	}

	envVal := os.Getenv("CODEDOCK_EDGE_TRUSTED_PROXIES")
	if envVal == "" {
		cfg := config.Get()
		if len(cfg.Traefik.TrustedProxies) > 0 {
			envVal = strings.Join(cfg.Traefik.TrustedProxies, ",")
		}
	}

	for _, extra := range ParseTrustedProxies(envVal) {
		add(extra)
	}

	return result
}

func ParseCIDRs(cidrs []string) []*net.IPNet {
	var nets []*net.IPNet
	for _, c := range cidrs {
		_, ipNet, err := net.ParseCIDR(c)
		if err == nil && ipNet != nil {
			nets = append(nets, ipNet)
		}
	}
	return nets
}

func IsIPInCIDRs(ip net.IP, nets []*net.IPNet) bool {
	if ip == nil {
		return false
	}
	for _, n := range nets {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}
