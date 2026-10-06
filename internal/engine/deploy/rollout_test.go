package deploy

import (
	"codedock.run/codedock/internal/models"
	"context"
	"encoding/json"
	"github.com/docker/docker/client"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type rolloutTestStore struct {
	DeployerStore
	saved *models.AppService
}

func (s *rolloutTestStore) GetServerSettings() (*models.ServerSettings, error) { return nil, nil }
func (s *rolloutTestStore) UpdateAppService(app *models.AppService) error {
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
		path := strings.TrimPrefix(r.URL.Path, "/v1.47")
		w.Header().Set("Content-Type", "application/json")
		if path == "/containers/json" {
			id := "old"
			if strings.Contains(r.URL.Query().Get("filters"), "codedock.rollout") {
				id = "new"
			}
			json.NewEncoder(w).Encode([]map[string]any{{"Id": id, "Names": []string{fixture.names[id]}, "State": "running"}})
			return
		}
		if path == "/containers/create" {
			fixture.names["new"] = "/" + r.URL.Query().Get("name")
			json.NewEncoder(w).Encode(map[string]string{"Id": "new"})
			return
		}
		parts := strings.Split(strings.Trim(path, "/"), "/")
		if len(parts) < 2 || parts[0] != "containers" {
			t.Errorf("unexpected Docker operation: %s", path)
			w.WriteHeader(500)
			return
		}
		id := parts[1]
		if r.Method == http.MethodDelete {
			fixture.removed[id] = true
			w.WriteHeader(204)
			return
		}
		if len(parts) < 3 {
			t.Errorf("unexpected Docker operation: %s", path)
			w.WriteHeader(500)
			return
		}
		switch parts[2] {
		case "rename":
			fixture.names[id] = "/" + r.URL.Query().Get("name")
			w.WriteHeader(204)
		case "start":
			fixture.starts[id]++
			w.WriteHeader(204)
		case "stop":
			w.WriteHeader(204)
		case "json":
			json.NewEncoder(w).Encode(map[string]any{"Id": id, "Name": fixture.names[id], "State": map[string]any{"Running": true, "Status": "running"}, "Config": map[string]any{}})
		default:
			t.Errorf("unexpected Docker operation: %s", path)
			w.WriteHeader(500)
		}
	}))
	t.Cleanup(server.Close)
	docker, err := client.NewClientWithOpts(client.WithHost(server.URL), client.WithHTTPClient(server.Client()), client.WithVersion("1.47"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := docker.Close(); err != nil {
			t.Error(err)
		}
	})
	store := &rolloutTestStore{}
	deployer := NewDeployer(docker, store)
	if err := deployer.SetRolloutDirectory(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	return deployer, store, fixture
}

func TestCancelledRolloutRestoresPreviousContainerAndMetadata(t *testing.T) {
	deployer, store, fixture := newRolloutFixture(t)
	app := &models.AppService{ID: "service", Name: "Application", ContainerID: "old", RuntimeMode: models.RuntimeModeWorker, Status: models.AppServiceStatusRunning}
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
	app := &models.AppService{ID: "service", ContainerID: "old", RuntimeMode: models.RuntimeModeWorker}
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
	journal := &models.RolloutJournal{ID: "interrupted", PreviousApp: models.AppService{ID: "service", ContainerID: "old"}, Previous: []models.RolloutContainer{{ID: "old", Name: "codedock-service", Running: true}}}
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
