package build

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
)

const (
	BuildCacheLabel        = "codedock.build-cache"
	BuildCacheProjectLabel = "codedock.build-cache.project"
	DefaultBuildCacheAge   = 7 * 24 * time.Hour
)

var (
	ErrBuildCacheOwnerMismatch = errors.New("the build cache volume belongs to another owner")
	ErrBuildCacheInUse         = errors.New("the project's build cache is still in use")
)

type DockerVolumeClient interface {
	VolumeCreate(ctx context.Context, options volume.CreateOptions) (volume.Volume, error)
	VolumeInspect(ctx context.Context, volumeID string) (volume.Volume, error)
	VolumeRemove(ctx context.Context, volumeID string, force bool) error
	VolumeList(ctx context.Context, options volume.ListOptions) (volume.ListResponse, error)
}

func ProjectBuildCacheVolume(projectID string) string {
	return fmt.Sprintf("codedock_buildkit_%s_state", strings.TrimSpace(projectID))
}

func BuildCacheOwner(vol volume.Volume) string {
	if vol.Labels == nil {
		return ""
	}
	if vol.Labels[BuildCacheLabel] != "1" {
		return ""
	}
	projectID := vol.Labels[BuildCacheProjectLabel]
	if projectID == "" {
		return ""
	}
	if vol.Driver != "" && vol.Driver != "local" {
		return ""
	}
	if vol.Name != ProjectBuildCacheVolume(projectID) {
		return ""
	}
	return projectID
}

func EnsureProjectBuildCache(ctx context.Context, docker DockerVolumeClient, projectID string) (bool, error) {
	name := ProjectBuildCacheVolume(projectID)
	vol, err := docker.VolumeInspect(ctx, name)
	if err == nil {
		if BuildCacheOwner(vol) != projectID {
			return false, ErrBuildCacheOwnerMismatch
		}
		return true, nil
	}

	created, err := docker.VolumeCreate(ctx, volume.CreateOptions{
		Name: name,
		Labels: map[string]string{
			BuildCacheLabel:        "1",
			BuildCacheProjectLabel: projectID,
		},
	})
	if err != nil {
		if volInspect, errInspect := docker.VolumeInspect(ctx, name); errInspect == nil {
			if BuildCacheOwner(volInspect) != projectID {
				return false, ErrBuildCacheOwnerMismatch
			}
			return true, nil
		}
		return false, fmt.Errorf("create build cache volume: %w", err)
	}

	if BuildCacheOwner(created) != projectID {
		return false, ErrBuildCacheOwnerMismatch
	}
	return false, nil
}

type ManagedCacheRecord struct {
	Name       string
	ProjectID  string
	Bytes      uint64
	References int
	CreatedAt  time.Time
	UsedAt     time.Time
}

type PruneOptions struct {
	KeepBytes           uint64
	ProjectID           string
	Now                 time.Time
	ProtectedProjectIDs map[string]bool
	ProtectAll          bool
	Storage             *BuildStorage
}

type PruneResult struct {
	CachesDeleted  []string
	SpaceReclaimed uint64
}

func PruneManagedBuildCaches(ctx context.Context, docker DockerVolumeClient, opts PruneOptions) (PruneResult, error) {
	keepBytes := opts.KeepBytes
	if keepBytes == 0 {
		keepBytes = DefaultBuildCacheKeepBytes
	}
	now := opts.Now
	if now.IsZero() {
		now = time.Now()
	}

	list, err := docker.VolumeList(ctx, volume.ListOptions{
		Filters: filters.NewArgs(filters.Arg("label", BuildCacheLabel+"=1")),
	})
	if err != nil {
		return PruneResult{}, fmt.Errorf("list build cache volumes: %w", err)
	}

	var records []ManagedCacheRecord
	var totalBytes uint64

	for _, vol := range list.Volumes {
		if vol == nil {
			continue
		}
		projectID := BuildCacheOwner(*vol)
		if projectID == "" {
			continue
		}
		if opts.ProjectID != "" && projectID != opts.ProjectID {
			continue
		}

		var size uint64
		var refs int
		if vol.UsageData != nil {
			if vol.UsageData.Size > 0 {
				size = uint64(vol.UsageData.Size)
			}
			if vol.UsageData.RefCount > 0 {
				refs = int(vol.UsageData.RefCount)
			}
		}

		createdAt := now
		if parsed, err := time.Parse(time.RFC3339Nano, vol.CreatedAt); err == nil {
			createdAt = parsed
		} else if parsed, err := time.Parse(time.RFC3339, vol.CreatedAt); err == nil {
			createdAt = parsed
		}

		records = append(records, ManagedCacheRecord{
			Name:       vol.Name,
			ProjectID:  projectID,
			Bytes:      size,
			References: refs,
			CreatedAt:  createdAt,
			UsedAt:     createdAt,
		})
		totalBytes += size
	}

	if opts.ProjectID != "" {
		if opts.ProtectAll || (opts.ProtectedProjectIDs != nil && opts.ProtectedProjectIDs[opts.ProjectID]) {
			return PruneResult{}, ErrBuildCacheInUse
		}
		for _, r := range records {
			if r.ProjectID == opts.ProjectID && r.References > 0 {
				return PruneResult{}, ErrBuildCacheInUse
			}
		}
	}

	sort.Slice(records, func(i, j int) bool {
		if records[i].UsedAt.Equal(records[j].UsedAt) {
			return records[i].Name < records[j].Name
		}
		return records[i].UsedAt.Before(records[j].UsedAt)
	})

	var result PruneResult
	storage := opts.Storage

	for _, rec := range records {
		if rec.References > 0 {
			continue
		}
		if opts.ProtectAll || (opts.ProtectedProjectIDs != nil && opts.ProtectedProjectIDs[rec.ProjectID]) {
			continue
		}

		expired := now.Sub(rec.UsedAt) > DefaultBuildCacheAge
		pressure := false
		if storage != nil {
			pressure = storage.FreeBytes < BuildStorageReserve(*storage) || (storage.FreeInodes != nil && *storage.FreeInodes == 0)
		}

		shouldPrune := false
		if opts.ProjectID != "" {
			shouldPrune = rec.ProjectID == opts.ProjectID
		} else {
			shouldPrune = expired || pressure || totalBytes > keepBytes
		}

		if !shouldPrune {
			continue
		}

		inspected, err := docker.VolumeInspect(ctx, rec.Name)
		if err != nil {
			continue
		}
		if BuildCacheOwner(inspected) != rec.ProjectID {
			return result, ErrBuildCacheOwnerMismatch
		}

		if err := docker.VolumeRemove(ctx, rec.Name, false); err != nil {
			if opts.ProjectID != "" {
				return result, ErrBuildCacheInUse
			}
			continue
		}

		result.CachesDeleted = append(result.CachesDeleted, rec.Name)
		result.SpaceReclaimed += rec.Bytes
		if totalBytes >= rec.Bytes {
			totalBytes -= rec.Bytes
		}
		if storage != nil {
			storage.FreeBytes += rec.Bytes
		}
	}

	return result, nil
}
