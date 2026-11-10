package deploy

import "testing"

func TestBuildTraefikLabelsRoutesEveryHostOnOnePort(t *testing.T) {
	manager := &ContainerManager{}
	labels := manager.buildTraefikLabels("svc", []RouteHost{
		{Host: "web.example.com"},
		{Host: "alias.example.com", PathPrefix: "/app"},
	}, 3000, "", nil)

	if got := labels["traefik.http.services.svc.loadbalancer.server.port"]; got != "3000" {
		t.Errorf("port label = %q", got)
	}
	if got := labels["traefik.http.routers.svc.rule"]; got != "Host(`web.example.com`)" {
		t.Errorf("primary rule = %q", got)
	}
	if got := labels["traefik.http.routers.svc.service"]; got != "svc" {
		t.Errorf("primary service = %q", got)
	}
	if got := labels["traefik.http.routers.svc-1.rule"]; got != "Host(`alias.example.com`) && PathPrefix(`/app`)" {
		t.Errorf("alias rule = %q", got)
	}
	if got := labels["traefik.http.routers.svc-1.service"]; got != "svc" {
		t.Errorf("alias service = %q", got)
	}
	if labels["traefik.http.routers.svc-1.tls"] != "true" || labels["traefik.http.routers.svc-1.tls.certresolver"] != "letsencrypt" {
		t.Errorf("alias tls labels missing: %#v", labels)
	}
}
