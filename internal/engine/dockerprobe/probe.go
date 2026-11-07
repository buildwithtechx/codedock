package dockerprobe

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Runner interface {
	Run(ctx context.Context, cmd string) (string, error)
	RunPipe(ctx context.Context, cmd string, stdin io.Reader, stdout io.Writer) error
}

type PsRow struct {
	ID    string `json:"ID"`
	Names string `json:"Names"`
	Image string `json:"Image"`
	State string `json:"State"`
	Ports string `json:"Ports"`
}

type InspectMount struct {
	Source      string `json:"Source"`
	Destination string `json:"Destination"`
	Type        string `json:"Type"`
	Name        string `json:"Name"`
}

type InspectResult struct {
	ID     string `json:"Id"`
	Name   string `json:"Name"`
	Config struct {
		Image  string            `json:"Image"`
		Env    []string          `json:"Env"`
		Labels map[string]string `json:"Labels"`
	} `json:"Config"`
	HostConfig struct {
		Binds        []string `json:"Binds"`
		PortBindings map[string][]struct {
			HostIP   string `json:"HostIp"`
			HostPort string `json:"HostPort"`
		} `json:"PortBindings"`
	} `json:"HostConfig"`
	Mounts []InspectMount `json:"Mounts"`
	State  struct {
		Status     string `json:"Status"`
		Running    bool   `json:"Running"`
		Restarting bool   `json:"Restarting"`
	} `json:"State"`
	NetworkSettings struct {
		Ports map[string][]struct {
			HostIP   string `json:"HostIp"`
			HostPort string `json:"HostPort"`
		} `json:"Ports"`
	} `json:"NetworkSettings"`
}

type Container struct {
	ID             string
	Name           string
	Image          string
	State          string
	Running        bool
	Ports          []string
	HostPorts      []string
	Env            map[string]string
	Volumes        []VolumeRef
	Labels         map[string]string
	ComposeProject string
}

type VolumeRef struct {
	Name        string
	Source      string
	Destination string
	Type        string
}

const MaskedSecret = "********"

func ListContainers(ctx context.Context, runner Runner) ([]PsRow, error) {
	out, err := runner.Run(ctx, `docker ps -a --format '{"ID":"{{.ID}}","Names":"{{.Names}}","Image":"{{.Image}}","State":"{{.State}}","Ports":"{{.Ports}}"}'`)
	if err != nil {
		return nil, fmt.Errorf("docker ps: %w", err)
	}
	var rows []PsRow
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var row PsRow
		if err := json.Unmarshal([]byte(line), &row); err == nil {
			rows = append(rows, row)
		}
	}
	return rows, nil
}

func InspectContainer(ctx context.Context, runner Runner, id string) (Container, error) {
	if id == "" || strings.ContainsAny(id, " \t\n;|&$`'\"\\") {
		return Container{}, fmt.Errorf("invalid container reference")
	}
	out, err := runner.Run(ctx, "docker inspect "+id)
	if err != nil {
		return Container{}, fmt.Errorf("docker inspect: %w", err)
	}
	var results []InspectResult
	if err := json.Unmarshal([]byte(out), &results); err != nil || len(results) == 0 {
		return Container{}, fmt.Errorf("docker inspect returned no data")
	}
	return convertInspect(results[0]), nil
}

func ProbeAll(ctx context.Context, runner Runner, progress func(step, detail string)) ([]Container, error) {
	if progress != nil {
		progress("LIST", "Listing Docker containers")
	}
	rows, err := ListContainers(ctx, runner)
	if err != nil {
		return nil, err
	}
	containers := make([]Container, 0, len(rows))
	for _, row := range rows {
		if progress != nil {
			progress("INSPECT", "Inspecting "+row.Names)
		}
		container, err := InspectContainer(ctx, runner, row.ID)
		if err != nil {
			continue
		}
		containers = append(containers, container)
	}
	return containers, nil
}

func convertInspect(res InspectResult) Container {
	container := Container{
		ID:      res.ID,
		Name:    strings.TrimPrefix(res.Name, "/"),
		Image:   res.Config.Image,
		State:   res.State.Status,
		Running: res.State.Running && !res.State.Restarting,
		Env:     map[string]string{},
		Labels:  res.Config.Labels,
	}
	if container.Labels == nil {
		container.Labels = map[string]string{}
	}
	for _, entry := range res.Config.Env {
		if key, value, ok := strings.Cut(entry, "="); ok {
			container.Env[key] = value
		}
	}
	for binding, published := range res.HostConfig.PortBindings {
		container.Ports = append(container.Ports, binding)
		for _, target := range published {
			if target.HostPort != "" {
				container.HostPorts = append(container.HostPorts, target.HostPort)
			}
		}
	}
	for binding, published := range res.NetworkSettings.Ports {
		seen := false
		for _, existing := range container.Ports {
			if existing == binding {
				seen = true
			}
		}
		if !seen {
			container.Ports = append(container.Ports, binding)
		}
		for _, target := range published {
			if target.HostPort != "" {
				container.HostPorts = append(container.HostPorts, target.HostPort)
			}
		}
	}
	for _, mount := range res.Mounts {
		name := mount.Name
		if name == "" {
			name = mount.Source
		}
		container.Volumes = append(container.Volumes, VolumeRef{Name: name, Source: mount.Source, Destination: mount.Destination, Type: mount.Type})
	}
	container.ComposeProject = container.Labels["com.docker.compose.project"]
	return container
}

func MaskEnv(env map[string]string) map[string]string {
	masked := make(map[string]string, len(env))
	for key := range env {
		masked[key] = MaskedSecret
	}
	return masked
}

func EnvKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for key := range env {
		keys = append(keys, key)
	}
	return keys
}

func DetectRoutes(container Container) []string {
	var routes []string
	for _, label := range []string{
		container.Labels["traefik.http.routers.migrated.rule"],
		container.Labels["com.codedock.domain"],
	} {
		if label != "" {
			routes = append(routes, label)
		}
	}
	for key, value := range container.Labels {
		if strings.HasPrefix(key, "traefik.http.routers.") && strings.HasSuffix(key, ".rule") {
			routes = append(routes, value)
		}
	}
	return routes
}
