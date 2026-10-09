package compose

import (
	"bytes"
	"codedock/internal/models"
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/google/uuid"
	"os"
	"testing"
	"time"
)

func TestLiveComposeRetainsVolumesAndResolvedVariables(t *testing.T) {
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
	runtime := NewStackRuntime(docker)
	id := uuid.NewString()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	content := `services:
  store:
    image: alpine:3.21
    environment:
      MESSAGE: ${MESSAGE:?required}
    command: [sh, -c, 'test -f /data/value || printf "%s" "$$MESSAGE" > /data/value; sleep 600']
    volumes: [data:/data]
    healthcheck:
      test: [CMD, test, -f, /data/value]
      interval: 1s
      timeout: 1s
      retries: 10
  web:
    image: alpine:3.21
    command: [sleep, '600']
    depends_on:
      store: {condition: service_healthy}
volumes:
  data: {}
`
	config, err := runtime.Canonical(ctx, id, content, map[string]string{"MESSAGE": "fixture-$literal"})
	if err != nil {
		t.Fatal(err)
	}
	stack := &models.ComposeStack{ID: id, Config: config}
	t.Cleanup(func() {
		cleanup, finish := context.WithTimeout(context.Background(), 30*time.Second)
		defer finish()
		if _, err := runtime.command(cleanup, id, config, nil, "down", "--volumes", "--remove-orphans"); err != nil {
			t.Errorf("dispose owned stack: %v", err)
		}
	})
	for attempt := 0; attempt < 2; attempt++ {
		if err := runtime.Apply(ctx, stack, func(string) error { return nil }); err != nil {
			t.Fatal(err)
		}
		results, err := runtime.Results(ctx, id)
		if err != nil || len(results) != 2 {
			t.Fatalf("per-service observation missing: %v %+v", err, results)
		}
		var storeID string
		for _, result := range results {
			if result.Name == "store" {
				storeID = result.ContainerID
				if result.Health != "healthy" {
					t.Fatal("dependent service was activated before health")
				}
			}
		}
		execResult, err := docker.ContainerExecCreate(ctx, storeID, container.ExecOptions{Cmd: []string{"sh", "-c", "cat /data/value; printf %s retained-marker > /data/value"}, AttachStdout: true, AttachStderr: true})
		if err != nil {
			t.Fatal(err)
		}
		attached, err := docker.ContainerExecAttach(ctx, execResult.ID, container.ExecAttachOptions{})
		if err != nil {
			t.Fatal(err)
		}
		var output bytes.Buffer
		_, err = stdcopy.StdCopy(&output, &output, attached.Reader)
		attached.Close()
		expected := "fixture-$literal"
		if attempt == 1 {
			expected = "retained-marker"
		}
		if err != nil || output.String() != expected {
			t.Fatalf("resolved variable changed through Compose: %q %v", output.String(), err)
		}
	}
	inspected, err := docker.VolumeInspect(ctx, stackProject(id)+"_data")
	if err != nil || inspected.Name == "" {
		t.Fatalf("owned volume missing after retry: %v", err)
	}

}
