package deploy

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/docker/docker/client"

	appsmod "codedock/internal/modules/apps"
	settingsmodule "codedock/internal/modules/settings"
)

type rolloutTestStore struct {
	DeployerStore
	saved *appsmod.AppService
}

func (s *rolloutTestStore) GetServerSettings() (*settingsmodule.ServerSettings, error) {
	return nil, nil
}
func (s *rolloutTestStore) UpdateAppService(app *appsmod.AppService) error {
	copy := *app
	s.saved = &copy
	return nil
}

type rolloutFixture struct {
	names   map[string]string
	removed map[string]bool
	starts  map[string]int
}

func newRolloutFixture(t *testing.T) (*Deployer, *rolloutTestStore, *rolloutFixture) {
	t.Helper()
	fixture := &rolloutFixture{names: map[string]string{"old": "/codedock-service"}, removed: map[string]bool{}, starts: map[string]int{}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/v1.47/containers/json"):
			id := "old"
			if strings.Contains(r.URL.Query().Get("filters"), "codedock.rollout") {
				id = "new"
			}
			if fixture.removed[id] {
				_ = json.NewEncoder(w).Encode([]map[string]any{})
				return
			}
			_ = json.NewEncoder(w).Encode([]map[string]any{{"Id": id, "Names": []string{fixture.names[id]}, "State": "running"}})
		case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/v1.47/containers/create"):
			fixture.names["new"] = "/codedock-service"
			_ = json.NewEncoder(w).Encode(map[string]string{"Id": "new"})
		case r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/rename"):
			parts := strings.Split(r.URL.Path, "/")
			id := parts[3]
			fixture.names[id] = "/" + r.URL.Query().Get("name")
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/start"):
			parts := strings.Split(r.URL.Path, "/")
			id := parts[3]
			fixture.starts[id]++
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/stop"):
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/json"):
			parts := strings.Split(r.URL.Path, "/")
			id := parts[3]
			_ = json.NewEncoder(w).Encode(map[string]any{
				"Id":    id,
				"Name":  fixture.names[id],
				"State": map[string]any{"Running": true, "Status": "running"},
			})
		case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/v1.47/containers/"):
			parts := strings.Split(r.URL.Path, "/")
			id := parts[3]
			fixture.removed[id] = true
			delete(fixture.names, id)
			w.WriteHeader(http.StatusNoContent)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)
	docker, err := client.NewClientWithOpts(client.WithHost(server.URL), client.WithVersion("1.47"), client.WithHTTPClient(server.Client()))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = docker.Close() })
	store := &rolloutTestStore{}
	deployer := NewDeployer(docker, store)
	if err := deployer.SetRolloutDirectory(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	return deployer, store, fixture
}

func TestCancelledRolloutRestoresPreviousContainerAndMetadata(t *testing.T) {
	deployer, store, fixture := newRolloutFixture(t)
	app := &appsmod.AppService{ID: "service", Name: "Application", ContainerID: "old", RuntimeMode: appsmod.RuntimeModeWorker, Status: appsmod.AppServiceStatusRunning}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = WithDeploymentProgress(ctx, func(phase string) error {
		if phase == "READINESS" {
			cancel()
			return ctx.Err()
		}
		return nil
	})
	_, err := deployer.replaceReplicas(ctx, app, ContainerRunOptions{ImageTag: "image", ServiceID: app.ID, RuntimeMode: app.RuntimeMode, InternalPort: 3000, ExtraLabels: map[string]string{}}, 1, nil)
	if err == nil {
		t.Fatal("cancelled activation succeeded")
	}
	if fixture.removed["old"] || !fixture.removed["new"] || fixture.names["old"] != "/codedock-service" || fixture.starts["old"] != 1 {
		t.Fatalf("previous workload was not restored: %+v", fixture)
	}
	if store.saved == nil || store.saved.ContainerID != "old" {
		t.Fatal("previous service metadata was not restored")
	}
	if _, err := os.Stat(deployer.rolloutPath(app.ID)); !os.IsNotExist(err) {
		t.Fatal("successful rollback left a recovery journal")
	}
}

func TestSuccessfulRolloutRetiresPreviousOnlyAfterReadiness(t *testing.T) {
	deployer, store, fixture := newRolloutFixture(t)
	app := &appsmod.AppService{ID: "service", ContainerID: "old", RuntimeMode: appsmod.RuntimeModeWorker}
	ctx := WithDeploymentProgress(context.Background(), func(phase string) error {
		if fixture.removed["old"] {
			t.Errorf("old container removed during phase %s", phase)
		}
		return nil
	})
	id, err := deployer.replaceReplicas(ctx, app, ContainerRunOptions{ImageTag: "image", ServiceID: app.ID, RuntimeMode: app.RuntimeMode, InternalPort: 3000, ExtraLabels: map[string]string{}}, 1, nil)
	if err != nil || id != "new" {
		t.Fatalf("rollout failed: %v", err)
	}
	if !fixture.removed["old"] || fixture.removed["new"] || store.saved.ContainerID != "new" {
		t.Fatal("successful rollout did not activate the new container")
	}
}

func TestDaemonRecoveryRestoresUncommittedRollout(t *testing.T) {
	deployer, store, fixture := newRolloutFixture(t)
	fixture.names["old"], fixture.names["new"] = "/codedock-service-previous-recovery", "/codedock-service"
	journal := &RolloutJournal{ID: "interrupted", PreviousApp: appsmod.AppService{ID: "service", ContainerID: "old"}, Previous: []RolloutContainer{{ID: "old", Name: "codedock-service", Running: true}}}
	if err := deployer.saveRollout(journal); err != nil {
		t.Fatal(err)
	}
	if err := deployer.RecoverRollouts(context.Background()); err != nil {
		t.Fatal(err)
	}
	if fixture.removed["old"] || !fixture.removed["new"] || fixture.names["old"] != "/codedock-service" || store.saved.ContainerID != "old" {
		t.Fatal("daemon recovery did not preserve the previous workload")
	}
}
