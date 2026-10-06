package backup

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"

	"codedock.run/codedock/internal/models"

	"codedock.run/codedock/internal/engine/compose"
)

func (bm *BackupManager) RestoreBackup(ctx context.Context, recordID string) error {
	rec, err := bm.store.GetBackupRecord(recordID)
	if err != nil || rec == nil {
		return fmt.Errorf("backup record not found: %w", err)
	}
	if rec.Status != "completed" {
		return errors.New("cannot restore: backup is not completed")
	}

	cfg, err := bm.store.GetBackupConfig(rec.BackupConfigID)
	if err != nil || cfg == nil {
		return fmt.Errorf("backup config not found: %w", err)
	}

	containerName, restoreCmd, err := bm.buildRestoreCommand(cfg)
	if err != nil {
		return fmt.Errorf("failed to build restore command: %w", err)
	}

	if bm.dockerClient == nil {
		return errors.New("database restore requires a Docker client")
	}
	archive, _, err := bm.OpenBackupArchive(ctx, recordID)
	if err != nil {
		return fmt.Errorf("open restore archive: %w", err)
	}
	data, readErr := io.ReadAll(archive)
	closeErr := archive.Close()
	if readErr != nil {
		return fmt.Errorf("read restore archive: %w", readErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close restore archive: %w", closeErr)
	}
	if len(data) == 0 {
		return errors.New("restore archive is empty")
	}

	return bm.executeRestore(ctx, containerName, restoreCmd, data)
}

func (bm *BackupManager) buildRestoreCommand(cfg *models.BackupConfig) (string, []string, error) {
	if cfg.DatabaseID != "" {
		db, err := bm.store.GetDatabase(cfg.DatabaseID)
		if err != nil || db == nil {
			return "", nil, fmt.Errorf("target database %s not found", cfg.DatabaseID)
		}
		containerName := databaseContainerIdentity(db)
		tmplMgr, err := compose.NewTemplateManager()
		if err != nil {
			return "", nil, fmt.Errorf("failed to init template manager: %v", err)
		}

		composeFile, err := tmplMgr.GetTemplate(strings.ToLower(string(db.Engine)))
		if err != nil {
			return "", nil, fmt.Errorf("unsupported database engine %s: %v", db.Engine, err)
		}

		tmplService, exists := composeFile.Services[strings.ToLower(string(db.Engine))]
		if !exists {
			for _, s := range composeFile.Services {
				tmplService = s
				break
			}
		}

		if tmplService.XCodedock != nil && tmplService.XCodedock.Restore != nil && len(tmplService.XCodedock.Restore.Command) > 0 {
			var cmd []string
			for _, c := range tmplService.XCodedock.Restore.Command {
				resolved := strings.ReplaceAll(c, "${db.password}", db.Password)
				resolved = strings.ReplaceAll(resolved, "${db.username}", db.Username)
				resolved = strings.ReplaceAll(resolved, "${db.database_name}", db.DatabaseName)
				cmd = append(cmd, resolved)
			}
			return containerName, cmd, nil
		}
		return "", nil, fmt.Errorf("database engine %s has no restore command", db.Engine)

	}

	return "", nil, errors.New("backup config requires databaseId")
}

func (bm *BackupManager) executeRestore(ctx context.Context, containerName string, restoreCmd []string, data []byte) error {
	if bm.dockerClient == nil {
		return errors.New("database restore requires a Docker client")
	}

	inspectResp, err := bm.dockerClient.ContainerInspect(ctx, containerName)
	if err != nil {
		return fmt.Errorf("inspect restore container: %w", err)
	}
	if inspectResp.State == nil || !inspectResp.State.Running {
		return fmt.Errorf("cannot restore: container %s is stopped or not running", containerName)
	}

	execConfig := container.ExecOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          restoreCmd,
	}

	execCreateResp, err := bm.dockerClient.ContainerExecCreate(ctx, inspectResp.ID, execConfig)
	if err != nil {
		return fmt.Errorf("docker exec create failed: %v", err)
	}

	attachResp, err := bm.dockerClient.ContainerExecAttach(ctx, execCreateResp.ID, container.ExecAttachOptions{})
	if err != nil {
		return fmt.Errorf("docker exec attach failed: %v", err)
	}
	defer attachResp.Close()
	stopCancellation := context.AfterFunc(ctx, attachResp.Close)
	defer stopCancellation()

	writeDone := make(chan error, 1)
	go func() {
		if _, err := io.Copy(attachResp.Conn, bytes.NewReader(data)); err != nil {
			attachResp.Close()
			writeDone <- fmt.Errorf("write restore archive: %w", err)
			return
		}
		writeDone <- attachResp.CloseWrite()
	}()

	var stdoutBuf, stderrBuf bytes.Buffer
	_, readErr := stdcopy.StdCopy(&stdoutBuf, &stderrBuf, attachResp.Reader)
	attachResp.Close()
	if readErr != nil {
		return fmt.Errorf("read restore stream: %w; stderr: %s", readErr, strings.TrimSpace(stderrBuf.String()))
	}
	if err := bm.waitExecSuccess(ctx, execCreateResp.ID, "restore"); err != nil {
		return fmt.Errorf("%w; stderr: %s", err, strings.TrimSpace(stderrBuf.String()))
	}
	select {
	case err := <-writeDone:
		if err != nil {
			return fmt.Errorf("send restore archive: %w; stderr: %s", err, strings.TrimSpace(stderrBuf.String()))
		}
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}
