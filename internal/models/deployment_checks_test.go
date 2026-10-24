package models

import (
	"reflect"
	"testing"
)

func TestResolveGitHubDeploymentChecks(t *testing.T) {
	def := DefaultGitHubDeploymentChecks()
	if !def.Enabled || !def.Deployment || !def.IncludeErrors {
		t.Fatalf("unexpected defaults: %+v", def)
	}
	if !reflect.DeepEqual(def.Services, []string{"all"}) {
		t.Fatalf("expected all services default, got: %v", def.Services)
	}

	nilResolved := ResolveGitHubDeploymentChecks(nil)
	if !reflect.DeepEqual(nilResolved, def) {
		t.Fatalf("expected nil to resolve to default, got: %+v", nilResolved)
	}

	custom := &GitHubDeploymentChecks{
		Enabled:       false,
		Deployment:    true,
		Services:      []string{"api", "web"},
		IncludeErrors: false,
	}
	resolved := ResolveGitHubDeploymentChecks(custom)
	if resolved.Enabled != false || !reflect.DeepEqual(resolved.Services, []string{"api", "web"}) {
		t.Fatalf("expected custom preferences preserved, got: %+v", resolved)
	}
}

func TestGitHubDeploymentChecks_TargetsService(t *testing.T) {
	checks := DefaultGitHubDeploymentChecks()
	if !checks.TargetsService("api") {
		t.Fatalf("expected 'all' to target any service")
	}

	checks.Services = []string{"web", "worker"}
	if !checks.TargetsService("web") {
		t.Fatalf("expected 'web' to be targeted")
	}
	if checks.TargetsService("db") {
		t.Fatalf("expected 'db' to not be targeted")
	}

	checks.Enabled = false
	if checks.TargetsService("web") {
		t.Fatalf("expected disabled checks to target nothing")
	}
}
