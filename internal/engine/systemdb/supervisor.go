package systemdb

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/containerd/errdefs"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	ContainerName = "codedock-postgres"
	DBUser        = "codedock"
	DBName        = "codedock"
)

type Supervisor struct {
	docker  *client.Client
	name    string
	dataDir string
	image   string
	port    int
}

func NewSupervisor(docker *client.Client, name, dataDir, image string, port int) *Supervisor {
	return &Supervisor{docker: docker, name: name, dataDir: dataDir, image: image, port: port}
}

func (s *Supervisor) EnsureRunning(ctx context.Context) (string, error) {
	if s.docker == nil {
		return "", fmt.Errorf("docker client is required to provision the embedded database")
	}
	password, err := loadOrGeneratePassword(s.dataDir)
	if err != nil {
		return "", err
	}
	if _, err := s.docker.ContainerInspect(ctx, s.name); err != nil {
		if !errdefs.IsNotFound(err) {
			return "", fmt.Errorf("failed to inspect postgres container: %w", err)
		}
		if err := s.createContainer(ctx, password); err != nil {
			return "", fmt.Errorf("failed to create postgres container: %w", err)
		}
	}
	if err := s.docker.ContainerStart(ctx, s.name, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start postgres container: %w", err)
	}
	databaseURL := s.databaseURL(password)
	if err := s.waitReady(ctx, databaseURL); err != nil {
		return "", err
	}
	slog.Info("embedded postgres is running", "container", s.name)
	return databaseURL, nil
}

func (s *Supervisor) databaseURL(password string) string {
	return "postgres://" + DBUser + ":" + password + "@127.0.0.1:" + strconv.Itoa(s.port) + "/" + DBName + "?sslmode=disable"
}

func (s *Supervisor) createContainer(ctx context.Context, password string) error {
	out, err := s.docker.ImagePull(ctx, s.image, image.PullOptions{})
	if err == nil {
		defer out.Close()
		io.Copy(io.Discard, out)
	}
	dataPath := filepath.Join(s.dataDir, "postgres", "data")
	if err := os.MkdirAll(dataPath, 0o700); err != nil {
		return fmt.Errorf("failed to create postgres data dir: %w", err)
	}
	port, err := nat.NewPort("tcp", "5432")
	if err != nil {
		return fmt.Errorf("failed to parse postgres port: %w", err)
	}
	hostConfig := &container.HostConfig{
		Binds: []string{dataPath + ":/var/lib/postgresql/data"},
		PortBindings: nat.PortMap{
			port: []nat.PortBinding{{HostIP: "127.0.0.1", HostPort: strconv.Itoa(s.port)}},
		},
		RestartPolicy: container.RestartPolicy{Name: "unless-stopped"},
	}
	resp, err := s.docker.ContainerCreate(ctx, &container.Config{
		Image: s.image,
		Env: []string{
			"POSTGRES_USER=" + DBUser,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + DBName,
		},
		ExposedPorts: nat.PortSet{port: struct{}{}},
		Labels: map[string]string{
			"codedock.managed": "systemdb",
		},
		Healthcheck: &container.HealthConfig{
			Test:     []string{"CMD-SHELL", "pg_isready -U " + DBUser + " -d " + DBName},
			Interval: 5 * time.Second,
			Timeout:  3 * time.Second,
			Retries:  12,
		},
	}, hostConfig, &network.NetworkingConfig{}, nil, s.name)
	if err != nil {
		return err
	}
	slog.Info("created postgres container", "containerID", resp.ID)
	return nil
}

func (s *Supervisor) waitReady(ctx context.Context, databaseURL string) error {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to open postgres for readiness: %w", err)
	}
	defer db.Close()
	deadline := time.Now().Add(90 * time.Second)
	for {
		if err := db.PingContext(ctx); err == nil {
			return nil
		} else if time.Now().After(deadline) {
			return fmt.Errorf("timed out waiting for postgres on port %d: %w", s.port, err)
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("cancelled waiting for postgres: %w", ctx.Err())
		case <-time.After(500 * time.Millisecond):
		}
	}
}
