package build

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/docker/docker/api/types/volume"
)

type mockDockerVolumeClient struct {
	volumes map[string]volume.Volume
}

func newMockDockerVolumeClient() *mockDockerVolumeClient {
	return &mockDockerVolumeClient{
		volumes: make(map[string]volume.Volume),
	}
}

func (m *mockDockerVolumeClient) VolumeCreate(ctx context.Context, options volume.CreateOptions) (volume.Volume, error) {
	vol := volume.Volume{
		Name:      options.Name,
		Driver:    options.Driver,
		Labels:    options.Labels,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	m.volumes[options.Name] = vol
	return vol, nil
}

func (m *mockDockerVolumeClient) VolumeInspect(ctx context.Context, volumeID string) (volume.Volume, error) {
	vol, ok := m.volumes[volumeID]
	if !ok {
		return volume.Volume{}, fmt.Errorf("volume not found: %s", volumeID)
	}
	return vol, nil
}

func (m *mockDockerVolumeClient) VolumeRemove(ctx context.Context, volumeID string, force bool) error {
	if _, ok := m.volumes[volumeID]; !ok {
		return fmt.Errorf("volume not found: %s", volumeID)
	}
	delete(m.volumes, volumeID)
	return nil
}

func (m *mockDockerVolumeClient) VolumeList(ctx context.Context, options volume.ListOptions) (volume.ListResponse, error) {
	var list []*volume.Volume
	for _, v := range m.volumes {
		vol := v
		list = append(list, &vol)
	}
	return volume.ListResponse{Volumes: list}, nil
}

func TestEnsureProjectBuildCache(t *testing.T) {
	client := newMockDockerVolumeClient()
	ctx := context.Background()

	t.Run("creates new cache volume", func(t *testing.T) {
		reused, err := EnsureProjectBuildCache(ctx, client, "p1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if reused {
			t.Fatal("expected reused=false on initial creation")
		}
	})

	t.Run("reuses existing cache volume", func(t *testing.T) {
		reused, err := EnsureProjectBuildCache(ctx, client, "p1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reused {
			t.Fatal("expected reused=true on second call")
		}
	})

	t.Run("owner mismatch error", func(t *testing.T) {
		client.volumes[ProjectBuildCacheVolume("p2")] = volume.Volume{
			Name: ProjectBuildCacheVolume("p2"),
			Labels: map[string]string{
				BuildCacheLabel:        "1",
				BuildCacheProjectLabel: "wrong-owner",
			},
		}
		_, err := EnsureProjectBuildCache(ctx, client, "p2")
		if err != ErrBuildCacheOwnerMismatch {
			t.Fatalf("expected ErrBuildCacheOwnerMismatch, got %v", err)
		}
	})
}

func TestPruneManagedBuildCaches(t *testing.T) {
	client := newMockDockerVolumeClient()
	ctx := context.Background()
	now := time.Now()

	// Cache 1: 10 days old, idle (size 2 GiB)
	vol1Name := ProjectBuildCacheVolume("p1")
	client.volumes[vol1Name] = volume.Volume{
		Name:      vol1Name,
		Driver:    "local",
		CreatedAt: now.Add(-10 * 24 * time.Hour).Format(time.RFC3339),
		Labels: map[string]string{
			BuildCacheLabel:        "1",
			BuildCacheProjectLabel: "p1",
		},
		UsageData: &volume.UsageData{Size: 2 * 1024 * 1024 * 1024, RefCount: 0},
	}

	// Cache 2: 1 hour old, active reference (size 1 GiB)
	vol2Name := ProjectBuildCacheVolume("p2")
	client.volumes[vol2Name] = volume.Volume{
		Name:      vol2Name,
		Driver:    "local",
		CreatedAt: now.Add(-1 * time.Hour).Format(time.RFC3339),
		Labels: map[string]string{
			BuildCacheLabel:        "1",
			BuildCacheProjectLabel: "p2",
		},
		UsageData: &volume.UsageData{Size: 1 * 1024 * 1024 * 1024, RefCount: 1},
	}

	res, err := PruneManagedBuildCaches(ctx, client, PruneOptions{
		Now:       now,
		KeepBytes: 10 * 1024 * 1024 * 1024,
	})
	if err != nil {
		t.Fatalf("unexpected prune error: %v", err)
	}

	if len(res.CachesDeleted) != 1 || res.CachesDeleted[0] != vol1Name {
		t.Fatalf("expected vol1 deleted due to age expiry, got %v", res.CachesDeleted)
	}
	if res.SpaceReclaimed != 2*1024*1024*1024 {
		t.Fatalf("unexpected space reclaimed: %d", res.SpaceReclaimed)
	}

	// Verify active reference volume p2 was kept
	if _, ok := client.volumes[vol2Name]; !ok {
		t.Fatal("expected active referenced volume p2 to be kept")
	}
}
