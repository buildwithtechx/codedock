package networking

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func TestProxyAppliesCertificateSettingsOnRestart(t *testing.T) {
	for _, check := range []struct {
		name, email string
		running     bool
	}{
		{"running proxy needs certificates", "", true},
		{"stopped proxy needs certificates", "", false},
		{"running proxy has certificates", "owner@example.com", true},
		{"stopped proxy has certificates", "owner@example.com", false},
	} {
		t.Run(check.name, func(t *testing.T) {
			existingEmail := check.email
			var mu sync.Mutex
			var actions []string
			record := func(action string) { mu.Lock(); defer mu.Unlock(); actions = append(actions, action) }
			docker := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				path := strings.TrimPrefix(r.URL.Path, "/v1.45")
				w.Header().Set("Content-Type", "application/json")
				switch {
				case path == "/networks/codedock-network":
					if err := json.NewEncoder(w).Encode(map[string]string{"Id": "network"}); err != nil {
						t.Error(err)
					}
				case path == "/containers/codedock-traefik/json":
					args := []string{}
					if existingEmail != "" {
						args = append(args, "--certificatesresolvers.letsencrypt.acme.email="+existingEmail)
					}
					if err := json.NewEncoder(w).Encode(map[string]any{"Id": "proxy", "Config": map[string]any{"Cmd": args}, "State": map[string]any{"Running": check.running}}); err != nil {
						t.Error(err)
					}
				case path == "/containers/codedock-traefik/stop":
					record("stop")
					if !check.running {
						w.WriteHeader(http.StatusNotModified)
						return
					}
					w.WriteHeader(204)
				case path == "/containers/codedock-traefik" && r.Method == http.MethodDelete:
					record("remove")
					w.WriteHeader(204)
				case path == "/images/create":
					if err := json.NewEncoder(w).Encode(map[string]string{"status": "pulled"}); err != nil {
						t.Error(err)
					}
				case path == "/containers/create":
					record("create")
					var cfg container.Config
					if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
						t.Error(err)
					}
					if traefikCertificateEmail(cfg.Cmd) != "owner@example.com" {
						t.Error("new proxy omitted configured certificate email")
					}
					if err := json.NewEncoder(w).Encode(map[string]string{"Id": "proxy"}); err != nil {
						t.Error(err)
					}
				case path == "/containers/codedock-traefik/start":
					record("start")
					w.WriteHeader(204)
				default:
					t.Errorf("unexpected Docker request: %s %s", r.Method, path)
					w.WriteHeader(500)
				}
			}))
			defer docker.Close()
			cli, err := client.NewClientWithOpts(client.WithHost(docker.URL), client.WithVersion("1.45"))
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := cli.Close(); err != nil {
					t.Error(err)
				}
			}()
			manager := NewTraefikManager(cli, "owner@example.com")
			if err := manager.EnsureTraefikRunning(context.Background()); err != nil {
				t.Fatal(err)
			}
			expected := []string{"start"}
			if existingEmail == "" {
				expected = []string{"remove", "create", "start"}
				if check.running {
					expected = append([]string{"stop"}, expected...)
				}
			}
			mu.Lock()
			defer mu.Unlock()
			if !slices.Equal(actions, expected) {
				t.Fatalf("expected proxy actions %v, got %v", expected, actions)
			}
		})
	}
}
