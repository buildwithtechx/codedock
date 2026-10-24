package deploy

import (
	"codedock/internal/models"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/containerd/errdefs"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"os"
	"path/filepath"
	"strings"
)

func (d *Deployer) SetRolloutDirectory(directory string) error {
	if err := os.MkdirAll(directory, 0700); err != nil {
		return fmt.Errorf("prepare rollout recovery: %w", err)
	}
	d.rolloutDirectory = directory
	return nil
}

func (d *Deployer) rolloutPath(service string) string {
	digest := sha256.Sum256([]byte(service))
	return filepath.Join(d.rolloutDirectory, hex.EncodeToString(digest[:])+".json")
}

func (d *Deployer) saveRollout(journal *models.RolloutJournal) error {
	if d.rolloutDirectory == "" {
		return nil
	}
	data, err := json.Marshal(journal)
	if err != nil {
		return err
	}
	file, err := os.CreateTemp(d.rolloutDirectory, ".rollout-")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	if _, err := file.Write(data); err != nil {
		file.Close()
		return err
	}
	if err := file.Sync(); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(file.Name(), d.rolloutPath(journal.PreviousApp.ID))
}

func (d *Deployer) clearRollout(service string) error {
	if d.rolloutDirectory == "" {
		return nil
	}
	if err := os.Remove(d.rolloutPath(service)); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func (d *Deployer) RecoverRollouts(ctx context.Context) error {
	if d.rolloutDirectory == "" || d.containerManager.dockerClient == nil {
		return nil
	}
	files, err := os.ReadDir(d.rolloutDirectory)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(d.rolloutDirectory, file.Name()))
		if err != nil {
			return err
		}
		var journal models.RolloutJournal
		if err := json.Unmarshal(data, &journal); err != nil {
			return fmt.Errorf("read rollout journal: %w", err)
		}
		if filepath.Base(d.rolloutPath(journal.PreviousApp.ID)) != file.Name() {
			return fmt.Errorf("rollout journal identity mismatch")
		}
		if journal.Committed {
			for _, old := range journal.Previous {
				if err := d.containerManager.StopAndRemove(ctx, old.ID); err != nil {
					return err
				}
			}
		} else {
			candidates, err := d.containerManager.dockerClient.ContainerList(ctx, container.ListOptions{All: true, Filters: filters.NewArgs(filters.Arg("label", "codedock.rollout="+journal.ID))})
			if err != nil {
				return err
			}
			for _, candidate := range candidates {
				if err := d.containerManager.StopAndRemove(ctx, candidate.ID); err != nil {
					return err
				}
			}
			for _, old := range journal.Previous {
				current, err := d.containerManager.Inspect(ctx, old.ID)
				if errdefs.IsNotFound(err) {
					continue
				}
				if err != nil {
					return err
				}
				if strings.TrimPrefix(current.Name, "/") != old.Name {
					if err := d.containerManager.dockerClient.ContainerRename(ctx, old.ID, old.Name); err != nil {
						return err
					}
				}
				if old.Running {
					if err := d.containerManager.dockerClient.ContainerStart(ctx, old.ID, container.StartOptions{}); err != nil {
						return err
					}
				}
			}
			if err := d.store.UpdateAppService(&journal.PreviousApp); err != nil {
				return err
			}
		}
		if err := d.clearRollout(journal.PreviousApp.ID); err != nil {
			return err
		}
	}
	return nil
}
