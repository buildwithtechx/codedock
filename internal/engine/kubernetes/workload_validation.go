package kubernetes

import (
	"codedock.run/codedock/internal/models"
	"fmt"
	"github.com/google/uuid"
	"net"
	"net/url"
	"path"
	"regexp"
	"strings"
)

func WorkloadIdentity(app *models.AppService) (string, string, error) {
	if _, err := uuid.Parse(app.ID); err != nil {
		return "", "", fmt.Errorf("invalid service identity")
	}
	if _, err := uuid.Parse(app.EnvironmentID); err != nil {
		return "", "", fmt.Errorf("invalid environment identity")
	}
	return "codedock-" + app.EnvironmentID, "app-" + app.ID, nil
}
func ValidateWorkload(workload *models.KubernetesWorkload) error {
	app, target := &workload.App, &workload.Target
	if _, _, err := WorkloadIdentity(app); err != nil {
		return err
	}
	if target.Kind != "kubernetes" || target.ClusterID == "" {
		return fmt.Errorf("select a Kubernetes cluster")
	}
	if workload.Image == "" || strings.ContainsAny(workload.Image, " \t\r\n") {
		return fmt.Errorf("select a deployable image reference")
	}
	if app.Replicas < 1 || app.Replicas > 64 {
		return fmt.Errorf("Kubernetes replica count must be between one and sixty-four")
	}
	if app.RuntimeMode != models.RuntimeModeWorker && (app.InternalPort < 1 || app.InternalPort > 65535) {
		return fmt.Errorf("select a valid internal port")
	}
	if app.RuntimeMode == models.RuntimeModeWorker && (app.HealthCheckPath != "" || app.Domain != "") {
		return fmt.Errorf("workers cannot expose HTTP routing or readiness paths")
	}
	if app.CPULimit < 0 || app.MemoryLimit < 0 {
		return fmt.Errorf("resource limits cannot be negative")
	}
	if app.HealthCheckPath != "" && (!strings.HasPrefix(app.HealthCheckPath, "/") || strings.ContainsAny(app.HealthCheckPath, "\r\n")) {
		return fmt.Errorf("readiness path must be an HTTP path")
	}
	if len(app.Volumes) > 0 {
		return fmt.Errorf("Docker host paths cannot be mapped to cluster storage; configure persistent claims explicitly")
	}
	if app.Domain != "" {
		host := app.Domain
		if strings.Contains(host, "://") {
			parsed, err := url.Parse(host)
			if err != nil || parsed.Path != "" && parsed.Path != "/" {
				return fmt.Errorf("invalid routing domain")
			}
			host = parsed.Hostname()
		}
		if net.ParseIP(host) != nil || !regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9.-]*[a-zA-Z0-9])?$`).MatchString(host) || !strings.Contains(host, ".") {
			return fmt.Errorf("routing requires a DNS hostname")
		}
	}
	names, mounts := map[string]bool{}, map[string]bool{}
	for _, volume := range target.Volumes {
		if !regexp.MustCompile(`^[a-z][a-z0-9-]{0,30}$`).MatchString(volume.Name) || names[volume.Name] {
			return fmt.Errorf("persistent volume names must be unique DNS labels")
		}
		if !path.IsAbs(volume.MountPath) || path.Clean(volume.MountPath) != volume.MountPath || volume.MountPath == "/" || mounts[volume.MountPath] {
			return fmt.Errorf("persistent volume mounts must be unique absolute paths below root")
		}
		if volume.SizeGiB < 1 || volume.SizeGiB > 16384 {
			return fmt.Errorf("persistent storage size must be between 1 and 16384 GiB")
		}
		if volume.StorageClass != "" && !regexp.MustCompile(`^[a-z0-9][a-z0-9.-]{0,252}$`).MatchString(volume.StorageClass) {
			return fmt.Errorf("invalid storage class")
		}
		if app.Replicas > 1 && !volume.Shared {
			return fmt.Errorf("multiple application replicas require shared ReadWriteMany storage; database replication is configured separately")
		}
		names[volume.Name], mounts[volume.MountPath] = true, true
	}
	for key := range workload.Variables {
		if !regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`).MatchString(key) {
			return fmt.Errorf("invalid environment variable name")
		}
	}
	return nil
}
