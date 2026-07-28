package backup

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"

	"codedock.run/codedock/internal/engine/compose"
)

func (bm *BackupManager) executeVolumeBackup(ctx context.Context, volumeName string) ([]byte, string, error) {
	if bm.dockerClient == nil {
		dumpBytes := []byte(fmt.Sprintf("-- Simulated volume backup dump for %s at %s --\n", volumeName, time.Now().UTC().Format(time.RFC3339)))
		return dumpBytes, "Docker client nil: simulated successful local dump.\n", nil
	}

	execCmd := []string{"tar", "-czf", "-", "-C", "/volume_data", "."}

	resp, err := bm.dockerClient.ContainerCreate(ctx, &container.Config{
		Image: "alpine",
		Cmd:   execCmd,
	}, &container.HostConfig{
		Binds: []string{fmt.Sprintf("%s:/volume_data:ro", volumeName)},
	}, nil, nil, "")
	if err != nil {
		return nil, "", fmt.Errorf("failed to create backup container for volume %s: %w", volumeName, err)
	}

	defer func() {
		_ = bm.dockerClient.ContainerRemove(context.Background(), resp.ID, container.RemoveOptions{Force: true})
	}()

	attachResp, err := bm.dockerClient.ContainerAttach(ctx, resp.ID, container.AttachOptions{
		Stdout: true,
		Stderr: true,
		Stream: true,
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to attach to backup container: %w", err)
	}
	defer attachResp.Close()

	if err := bm.dockerClient.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return nil, "", fmt.Errorf("failed to start backup container: %w", err)
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	if _, err := stdcopy.StdCopy(&stdoutBuf, &stderrBuf, attachResp.Reader); err != nil {
		_, _ = io.Copy(&stdoutBuf, attachResp.Reader)
	}

	statusCh, errCh := bm.dockerClient.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return nil, "", fmt.Errorf("error waiting for backup container: %w", err)
		}
	case status := <-statusCh:
		if status.StatusCode != 0 {
			return nil, "", fmt.Errorf("backup container exited with status %d, stderr: %s", status.StatusCode, stderrBuf.String())
		}
	}

	return stdoutBuf.Bytes(), stderrBuf.String(), nil
}

func (bm *BackupManager) buildDumpCommand(cfg *models.BackupConfig) (string, []string, string, error) {
	if cfg.DatabaseID != "" {
		db, err := bm.store.GetDatabase(cfg.DatabaseID)
		if err != nil || db == nil {
			return "", nil, "", fmt.Errorf("target database %s not found", cfg.DatabaseID)
		}
		containerName := utils.NormalizeContainerName(db.ID)
		tmplMgr, err := compose.NewTemplateManager()
		if err != nil {
			return "", nil, "", fmt.Errorf("failed to init template manager: %v", err)
		}

		composeFile, err := tmplMgr.GetTemplate(strings.ToLower(string(db.Engine)))
		if err != nil {
			return "", nil, "", fmt.Errorf("unsupported database engine %s: %v", db.Engine, err)
		}

		tmplService, exists := composeFile.Services[strings.ToLower(string(db.Engine))]
		if !exists {
			for _, s := range composeFile.Services {
				tmplService = s
				break
			}
		}

		if tmplService.XCodedock != nil && tmplService.XCodedock.Backup != nil && len(tmplService.XCodedock.Backup.Command) > 0 {
			var dumpCmd []string
			for _, c := range tmplService.XCodedock.Backup.Command {
				resolved := strings.ReplaceAll(c, "${db.password}", db.Password)
				resolved = strings.ReplaceAll(resolved, "${db.username}", db.Username)
				resolved = strings.ReplaceAll(resolved, "${db.database_name}", db.DatabaseName)
				dumpCmd = append(dumpCmd, resolved)
			}
			return containerName, dumpCmd, tmplService.XCodedock.Backup.FileExtension, nil
		}
		return containerName, []string{"sh", "-c", "echo 'Generic volume snapshot'"}, ".tar.gz", nil
	}

	return "", nil, "", errors.New("backup config requires databaseId")
}

func (bm *BackupManager) executeDump(ctx context.Context, containerName string, dumpCmd []string, backupName string) ([]byte, string, error) {
	if bm.dockerClient == nil {
		dumpBytes := []byte(fmt.Sprintf("-- Simulated backup dump for %s at %s --\n", backupName, time.Now().UTC().Format(time.RFC3339)))
		return dumpBytes, "Docker client nil: simulated successful local dump.\n", nil
	}

	inspectResp, err := bm.dockerClient.ContainerInspect(ctx, containerName)
	if err != nil || !inspectResp.State.Running {
		return nil, "", fmt.Errorf("cannot backup: container %s is stopped or not running", containerName)
	}

	execConfig := container.ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          dumpCmd,
	}

	execCreateResp, err := bm.dockerClient.ContainerExecCreate(ctx, inspectResp.ID, execConfig)
	if err != nil {
		return nil, "", fmt.Errorf("docker exec create failed: %v", err)
	}

	attachResp, err := bm.dockerClient.ContainerExecAttach(ctx, execCreateResp.ID, container.ExecAttachOptions{})
	if err != nil {
		return nil, "", fmt.Errorf("docker exec attach failed: %v", err)
	}
	defer attachResp.Close()

	var stdoutBuf, stderrBuf bytes.Buffer
	if _, err := stdcopy.StdCopy(&stdoutBuf, &stderrBuf, attachResp.Reader); err != nil {
		_, _ = io.Copy(&stdoutBuf, attachResp.Reader)
	}

	return stdoutBuf.Bytes(), stderrBuf.String(), nil
}

func (bm *BackupManager) handleS3Upload(ctx context.Context, cfg *models.BackupConfig, fileName string, dumpBytes []byte, execLogs string) (string, string, error) {
	if cfg.S3DestinationID == "" {
		return "", execLogs, fmt.Errorf("S3 destination ID missing for config %s", cfg.ID)
	}
	dest, err := bm.store.GetS3Destination(cfg.S3DestinationID)
	if err != nil {
		return "", execLogs, fmt.Errorf("failed to get S3 destination %s: %w", cfg.S3DestinationID, err)
	}
	if dest == nil {
		return "", execLogs, fmt.Errorf("S3 destination %s not found", cfg.S3DestinationID)
	}
	s3URL, err := bm.uploadToS3(ctx, dest, fileName, dumpBytes)
	if err != nil {
		return "", execLogs, fmt.Errorf("S3 upload failed: %w", err)
	}
	execLogs += fmt.Sprintf("\n✅ Successfully uploaded backup to S3/MinIO destination: %s", s3URL)
	return s3URL, execLogs, nil
}
