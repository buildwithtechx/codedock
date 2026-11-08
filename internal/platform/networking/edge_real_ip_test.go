package networking

import (
	"os"
	"slices"
	"strings"
	"testing"
)

func TestParseTrustedProxies(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want []string
	}{
		{
			name: "empty input",
			raw:  "",
			want: nil,
		},
		{
			name: "single ipv4 cidr",
			raw:  "192.168.1.0/24",
			want: []string{"192.168.1.0/24"},
		},
		{
			name: "bare ipv4 address",
			raw:  "10.0.0.1",
			want: []string{"10.0.0.1/32"},
		},
		{
			name: "bare ipv6 address",
			raw:  "2001:db8::1",
			want: []string{"2001:db8::1/128"},
		},
		{
			name: "mixed comma and whitespace separated with duplicates",
			raw:  "10.0.0.0/8, 172.16.0.0/12\n192.168.0.0/16 10.0.0.0/8",
			want: []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"},
		},
		{
			name: "drops invalid prefix lengths and malformed entries",
			raw:  "10.0.0.0/33, 192.168.1.1/-5, not-an-ip, 2001:db8::/129, fe80::1%eth0, 1.2.3.4",
			want: []string{"1.2.3.4/32"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseTrustedProxies(tt.raw)
			if len(got) != len(tt.want) {
				t.Fatalf("ParseTrustedProxies(%q) len = %d; want %d (%v vs %v)", tt.raw, len(got), len(tt.want), got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("got[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestValidateRealIPHeader(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{name: "empty string", raw: "", want: "CF-Connecting-IP"},
		{name: "standard CF header", raw: "CF-Connecting-IP", want: "CF-Connecting-IP"},
		{name: "standard XFF", raw: "X-Forwarded-For", want: "X-Forwarded-For"},
		{name: "custom header", raw: "X-Real-IP", want: "X-Real-IP"},
		{name: "injection attempt with spaces", raw: "X-Real-IP; evil", want: "CF-Connecting-IP"},
		{name: "injection attempt with braces", raw: "X-Real-IP { test }", want: "CF-Connecting-IP"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateRealIPHeader(tt.raw); got != tt.want {
				t.Errorf("ValidateRealIPHeader(%q) = %q; want %q", tt.raw, got, tt.want)
			}
		})
	}
}

func TestGetEdgeTrustedProxies(t *testing.T) {
	os.Setenv("CODEDOCK_EDGE_TRUSTED_PROXIES", "10.100.0.0/16, 127.0.0.1")
	defer os.Unsetenv("CODEDOCK_EDGE_TRUSTED_PROXIES")

	proxies := GetEdgeTrustedProxies()
	if !slices.Contains(proxies, "173.245.48.0/20") {
		t.Errorf("expected Cloudflare IPv4 range to be present")
	}
	if !slices.Contains(proxies, "2400:cb00::/32") {
		t.Errorf("expected Cloudflare IPv6 range to be present")
	}
	if !slices.Contains(proxies, "10.100.0.0/16") {
		t.Errorf("expected operator declared proxy 10.100.0.0/16 to be present")
	}
	if !slices.Contains(proxies, "127.0.0.1/32") {
		t.Errorf("expected operator declared proxy 127.0.0.1/32 to be present")
	}
}

func TestTraefikCmdArgsTrustedProxies(t *testing.T) {
	os.Setenv("CODEDOCK_EDGE_TRUSTED_PROXIES", "10.200.0.0/24")
	defer os.Unsetenv("CODEDOCK_EDGE_TRUSTED_PROXIES")

	mgr := &TraefikManager{}
	args := mgr.buildTraefikCmdArgs()

	var foundWeb, foundWebsecure bool
	for _, arg := range args {
		if strings.HasPrefix(arg, "--entrypoints.web.forwardedheaders.trustedips=") {
			foundWeb = true
			if !strings.Contains(arg, "10.200.0.0/24") {
				t.Errorf("web entrypoint missing operator proxy: %s", arg)
			}
		}
		if strings.HasPrefix(arg, "--entrypoints.websecure.forwardedheaders.trustedips=") {
			foundWebsecure = true
			if !strings.Contains(arg, "10.200.0.0/24") {
				t.Errorf("websecure entrypoint missing operator proxy: %s", arg)
			}
		}
	}

	if !foundWeb || !foundWebsecure {
		t.Fatalf("expected trustedips flags on web and websecure entrypoints, got: %v", args)
	}
}
