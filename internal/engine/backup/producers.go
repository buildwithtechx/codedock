package backup

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"codedock/internal/models"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
)

func (bm *BackupManager) execInContainer(ctx context.Context, containerName string, command []string) (string, error) {
	if bm.dockerClient == nil {
		return "", fmt.Errorf("quiesce requires a Docker client")
	}
	inspected, err := bm.dockerClient.ContainerInspect(ctx, containerName)
	if err != nil {
		return "", err
	}
	if inspected.State == nil || !inspected.State.Running {
		return "", fmt.Errorf("container %s is not running", containerName)
	}
	created, err := bm.dockerClient.ContainerExecCreate(ctx, inspected.ID, container.ExecOptions{AttachStdout: true, AttachStderr: true, Cmd: command})
	if err != nil {
		return "", err
	}
	attached, err := bm.dockerClient.ContainerExecAttach(ctx, created.ID, container.ExecAttachOptions{})
	if err != nil {
		return "", err
	}
	defer attached.Close()
	stop := context.AfterFunc(ctx, attached.Close)
	defer stop()
	var stdout, stderr bytes.Buffer
	if _, err := stdcopy.StdCopy(&stdout, &stderr, attached.Reader); err != nil {
		return "", err
	}
	if err := bm.waitExecSuccess(ctx, created.ID, "quiesce"); err != nil {
		return "", fmt.Errorf("%w; stderr: %s", err, strings.TrimSpace(stderr.String()))
	}
	return stdout.String(), nil
}

func splitShellCommand(command string) []string {
	return []string{"sh", "-ec", command}
}

func (bm *BackupManager) quiesceTarget(ctx context.Context, cfg *models.BackupConfig, containerName string) (func(), error) {
	quiesce := strings.TrimSpace(cfg.QuiesceCommand)
	unquiesce := strings.TrimSpace(cfg.UnquiesceCommand)
	if quiesce == "" && unquiesce == "" {
		return func() {}, nil
	}
	if containerName == "" {
		return nil, fmt.Errorf("quiesce commands need a running producer container")
	}
	if quiesce != "" {
		if _, err := bm.execInContainer(ctx, containerName, splitShellCommand(quiesce)); err != nil {
			return nil, fmt.Errorf("quiesce producer: %w", err)
		}
	}
	release := func() {
		if unquiesce == "" {
			return
		}
		unlock, cancel := context.WithTimeout(context.Background(), 120*1000000000)
		defer cancel()
		_, _ = bm.execInContainer(unlock, containerName, splitShellCommand(unquiesce))
	}
	return release, nil
}

func (bm *BackupManager) executeVolumeBackupIncremental(ctx context.Context, cfg *models.BackupConfig, volumeName string) ([]byte, string, error) {
	if bm.dockerClient == nil {
		return nil, "", fmt.Errorf("volume backup requires a Docker client")
	}
	if bm.volumeOperations != nil {
		release, err := bm.volumeOperations.AcquireVolume(volumeName)
		if err != nil {
			return nil, "", err
		}
		defer release()
	}
	stateDir := filepath.Join(bm.backupDir, "incremental", cfg.ID)
	if err := os.MkdirAll(stateDir, 0o700); err != nil {
		return nil, "", err
	}
	snapshot := filepath.Join(stateDir, "snapshot.snar")
	execCmd := []string{"tar", "--listed-incremental=/snapshot/snapshot.snar", "-czf", "-", "-C", "/volume_data", "."}
	resp, err := bm.dockerClient.ContainerCreate(ctx, &container.Config{Image: "alpine", Cmd: execCmd}, &container.HostConfig{Binds: []string{volumeName + ":/volume_data:ro", stateDir + ":/snapshot:rw"}}, nil, nil, "")
	if err != nil {
		return nil, "", fmt.Errorf("create incremental backup container: %w", err)
	}
	defer func() {
		_ = bm.dockerClient.ContainerRemove(context.Background(), resp.ID, container.RemoveOptions{Force: true})
	}()
	attachResp, err := bm.dockerClient.ContainerAttach(ctx, resp.ID, container.AttachOptions{Stdout: true, Stderr: true, Stream: true})
	if err != nil {
		return nil, "", fmt.Errorf("attach incremental backup container: %w", err)
	}
	defer attachResp.Close()
	stop := context.AfterFunc(ctx, attachResp.Close)
	defer stop()
	if err := bm.dockerClient.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return nil, "", fmt.Errorf("start incremental backup container: %w", err)
	}
	var stdoutBuf, stderrBuf bytes.Buffer
	if _, err := stdcopy.StdCopy(&stdoutBuf, &stderrBuf, attachResp.Reader); err != nil {
		return nil, stderrBuf.String(), fmt.Errorf("read incremental stream: %w", err)
	}
	statusCh, errCh := bm.dockerClient.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return nil, "", fmt.Errorf("wait incremental container: %w", err)
		}
	case status := <-statusCh:
		if status.StatusCode != 0 {
			return nil, "", fmt.Errorf("incremental container exited %d: %s", status.StatusCode, stderrBuf.String())
		}
	}
	_ = snapshot
	return stdoutBuf.Bytes(), stderrBuf.String(), nil
}

func (bm *BackupManager) executeFileBackup(ctx context.Context, sourcePath string) ([]byte, string, error) {	if bm.dockerClient == nil {
		return nil, "", fmt.Errorf("file backup requires a Docker client")
	}
	cleaned := filepath.Clean("/" + strings.Trim(sourcePath, "/"))
	if cleaned == "/" || strings.Contains(cleaned, "..") {
		return nil, "", fmt.Errorf("invalid file source path")
	}
	parent, leaf := filepath.Split(strings.TrimSuffix(cleaned, "/"))
	execCmd := []string{"tar", "-czf", "-", "-C", parent, leaf}
	resp, err := bm.dockerClient.ContainerCreate(ctx, &container.Config{Image: "alpine", Cmd: execCmd}, &container.HostConfig{Binds: []string{parent + ":/source:ro"}}, nil, nil, "")
	if err != nil {
		return nil, "", fmt.Errorf("create file backup container: %w", err)
	}
	defer func() {
		_ = bm.dockerClient.ContainerRemove(context.Background(), resp.ID, container.RemoveOptions{Force: true})
	}()
	attachResp, err := bm.dockerClient.ContainerAttach(ctx, resp.ID, container.AttachOptions{Stdout: true, Stderr: true, Stream: true})
	if err != nil {
		return nil, "", fmt.Errorf("attach file backup container: %w", err)
	}
	defer attachResp.Close()
	stop := context.AfterFunc(ctx, attachResp.Close)
	defer stop()
	if err := bm.dockerClient.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return nil, "", fmt.Errorf("start file backup container: %w", err)
	}
	var stdoutBuf, stderrBuf bytes.Buffer
	if _, err := stdcopy.StdCopy(&stdoutBuf, &stderrBuf, attachResp.Reader); err != nil {
		return nil, stderrBuf.String(), fmt.Errorf("read file backup stream: %w", err)
	}
	statusCh, errCh := bm.dockerClient.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return nil, "", fmt.Errorf("wait file backup container: %w", err)
		}
	case status := <-statusCh:
		if status.StatusCode != 0 {
			return nil, "", fmt.Errorf("file backup exited %d: %s", status.StatusCode, stderrBuf.String())
		}
	}
	return stdoutBuf.Bytes(), stderrBuf.String(), nil
}

func (bm *BackupManager) serviceContainer(ctx context.Context, serviceID string) (string, error) {	finder, ok := bm.store.(interface {
		GetAppService(string) (*models.AppService, error)
	})
	if !ok {
		return "", fmt.Errorf("service container lookup unavailable")
	}
	app, err := finder.GetAppService(serviceID)
	if err != nil || app == nil || app.ContainerID == "" {
		return "", fmt.Errorf("service has no running container")
	}
	return app.ContainerID, nil
}

func (bm *BackupManager) executeCustomBackup(ctx context.Context, cfg *models.BackupConfig) ([]byte, string, error) {
	if strings.TrimSpace(cfg.CustomBackupCommand) == "" {
		return nil, "", fmt.Errorf("custom backup command is empty")
	}
	containerName, err := bm.serviceContainer(ctx, cfg.ServiceID)
	if err != nil {
		return nil, "", err
	}
	return bm.executeDump(ctx, containerName, splitShellCommand(cfg.CustomBackupCommand), cfg.Name)
}

func (bm *BackupManager) restoreFile(ctx context.Context, sourcePath string, archive []byte) error {
	if bm.dockerClient == nil {
		return fmt.Errorf("file restore requires a Docker client")
	}
	cleaned := filepath.Clean("/" + strings.Trim(sourcePath, "/"))
	if cleaned == "/" || strings.Contains(cleaned, "..") {
		return fmt.Errorf("invalid file source path")
	}
	parent, leaf := filepath.Split(strings.TrimSuffix(cleaned, "/"))
	if leaf == "" {
		return fmt.Errorf("invalid file source path")
	}
	resp, err := bm.dockerClient.ContainerCreate(ctx, &container.Config{Image: "alpine", Cmd: []string{"tar", "-xzf", "-", "-C", "/target"}}, &container.HostConfig{Binds: []string{parent + ":/target:rw"}}, nil, nil, "")
	if err != nil {
		return fmt.Errorf("create file restore container: %w", err)
	}
	defer func() {
		_ = bm.dockerClient.ContainerRemove(context.Background(), resp.ID, container.RemoveOptions{Force: true})
	}()
	attached, err := bm.dockerClient.ContainerAttach(ctx, resp.ID, container.AttachOptions{Stdin: true, Stdout: true, Stderr: true, Stream: true})
	if err != nil {
		return fmt.Errorf("attach file restore container: %w", err)
	}
	defer attached.Close()
	stop := context.AfterFunc(ctx, attached.Close)
	defer stop()
	if err := bm.dockerClient.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("start file restore container: %w", err)
	}
	writeDone := make(chan error, 1)
	go func() {
		if _, err := io.Copy(attached.Conn, bytes.NewReader(archive)); err != nil {
			writeDone <- err
			return
		}
		writeDone <- attached.CloseWrite()
	}()
	var stdoutBuf, stderrBuf bytes.Buffer
	if _, err := stdcopy.StdCopy(&stdoutBuf, &stderrBuf, attached.Reader); err != nil {
		return fmt.Errorf("read file restore stream: %w", err)
	}
	statusCh, errCh := bm.dockerClient.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("wait file restore container: %w", err)
		}
	case status := <-statusCh:
		if status.StatusCode != 0 {
			return fmt.Errorf("file restore exited %d: %s", status.StatusCode, stderrBuf.String())
		}
	}
	select {
	case err := <-writeDone:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
