package compose

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/container"
	"gopkg.in/yaml.v3"
)

func ValidateStackInput(content string) error {
	if len(content) == 0 || len(content) > 1024*1024 {
		return fmt.Errorf("Compose content must be between 1 byte and 1 MB")
	}
	var document map[string]any
	if err := yaml.Unmarshal([]byte(content), &document); err != nil {
		return fmt.Errorf("invalid Compose document: %w", err)
	}
	if err := allowedComposeKeys(document, "document", "name services networks volumes version"); err != nil {
		return err
	}
	services, ok := document["services"].(map[string]any)
	if !ok || len(services) == 0 {
		return fmt.Errorf("Compose services are required")
	}
	dependencies := make(map[string][]string)
	for name, value := range services {
		service, ok := value.(map[string]any)
		if !ok {
			return fmt.Errorf("invalid service: %s", name)
		}
		if err := allowedComposeKeys(service, name, "image build environment command entrypoint ports expose volumes networks depends_on healthcheck restart deploy cpus mem_limit mem_reservation user working_dir stop_signal stop_grace_period init tty stdin_open labels logging pull_policy scale attach"); err != nil {
			return err
		}
		if service["image"] == nil && service["build"] == nil {
			return fmt.Errorf("service %s needs image or build", name)
		}
		if build, ok := service["build"].(map[string]any); ok {
			if err := allowedComposeKeys(build, name+".build", "context dockerfile dockerfile_inline args target cache_from cache_to labels no_cache pull platforms network shm_size extra_hosts isolation"); err != nil {
				return err
			}
			if err := validateRemoteContext(fmt.Sprint(build["context"])); err != nil {
				return err
			}
		} else if build, ok := service["build"].(string); ok {
			if err := validateRemoteContext(build); err != nil {
				return err
			}
		}
		if volumes, ok := service["volumes"].([]any); ok {
			for _, value := range volumes {
				if mount, ok := value.(string); ok {
					parts := strings.Split(mount, ":")
					if len(parts) < 2 || strings.ContainsAny(parts[0], "/\\.") {
						return fmt.Errorf("service %s: use declared named volumes; uploaded Compose bind paths and anonymous volumes are unsupported", name)
					}
				} else if mount, ok := value.(map[string]any); !ok || mount["type"] != "volume" || mount["source"] == nil {
					return fmt.Errorf("service %s: only named volume mounts are supported", name)
				}
			}
		}
		if labels, ok := service["labels"].(map[string]any); ok {
			for key := range labels {
				if strings.HasPrefix(key, "com.docker.compose.") || strings.HasPrefix(key, "codedock.") {
					return fmt.Errorf("service %s uses reserved ownership labels", name)
				}
			}
		}
		switch depends := service["depends_on"].(type) {
		case []any:
			for _, dependency := range depends {
				dependencies[name] = append(dependencies[name], fmt.Sprint(dependency))
			}
		case map[string]any:
			for dependency := range depends {
				dependencies[name] = append(dependencies[name], dependency)
			}
		}
	}
	for _, kind := range []string{"networks", "volumes"} {
		resources, _ := document[kind].(map[string]any)
		for name, value := range resources {
			settings, _ := value.(map[string]any)
			if external, ok := settings["external"].(bool); ok && external {
				return fmt.Errorf("external %s %s cannot be safely owned by this stack", kind, name)
			}
			if kind == "volumes" {
				if err := allowedComposeKeys(settings, "volume "+name, "driver labels name external"); err != nil {
					return err
				}
				if driver, ok := settings["driver"]; ok && driver != "local" {
					return fmt.Errorf("unsupported volume driver for %s", name)
				}
			}
		}
	}
	active, visited := map[string]bool{}, map[string]bool{}
	var visit func(string) error
	visit = func(name string) error {
		if _, ok := services[name]; !ok {
			return fmt.Errorf("unknown dependency: %s", name)
		}
		if active[name] {
			return fmt.Errorf("dependency cycle at %s", name)
		}
		if visited[name] {
			return nil
		}
		active[name] = true
		for _, dependency := range dependencies[name] {
			if err := visit(dependency); err != nil {
				return err
			}
		}
		active[name], visited[name] = false, true
		return nil
	}
	for name := range services {
		if err := visit(name); err != nil {
			return err
		}
	}
	return nil
}

func allowedComposeKeys(value map[string]any, path, allowed string) error {
	for key := range value {
		if strings.HasPrefix(key, "x-") {
			continue
		}
		if !strings.Contains(" "+allowed+" ", " "+key+" ") {
			return fmt.Errorf("unsupported Compose field %s.%s; remove or configure it before applying", path, key)
		}
	}
	return nil
}

func validateRemoteContext(value string) error {
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme != "https" || parsed.User != nil || parsed.Port() != "" {
		return fmt.Errorf("uploaded Compose builds require an HTTPS Git context; local build directories require repository source support")
	}
	if parsed.Hostname() != "github.com" && parsed.Hostname() != "gitlab.com" && parsed.Hostname() != "bitbucket.org" {
		return fmt.Errorf("unsupported remote build context host")
	}
	return nil
}

func (r *StackRuntime) ValidatePorts(ctx context.Context, id, content string) error {
	if r.docker == nil {
		return fmt.Errorf("local Docker unavailable")
	}
	var config struct {
		Services map[string]struct {
			Ports []struct {
				Published string `json:"published"`
				Protocol  string `json:"protocol"`
				HostIP    string `json:"host_ip"`
			} `json:"ports"`
		} `json:"services"`
	}
	if err := json.Unmarshal([]byte(content), &config); err != nil {
		return err
	}
	ports := map[string]string{}
	for name, service := range config.Services {
		for _, port := range service.Ports {
			number, err := strconv.Atoi(port.Published)
			if err != nil || number < 1 || number > 65535 {
				return fmt.Errorf("service %s: publish a fixed port between 1 and 65535", name)
			}
			key := port.Published + "/" + port.Protocol
			if previous := ports[key]; previous != "" {
				return fmt.Errorf("published port conflict between %s and %s", name, previous)
			}
			ports[key] = name
		}
	}
	existing, err := r.docker.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return fmt.Errorf("validate published ports: %w", err)
	}
	for _, current := range existing {
		if current.Labels["com.docker.compose.project"] == stackProject(id) {
			continue
		}
		for _, port := range current.Ports {
			key := strconv.Itoa(int(port.PublicPort)) + "/" + port.Type
			if service := ports[key]; service != "" {
				return fmt.Errorf("service %s: port %s is owned by another container", service, key)
			}
		}
	}
	return nil
}
