package deploy

import (
	"testing"

	"codedock/internal/models"
)

func TestNormalizeRouteHostsKeepsAliasesSharingAPort(t *testing.T) {
	hosts := NormalizeRouteHosts("web.example.com", []models.DomainConfig{
		{DomainName: "web.example.com"},
		{DomainName: "alias.example.com"},
		{DomainName: "alias.example.com"},
		{DomainName: "  Other.Example.com ", PathPrefix: "app"},
	})
	want := []RouteHost{
		{Host: "web.example.com", PathPrefix: ""},
		{Host: "alias.example.com", PathPrefix: ""},
		{Host: "other.example.com", PathPrefix: "/app"},
	}
	if len(hosts) != len(want) {
		t.Fatalf("expected %d hosts, got %d (%#v)", len(want), len(hosts), hosts)
	}
	for i := range want {
		if hosts[i] != want[i] {
			t.Errorf("host %d = %#v, want %#v", i, hosts[i], want[i])
		}
	}
}

func TestNormalizeRouteHostsPrimaryFirstAndEmpty(t *testing.T) {
	if hosts := NormalizeRouteHosts("", nil); len(hosts) != 0 {
		t.Fatalf("expected no hosts, got %#v", hosts)
	}
	hosts := NormalizeRouteHosts("primary.example.com", []models.DomainConfig{{DomainName: ""}})
	if len(hosts) != 1 || hosts[0].Host != "primary.example.com" {
		t.Fatalf("expected only the primary, got %#v", hosts)
	}
}
