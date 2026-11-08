package deployments

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type ArchiveService struct {
	deploymentService *DeploymentService
}

func NewArchiveService(ds *DeploymentService) *ArchiveService {
	return &ArchiveService{
		deploymentService: ds,
	}
}

type ArchiveDeployResult struct {
	DeploymentID string `json:"deploymentId"`
	ServiceID    string `json:"serviceId"`
	ProjectID    string `json:"projectId"`
}

func (s *ArchiveService) Deploy(ctx context.Context, serviceID, projectID, archivePath string) (*ArchiveDeployResult, error) {
	srcDir, err := ExtractToTemp(archivePath)
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(filepath.Dir(srcDir))

	depID, err := s.deploymentService.TriggerDeployment(ctx, serviceID, "", projectID, "archive", "Archive Deploy", "Deploying from uploaded archive...\n")
	if err != nil {
		return nil, fmt.Errorf("deployment failed: %w", err)
	}

	return &ArchiveDeployResult{
		DeploymentID: depID,
		ServiceID:    serviceID,
		ProjectID:    projectID,
	}, nil
}

func ExtractToTemp(archivePath string) (string, error) {
	tmpDir := filepath.Join(os.TempDir(), "codedock-archive", uuid.New().String())
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	src, err := os.Open(archivePath)
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", fmt.Errorf("failed to open archive: %w", err)
	}
	defer src.Close()

	if err := UntarStream(tmpDir, src); err != nil {
		os.RemoveAll(tmpDir)
		return "", fmt.Errorf("failed to extract: %w", err)
	}

	return FindSourceDir(tmpDir), nil
}

func UntarStream(dest string, src io.Reader) error {
	gzr, err := gzip.NewReader(src)
	if err != nil {
		return untarReader(dest, tar.NewReader(src))
	}
	defer gzr.Close()
	return untarReader(dest, tar.NewReader(gzr))
}

func untarReader(dest string, tr *tar.Reader) error {
	destClean := filepath.Clean(dest)
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}

		path := filepath.Join(dest, filepath.Clean(hdr.Name))
		if !strings.HasPrefix(path, destClean+string(os.PathSeparator)) {
			continue
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(path, 0o755)
		case tar.TypeReg:
			os.MkdirAll(filepath.Dir(path), 0o755)
			f, _ := os.Create(path)
			if f != nil {
				_, _ = io.Copy(f, tr)
				f.Close()
			}
		}
	}
	return nil
}

func FindSourceDir(dir string) string {
	entries, _ := os.ReadDir(dir)
	if len(entries) == 1 && entries[0].IsDir() {
		return filepath.Join(dir, entries[0].Name())
	}
	return dir
}
