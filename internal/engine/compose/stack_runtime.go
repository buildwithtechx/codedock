package compose

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"codedock/internal/models"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type StackRuntime struct{ docker *client.Client }

func NewStackRuntime(docker *client.Client) *StackRuntime { return &StackRuntime{docker: docker} }

func stackProject(id string) string { return "codedock-stack-" + strings.ReplaceAll(id, "-", "") }

func (r *StackRuntime) Canonical(ctx context.Context, id, content string, variables map[string]string) (string, error) {
	if err := ValidateStackInput(content); err != nil {
		return "", err
	}
	output, err := r.command(ctx, id, content, variables, "config", "--format", "json")
	if err != nil {
		return "", fmt.Errorf("Compose validation failed. Check required variables and supported configuration; Docker Compose v2 must be installed: %w", err)
	}
	if err := ValidateStackInput(string(output)); err != nil {
		return "", err
	}
	if err := validateStackOwnership(string(output), id); err != nil {
		return "", err
	}
	return string(output), nil
}

func (r *StackRuntime) Apply(ctx context.Context, stack *models.ComposeStack, progress func(string) error) error {
	if r.docker == nil {
		return fmt.Errorf("local Docker is unavailable")
	}
	if err := r.ValidatePorts(ctx, stack.ID, stack.Config); err != nil {
		return err
	}
	if err := progress("BUILDING"); err != nil {
		return err
	}
	if _, err := r.command(ctx, stack.ID, stack.Config, nil, "build"); err != nil {
		return fmt.Errorf("stack build failed; existing containers were preserved: %w", err)
	}
	var config struct {
		Services map[string]map[string]any `json:"services"`
	}
	if err := json.Unmarshal([]byte(stack.Config), &config); err != nil {
		return err
	}
	images := []string{}
	for name, service := range config.Services {
		if service["build"] == nil {
			images = append(images, name)
		}
	}
	if len(images) > 0 {
		args := append([]string{"pull"}, images...)
		if _, err := r.command(ctx, stack.ID, stack.Config, nil, args...); err != nil {
			return fmt.Errorf("stack image preparation failed; existing containers were preserved: %w", err)
		}
	}
	if err := progress("STARTING"); err != nil {
		return err
	}
	if _, err := r.command(ctx, stack.ID, stack.Config, nil, "up", "--detach", "--no-build", "--pull", "never", "--wait", "--wait-timeout", "120"); err != nil {
		return fmt.Errorf("stack activation was interrupted or failed; some services may have changed. Review per-service states before retrying: %w", err)
	}
	return progress("READINESS")
}

func (r *StackRuntime) Results(ctx context.Context, id string) ([]models.ComposeServiceResult, error) {
	if r.docker == nil {
		return nil, fmt.Errorf("local Docker is unavailable")
	}
	containers, err := r.docker.ContainerList(ctx, container.ListOptions{All: true, Filters: filters.NewArgs(filters.Arg("label", "com.docker.compose.project="+stackProject(id)))})
	if err != nil {
		return nil, fmt.Errorf("observe stack: %w", err)
	}
	results := []models.ComposeServiceResult{}
	for _, current := range containers {
		inspected, err := r.docker.ContainerInspect(ctx, current.ID)
		if err != nil {
			return nil, fmt.Errorf("inspect stack service: %w", err)
		}
		if inspected.State == nil {
			return nil, fmt.Errorf("stack service state unavailable")
		}
		health := "unconfigured"
		if inspected.State.Health != nil {
			health = inspected.State.Health.Status
		}
		results = append(results, models.ComposeServiceResult{Name: current.Labels["com.docker.compose.service"], ContainerID: current.ID, State: inspected.State.Status, Health: health, ExitCode: inspected.State.ExitCode})
	}
	return results, nil
}

func (r *StackRuntime) command(ctx context.Context, id, content string, variables map[string]string, args ...string) ([]byte, error) {
	directory, err := os.MkdirTemp("", "codedock-compose-")
	if err != nil {
		return nil, fmt.Errorf("prepare Compose command: %w", err)
	}
	defer os.RemoveAll(directory)
	configPath := filepath.Join(directory, "compose.json")
	if err := os.WriteFile(configPath, []byte(content), 0600); err != nil {
		return nil, err
	}
	envPath := filepath.Join(directory, ".env")
	if err := os.WriteFile(envPath, nil, 0600); err != nil {
		return nil, err
	}
	commandArgs := append([]string{"compose", "--project-name", stackProject(id), "--env-file", envPath, "--file", configPath}, args...)
	command := exec.CommandContext(ctx, "docker", commandArgs...)
	command.Dir = directory
	for _, key := range []string{"PATH", "Path", "SystemRoot", "TEMP", "TMP", "HOME", "USERPROFILE", "DOCKER_CONFIG", "DOCKER_TLS_VERIFY", "DOCKER_CERT_PATH"} {
		if value := os.Getenv(key); value != "" {
			command.Env = append(command.Env, key+"="+value)
		}
	}
	if r.docker != nil {
		command.Env = append(command.Env, "DOCKER_HOST="+r.docker.DaemonHost())
	}
	for key, value := range variables {
		normalized := strings.ToUpper(key)
		if !regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`).MatchString(key) || strings.HasPrefix(normalized, "COMPOSE_") || strings.HasPrefix(normalized, "DOCKER_") || normalized == "PATH" || normalized == "HOME" || normalized == "SYSTEMROOT" || normalized == "USERPROFILE" || normalized == "TEMP" || normalized == "TMP" {
			return nil, fmt.Errorf("reserved interpolation variable: %s", key)
		}
		command.Env = append(command.Env, key+"="+value)
	}
	var output limitedComposeOutput
	var stderr limitedComposeOutput
	command.Stdout, command.Stderr = &output, &stderr
	if err := command.Run(); err != nil {
		return nil, fmt.Errorf("Docker Compose command %s failed (%w)", args[0], err)
	}
	if args[0] == "config" && strings.Contains(stderr.String(), "variable is not set") {
		return nil, fmt.Errorf("an interpolation variable is missing; supply it explicitly")
	}
	return output.Bytes(), nil
}

type limitedComposeOutput struct{ bytes.Buffer }

func (w *limitedComposeOutput) Write(data []byte) (int, error) {
	if w.Len()+len(data) > 2*1024*1024 {
		return 0, fmt.Errorf("Compose output exceeds size limit")
	}
	return w.Buffer.Write(data)
}

var _ io.Writer = (*limitedComposeOutput)(nil)
