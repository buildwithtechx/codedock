package backup

import (
	"codedock.run/codedock/internal/models"
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"io"
	"log/slog"
	"strconv"
	"time"
)

type VolumeRestoreOwner interface {
	VolumeRestoreOwner(context.Context, *models.BackupConfig) (string, error)
}

func (bm *BackupManager) ValidateVolumeRestore(ctx context.Context, recordID string) (*models.VolumeRestoreTarget, error) {
	if bm.dockerClient == nil {
		return nil, errors.New("volume restore requires Docker")
	}
	rec, err := bm.store.GetBackupRecord(recordID)
	if err != nil {
		return nil, fmt.Errorf("load volume backup record: %w", err)
	}
	if rec == nil || rec.Status != models.BackupRecordStatusCompleted {
		return nil, errors.New("volume restore requires a completed backup")
	}
	cfg, err := bm.store.GetBackupConfig(rec.BackupConfigID)
	if err != nil {
		return nil, fmt.Errorf("load volume backup config: %w", err)
	}
	if cfg == nil || cfg.VolumeName == "" {
		return nil, errors.New("backup is not a named volume archive")
	}
	owner, ok := bm.store.(VolumeRestoreOwner)
	if !ok {
		return nil, errors.New("volume ownership validation unavailable")
	}
	containerID, err := owner.VolumeRestoreOwner(ctx, cfg)
	if err != nil {
		return nil, err
	}
	volume, err := bm.dockerClient.VolumeInspect(ctx, cfg.VolumeName)
	if err != nil {
		return nil, fmt.Errorf("target volume must already exist: %w", err)
	}
	if volume.Name != cfg.VolumeName {
		return nil, errors.New("target volume identity mismatch")
	}
	inspected, err := bm.dockerClient.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, fmt.Errorf("inspect volume owner: %w", err)
	}
	bound := false
	for _, m := range inspected.Mounts {
		if m.Type == mount.TypeVolume && m.Name == cfg.VolumeName {
			bound = true
		}
	}
	if !bound {
		return nil, errors.New("named volume is not mounted by its registered owner; host bind mounts are unsupported")
	}
	consumers, err := bm.dockerClient.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("check active volume consumers: %w", err)
	}
	for _, consumer := range consumers {
		for _, m := range consumer.Mounts {
			if m.Type == mount.TypeVolume && m.Name == cfg.VolumeName {
				return nil, errors.New("stop every container using the volume before restoring")
			}
		}
	}
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 3600
	}
	if timeout > 86400 {
		timeout = 86400
	}
	return &models.VolumeRestoreTarget{RecordID: recordID, VolumeName: cfg.VolumeName, ContainerID: containerID, VolumeCreatedAt: volume.CreatedAt, TimeoutSeconds: timeout}, nil
}

func (bm *BackupManager) CancelVolumeRestore(recordID string) bool {
	bm.restoreMu.Lock()
	defer bm.restoreMu.Unlock()
	cancel, ok := bm.restores[recordID]
	if ok {
		cancel()
	}
	return ok
}

func (bm *BackupManager) RestoreVolume(ctx context.Context, recordID, confirmedVolume string) error {
	target, err := bm.ValidateVolumeRestore(ctx, recordID)
	if err != nil {
		return err
	}
	if target.VolumeName != confirmedVolume {
		return errors.New("confirmation must match the target volume name")
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(target.TimeoutSeconds)*time.Second)
	defer cancel()
	bm.restoreMu.Lock()
	if _, exists := bm.restoreVolumes[target.VolumeName]; exists {
		bm.restoreMu.Unlock()
		return errors.New("a restore is already running for this volume")
	}
	bm.restores[recordID] = cancel
	bm.restoreVolumes[target.VolumeName] = recordID
	bm.restoreMu.Unlock()
	defer func() {
		bm.restoreMu.Lock()
		delete(bm.restores, recordID)
		delete(bm.restoreVolumes, target.VolumeName)
		bm.restoreMu.Unlock()
	}()
	archive, _, err := bm.OpenBackupArchive(ctx, recordID)
	if err != nil {
		return err
	}
	defer func() {
		if err := archive.Close(); err != nil {
			slog.Warn("close volume restore archive", "error", err)
		}
	}()
	result, err := bm.dockerClient.ContainerCreate(ctx, &container.Config{
		Image: "alpine", OpenStdin: true, StdinOnce: true, AttachStdin: true, AttachStdout: true, AttachStderr: true,
		Cmd: []string{"timeout", strconv.Itoa(target.TimeoutSeconds), "tar", "-xzf", "-", "-C", "/restore"},
	}, &container.HostConfig{NetworkMode: "none", Mounts: []mount.Mount{{Type: mount.TypeVolume, Source: target.VolumeName, Target: "/restore"}}}, nil, nil, "")
	if err != nil {
		return fmt.Errorf("create volume restore helper: %w", err)
	}
	defer func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := bm.dockerClient.ContainerRemove(cleanup, result.ID, container.RemoveOptions{Force: true}); err != nil && !client.IsErrNotFound(err) {
			slog.Warn("remove volume restore helper", "container_id", result.ID, "error", err)
		}
	}()
	stop := context.AfterFunc(ctx, func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := bm.dockerClient.ContainerRemove(cleanup, result.ID, container.RemoveOptions{Force: true}); err != nil && !client.IsErrNotFound(err) {
			slog.Warn("remove volume restore helper", "container_id", result.ID, "error", err)
		}
	})
	defer stop()
	attached, err := bm.dockerClient.ContainerAttach(ctx, result.ID, container.AttachOptions{Stdin: true, Stdout: true, Stderr: true, Stream: true})
	if err != nil {
		return fmt.Errorf("attach volume restore input: %w", err)
	}
	defer attached.Close()
	closeOnCancel := context.AfterFunc(ctx, attached.Close)
	defer closeOnCancel()
	validated, err := bm.ValidateVolumeRestore(ctx, recordID)
	if err != nil {
		return fmt.Errorf("revalidate volume restore target: %w", err)
	}
	if validated.VolumeName != target.VolumeName || validated.VolumeCreatedAt != target.VolumeCreatedAt {
		return errors.New("volume restore target changed before startup")
	}
	if err := bm.dockerClient.ContainerStart(ctx, result.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("start volume restore: %w", err)
	}
	go func() {
		if _, err := io.Copy(io.Discard, attached.Reader); err != nil && ctx.Err() == nil {
			slog.Warn("read volume restore output", "error", err)
		}
	}()
	sent := make(chan error, 1)
	go func() {
		_, copyErr := io.Copy(attached.Conn, archive)
		closeErr := attached.CloseWrite()
		if copyErr != nil {
			sent <- copyErr
		} else {
			sent <- closeErr
		}
	}()
	status, failures := bm.dockerClient.ContainerWait(ctx, result.ID, container.WaitConditionNotRunning)
	select {
	case <-ctx.Done():
		attached.Close()
		return fmt.Errorf("volume restore interrupted; target may contain partial data: %w", ctx.Err())
	case err := <-failures:
		attached.Close()
		if err != nil {
			return fmt.Errorf("wait for volume restore: %w", err)
		}
		return errors.New("volume restore ended without an exit result")
	case result := <-status:
		attached.Close()
		if result.StatusCode != 0 {
			return fmt.Errorf("volume restore exited with status %d; target may contain partial data", result.StatusCode)
		}
	}
	select {
	case err := <-sent:
		if err != nil {
			return fmt.Errorf("send volume archive: %w", err)
		}
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}
