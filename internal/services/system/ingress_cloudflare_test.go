package system

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCanonicalFormatting(t *testing.T) {
	input := map[string]any{
		"z":      1,
		"a":      "hello",
		"empty":  map[string]any{},
		"nilVal": nil,
		"nested": map[string]any{"b": 2, "a": 1},
	}
	res := Canonical(input)
	expected := `{"a":"hello","nested":{"a":1,"b":2},"z":1}`
	if res != expected {
		t.Fatalf("expected %s, got %s", expected, res)
	}
}

func TestMergeIngress(t *testing.T) {
	cfg := TunnelConfig{
		Ingress: []IngressRule{
			{Service: "http_status:404"},
		},
	}

	newRule := &IngressRule{
		Hostname: "app.example.com",
		Service:  "http://127.0.0.1:8080",
	}

	updated, err := MergeIngress(cfg, "app.example.com", nil, nil, newRule)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(updated.Ingress) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(updated.Ingress))
	}
	if updated.Ingress[0].Hostname != "app.example.com" {
		t.Fatalf("expected app.example.com as first rule, got %s", updated.Ingress[0].Hostname)
	}

	// Reject if existing rule modified outside CodeDock
	prevRule := &IngressRule{
		Hostname: "app.example.com",
		Service:  "http://127.0.0.1:9090",
	}
	_, err = MergeIngress(updated, "app.example.com", prevRule, nil, newRule)
	if err == nil {
		t.Fatalf("expected error on out-of-band modification")
	}

	// Wildcard conflict
	wildcardCfg := TunnelConfig{
		Ingress: []IngressRule{
			{Hostname: "*.example.com", Service: "http://127.0.0.1:80"},
			{Service: "http_status:404"},
		},
	}
	_, err = MergeIngress(wildcardCfg, "sub.example.com", nil, nil, newRule)
	if err == nil {
		t.Fatalf("expected wildcard collision error")
	}
}

func TestParseAccessAppDomainLess(t *testing.T) {
	// Launcher/SaaS application without domain or aud (upstream #1302)
	raw := map[string]any{
		"id":   "app-saas-1",
		"name": "Google Workspace SSO",
		"type": "saas",
		"policies": []any{
			map[string]any{"id": "pol-1"},
		},
	}
	app, err := ParseAccessApp(raw)
	if err != nil {
		t.Fatalf("unexpected error parsing domain-less access app: %v", err)
	}
	if app.ID != "app-saas-1" || app.Name != "Google Workspace SSO" || len(app.Policies) != 1 || app.Policies[0] != "pol-1" {
		t.Fatalf("unexpected parsed app: %+v", app)
	}
}

func TestCloudflareIngressServiceAPI(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/accounts":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"success":true,"result":[{"id":"acc-1","name":"Main Account"}]}`))
		case "/accounts/acc-1/cfd_tunnel":
			if r.Method == http.MethodPost {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"success":true,"result":{"id":"tun-new"}}`))
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"success":true,"result":[{"id":"tun-1","name":"existing","status":"healthy","remote_config":true}]}`))
			}
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	svc := NewCloudflareIngressService(ts.Client())
	svc.baseURL = ts.URL

	ctx := context.Background()
	accs, err := svc.ListAccounts(ctx, "test-token")
	if err != nil || len(accs) != 1 || accs[0].ID != "acc-1" {
		t.Fatalf("ListAccounts failed: %v, %+v", err, accs)
	}

	tuns, err := svc.ListTunnels(ctx, "test-token", "acc-1")
	if err != nil || len(tuns) != 1 || tuns[0].ID != "tun-1" || !tuns[0].Remote {
		t.Fatalf("ListTunnels failed: %v, %+v", err, tuns)
	}

	tunID, err := svc.CreateTunnel(ctx, "test-token", "acc-1", "new-tunnel")
	if err != nil || tunID != "tun-new" {
		t.Fatalf("CreateTunnel failed: %v, tunID: %s", err, tunID)
	}
}
