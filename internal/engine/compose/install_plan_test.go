package compose

import (
	"strings"
	"testing"
)

func installTestTemplate() ComposeTemplate {
	return ComposeTemplate{
		Services: map[string]ComposeService{
			"web": {
				Image:       "example/web:1.0",
				Ports:       []string{"8080:80"},
				Environment: []string{"APP_ENV=production", "APP_SECRET=${secret:APP_SECRET}"},
				Volumes:     []string{"web-data:/data"},
				DependsOn:   []string{"worker"},
			},
			"worker": {
				Image: "example/worker:1.0",
			},
		},
		Volumes: map[string]any{"web-data": nil},
		XCodedock: &CodedockMetadata{
			IsOneClick:  true,
			Name:        "Example",
			Description: "Example application",
			Secrets:     []CodedockSecretSpec{{Var: "APP_SECRET", Label: "App secret", Length: 16}},
		},
	}
}

func fixedSecretGenerator(value string) SecretGenerator {
	return func(length int) (string, error) { return value, nil }
}

func TestResolveInstallPlan(t *testing.T) {
	plan, err := ResolveInstallPlan("example", installTestTemplate(), InstallRequest{
		Name:     "Demo",
		HostPort: 9090,
		Domain:   "demo.example.com",
	}, fixedSecretGenerator("generated-secret-value"))
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.GeneratedSecrets) != 1 || plan.GeneratedSecrets[0] != "APP_SECRET" {
		t.Fatalf("unexpected generated secrets: %v", plan.GeneratedSecrets)
	}
	if plan.Digest == "" {
		t.Fatal("expected plan digest")
	}
	var web *InstallServicePlan
	for i := range plan.Services {
		if plan.Services[i].Service == "web" {
			web = &plan.Services[i]
		}
	}
	if web == nil {
		t.Fatal("expected web service plan")
	}
	joined := strings.Join(web.Env, "\n")
	if !strings.Contains(joined, "APP_SECRET=generated-secret-value") {
		t.Errorf("expected resolved secret, got %v", web.Env)
	}
	if len(web.Ports) != 1 || web.Ports[0] != "9090:80" {
		t.Errorf("expected host port override, got %v", web.Ports)
	}
	if len(plan.Volumes) != 1 || plan.Volumes[0].Name != "app-demo-web-data" {
		t.Errorf("unexpected volumes: %+v", plan.Volumes)
	}
	if !strings.Contains(plan.ComposeYAML, "demo.example.com") {
		t.Errorf("expected domain labels in document:\n%s", plan.ComposeYAML)
	}
	if !strings.Contains(plan.ComposeYAML, "app-demo-web-data") {
		t.Errorf("expected prefixed volume in document:\n%s", plan.ComposeYAML)
	}
}

func TestInstallPlanMasked(t *testing.T) {
	plan, err := ResolveInstallPlan("example", installTestTemplate(), InstallRequest{Name: "Demo"}, fixedSecretGenerator("top-secret-value"))
	if err != nil {
		t.Fatal(err)
	}
	masked := plan.Masked()
	for _, service := range masked.Services {
		for _, entry := range service.Env {
			if strings.Contains(entry, "top-secret-value") {
				t.Fatalf("masked plan leaks secret: %v", masked.Services)
			}
		}
	}
	if strings.Contains(masked.ComposeYAML, "top-secret-value") {
		t.Fatalf("masked document leaks secret:\n%s", masked.ComposeYAML)
	}
	if masked.Digest != plan.Digest {
		t.Error("masked plan must keep the digest")
	}
}

func TestResolveInstallPlanRejectsBadInput(t *testing.T) {
	tmpl := installTestTemplate()
	if _, err := ResolveInstallPlan("example", tmpl, InstallRequest{}, fixedSecretGenerator("value")); err == nil {
		t.Error("expected missing name error")
	}
	badPort := InstallRequest{Name: "Demo", HostPort: 99999}
	if _, err := ResolveInstallPlan("example", tmpl, badPort, nil); err == nil {
		t.Error("expected host port range error")
	}
	shortSecret := InstallRequest{Name: "Demo", Secrets: map[string]string{"APP_SECRET": "short"}}
	if _, err := ResolveInstallPlan("example", tmpl, shortSecret, nil); err == nil {
		t.Error("expected short secret error")
	}
	badDomain := InstallRequest{Name: "Demo", Domain: "not a domain"}
	if _, err := ResolveInstallPlan("example", tmpl, badDomain, fixedSecretGenerator("generated-secret-value")); err == nil {
		t.Error("expected domain error")
	}
	provided := InstallRequest{Name: "Demo", Secrets: map[string]string{"APP_SECRET": "provided-secret-value"}}
	plan, err := ResolveInstallPlan("example", tmpl, provided, fixedSecretGenerator("generated-secret-value"))
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.GeneratedSecrets) != 0 {
		t.Errorf("provided secrets must not be regenerated: %v", plan.GeneratedSecrets)
	}
}
