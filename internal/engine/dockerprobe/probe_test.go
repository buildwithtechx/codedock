package dockerprobe

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"
)

type scriptRunner struct {
	outputs map[string]string
	calls   []string
	pipes   []string
	pipeErr error
}

func (r *scriptRunner) Run(ctx context.Context, cmd string) (string, error) {
	r.calls = append(r.calls, cmd)
	for prefix, out := range r.outputs {
		if strings.HasPrefix(cmd, prefix) {
			return out, nil
		}
	}
	return "", fmt.Errorf("unexpected command: %s", cmd)
}

func (r *scriptRunner) RunPipe(ctx context.Context, cmd string, stdin io.Reader, stdout io.Writer) error {
	r.pipes = append(r.pipes, cmd)
	if r.pipeErr != nil {
		return r.pipeErr
	}
	if stdin != nil {
		_, _ = io.Copy(io.Discard, stdin)
	}
	if stdout != nil {
		_, _ = stdout.Write([]byte("piped-bytes"))
	}
	return nil
}

const psFixture = `{"ID":"abc123def456","Names":"web","Image":"nginx:1.25","State":"running","Ports":"0.0.0.0:8080->80/tcp"}
{"ID":"def456abc123","Names":"db","Image":"postgres:16","State":"exited","Ports":""}
`

const inspectFixture = `[{
	"Id": "abc123def456",
	"Name": "/web",
	"Config": {"Image": "nginx:1.25", "Env": ["APP_SECRET=topsecret", "PORT=80"], "Labels": {"com.docker.compose.project": "shop"}},
	"HostConfig": {"Binds": [], "PortBindings": {"80/tcp": [{"HostIp": "0.0.0.0", "HostPort": "8080"}]}},
	"Mounts": [{"Source": "webdata", "Destination": "/usr/share/nginx/html", "Type": "volume", "Name": "webdata"}],
	"State": {"Status": "running", "Running": true, "Restarting": false},
	"NetworkSettings": {"Ports": {}}
}]`

func TestListContainersParsesRows(t *testing.T) {
	runner := &scriptRunner{outputs: map[string]string{"docker ps": psFixture}}
	rows, err := ListContainers(context.Background(), runner)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(rows) != 2 || rows[0].Names != "web" || rows[1].Image != "postgres:16" {
		t.Fatalf("unexpected rows: %#v", rows)
	}
}

func TestInspectContainerParsesDetail(t *testing.T) {
	runner := &scriptRunner{outputs: map[string]string{"docker inspect": inspectFixture}}
	container, err := InspectContainer(context.Background(), runner, "abc123def456")
	if err != nil {
		t.Fatalf("inspect: %v", err)
	}
	if container.Name != "web" || container.Image != "nginx:1.25" || !container.Running {
		t.Fatalf("unexpected container: %#v", container)
	}
	if container.Env["APP_SECRET"] != "topsecret" {
		t.Fatalf("expected real env, got %#v", container.Env)
	}
	if len(container.HostPorts) != 1 || container.HostPorts[0] != "8080" {
		t.Fatalf("expected host port 8080, got %#v", container.HostPorts)
	}
	if len(container.Volumes) != 1 || container.Volumes[0].Name != "webdata" {
		t.Fatalf("expected webdata volume, got %#v", container.Volumes)
	}
	if container.ComposeProject != "shop" {
		t.Fatalf("expected compose project shop, got %q", container.ComposeProject)
	}
}

func TestInspectContainerRejectsInjection(t *testing.T) {
	runner := &scriptRunner{}
	if _, err := InspectContainer(context.Background(), runner, "abc; rm -rf /"); err == nil {
		t.Fatal("expected invalid reference to fail")
	}
	if len(runner.calls) != 0 {
		t.Fatal("expected no command to run")
	}
}

func TestMaskEnvHidesValues(t *testing.T) {
	masked := MaskEnv(map[string]string{"A": "1", "B": "2"})
	for key, value := range masked {
		if value != MaskedSecret {
			t.Fatalf("expected masked %s, got %q", key, value)
		}
	}
	if len(EnvKeys(map[string]string{"A": "1"})) != 1 {
		t.Fatal("expected env keys")
	}
}

func TestTransferImageStreamsSaveToLoad(t *testing.T) {
	source := &scriptRunner{}
	target := &scriptRunner{}
	if err := TransferImage(context.Background(), source, target, "nginx:1.25", nil); err != nil {
		t.Fatalf("transfer: %v", err)
	}
	if len(source.pipes) != 1 || !strings.Contains(source.pipes[0], "docker save nginx:1.25") {
		t.Fatalf("expected docker save, got %v", source.pipes)
	}
	if len(target.pipes) != 1 || !strings.Contains(target.pipes[0], "docker load") {
		t.Fatalf("expected docker load, got %v", target.pipes)
	}
	if err := TransferImage(context.Background(), source, target, "evil;id", nil); err == nil {
		t.Fatal("expected invalid image to fail")
	}
}

func TestTransferVolumeCreatesTargetFirst(t *testing.T) {
	source := &scriptRunner{}
	target := &scriptRunner{outputs: map[string]string{"docker volume ls": "", "docker volume create": "webdata\n"}}
	if err := TransferVolume(context.Background(), source, target, "webdata", "webdata", nil); err != nil {
		t.Fatalf("transfer: %v", err)
	}
	joined := strings.Join(target.calls, "\n")
	if !strings.Contains(joined, "docker volume create webdata") {
		t.Fatalf("expected volume create, got %v", target.calls)
	}
}

func TestVolumeHasDataReadsListing(t *testing.T) {
	full := &scriptRunner{outputs: map[string]string{"docker run": "index.html\n"}}
	has, err := VolumeHasData(context.Background(), full, "webdata")
	if err != nil || !has {
		t.Fatalf("expected data, got %v %v", has, err)
	}
	empty := &scriptRunner{outputs: map[string]string{"docker run": "\n"}}
	has, err = VolumeHasData(context.Background(), empty, "webdata")
	if err != nil || has {
		t.Fatalf("expected empty, got %v %v", has, err)
	}
}
