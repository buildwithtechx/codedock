package deploy

import (
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
	"io"
	"os"
	"testing"
	"time"
)

func TestLiveCancelledRolloutKeepsPreviousWorker(t *testing.T) {
	if os.Getenv("CODEDOCK_TEST_LIVE_DOCKER") != "1" {
		t.Skip("disposable Docker integration is opt-in")
	}
	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := docker.Close(); err != nil {
			t.Error(err)
		}
	})
	ctx, finish := context.WithTimeout(context.Background(), 2*time.Minute)
	defer finish()
	isolatedNetwork := "codedock-fixture-" + uuid.NewString()
	createdNetwork, err := docker.NetworkCreate(ctx, isolatedNetwork, network.CreateOptions{Labels: map[string]string{"codedock.fixture": "true"}})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := docker.NetworkRemove(cleanup, createdNetwork.ID); err != nil {
			t.Error(err)
		}
	})
	pulled, err := docker.ImagePull(ctx, "alpine:3.21", image.PullOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(io.Discard, pulled); err != nil {
		pulled.Close()
		t.Fatal(err)
	}
	if err := pulled.Close(); err != nil {
		t.Fatal(err)
	}
	store := &rolloutTestStore{}
	deployer := NewDeployer(docker, store)
	if err := deployer.SetRolloutDirectory(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	app := &models.AppService{ID: uuid.NewString(), RuntimeMode: models.RuntimeModeWorker, Status: models.AppServiceStatusRunning}
	opts := ContainerRunOptions{Name: utils.NormalizeContainerName(app.ID), ImageTag: "alpine:3.21", ServiceID: app.ID, RuntimeMode: app.RuntimeMode, InternalPort: 3000, Cmd: []string{"sleep", "600"}, ExtraLabels: map[string]string{}, NetworkName: isolatedNetwork}
	previous, err := deployer.containerManager.CreateAndStart(ctx, opts)
	if err != nil {
		t.Fatal(err)
	}
	app.ContainerID = previous
	t.Cleanup(func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := docker.ContainerRemove(cleanup, previous, container.RemoveOptions{Force: true}); err != nil {
			t.Error(err)
		}
	})
	operation, cancel := context.WithCancel(ctx)
	operation = WithDeploymentProgress(operation, func(phase string) error {
		if phase == "READINESS" {
			cancel()
			return operation.Err()
		}
		return nil
	})
	if _, err := deployer.replaceReplicas(operation, app, opts, 1, nil); err == nil {
		t.Fatal("cancelled rollout succeeded")
	}
	restored, err := docker.ContainerInspect(ctx, previous)
	if err != nil || restored.State == nil || !restored.State.Running || restored.Name != "/"+opts.Name {
		t.Fatalf("previous worker was not restored: %v", err)
	}
	if store.saved == nil || store.saved.ContainerID != previous {
		t.Fatal("previous runtime identity was not retained")
	}
}
