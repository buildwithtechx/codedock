package deploy

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path"
	"strings"

	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

const ImageVolumeLabelProject = "codedock.project"
const ImageVolumeLabelService = "codedock.service_id"

// ImageVolumeName generates a deterministic volume name for a workload destination.
func ImageVolumeName(workload, destination string) string {
	sum := sha256.Sum256(fmt.Appendf(nil, "%s:%s", workload, destination))
	return fmt.Sprintf("codedock-vol-%s", hex.EncodeToString(sum[:])[:24])
}

// CleanVolumeDestination normalizes POSIX destination path.
func CleanVolumeDestination(target string) (string, error) {
	if !strings.HasPrefix(target, "/") || strings.Contains(target, ":") {
		return "", fmt.Errorf("invalid persistent volume destination: %s", target)
	}
	clean := path.Clean(target)
	if clean == "" {
		clean = "/"
	}
	return clean, nil
}

// ResolveImageVolumeBinds ensures image-declared VOLUME paths are backed by
// persistent named volumes across redeploys instead of ephemeral anonymous volumes.
func ResolveImageVolumeBinds(
	ctx context.Context,
	dockerCli *client.Client,
	serviceID string,
	projectID string,
	declaredBinds []string,
	imageTag string,
	previousContainerID string,
) ([]string, error) {
	if dockerCli == nil || imageTag == "" {
		return declaredBinds, nil
	}

	inspect, _, err := dockerCli.ImageInspectWithRaw(ctx, imageTag)
	if err != nil {
		// If image cannot be inspected locally (e.g. remote mock in unit tests), return original binds
		return declaredBinds, nil
	}

	if inspect.Config == nil || len(inspect.Config.Volumes) == 0 {
		return declaredBinds, nil
	}

	covered := make(map[string]bool)
	for _, b := range declaredBinds {
		parts := strings.Split(b, ":")
		if len(parts) >= 2 {
			if clean, err := CleanVolumeDestination(parts[1]); err == nil {
				covered[clean] = true
			}
		}
	}

	// Check if previous container exists and had existing mounts
	previousMounts := make(map[string]string)
	if previousContainerID != "" {
		prevInspect, err := dockerCli.ContainerInspect(ctx, previousContainerID)
		if err == nil && prevInspect.Config != nil {
			for _, m := range prevInspect.Mounts {
				if m.Type == "volume" && m.Name != "" {
					if clean, err := CleanVolumeDestination(m.Destination); err == nil {
						previousMounts[clean] = m.Name
					}
				}
			}
		}
	}

	finalBinds := append([]string{}, declaredBinds...)

	for rawVol := range inspect.Config.Volumes {
		dest, err := CleanVolumeDestination(rawVol)
		if err != nil || covered[dest] {
			continue
		}

		volName := previousMounts[dest]
		if volName == "" {
			workload := serviceID
			if workload == "" {
				workload = projectID
			}
			volName = ImageVolumeName(workload, dest)
		}

		// Ensure named volume exists with project/service labels
		labels := map[string]string{
			ImageVolumeLabelService: serviceID,
		}
		if projectID != "" {
			labels[ImageVolumeLabelProject] = projectID
		}

		_, inspectErr := dockerCli.VolumeInspect(ctx, volName)
		if inspectErr != nil {
			_, createErr := dockerCli.VolumeCreate(ctx, volume.CreateOptions{
				Name:   volName,
				Labels: labels,
			})
			if createErr != nil && !strings.Contains(createErr.Error(), "already exists") {
				return nil, fmt.Errorf("ensure persistent volume %s: %w", volName, createErr)
			}
		}

		finalBinds = append(finalBinds, fmt.Sprintf("%s:%s", volName, dest))
		covered[dest] = true
	}

	return finalBinds, nil
}
