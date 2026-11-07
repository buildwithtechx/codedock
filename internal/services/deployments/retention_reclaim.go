package deployments

import (
	"context"
	"log/slog"
	"strings"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"

	"codedock/internal/models"
	"codedock/internal/utils"
)

func isRestorableStatus(status models.DeploymentStatus) bool {
	s := strings.ToUpper(string(status))
	return s == "READY" || s == "ACTIVE" || s == "SUCCESS" || s == "PARTIAL_FAILURE"
}

// HasRedactedConfiguration delegates to utils.HasRedactedConfiguration.
func HasRedactedConfiguration(values ...string) bool {
	return utils.HasRedactedConfiguration(values...)
}

func isReproducible(dep *models.Deployment) bool {
	if HasRedactedConfiguration(dep.BuildLogs, dep.CommitMessage) {
		return false
	}
	if strings.TrimSpace(dep.CommitHash) != "" {
		return true
	}
	ref := strings.TrimSpace(dep.ImageRef)
	return strings.Contains(ref, "@sha256:")
}

func (s *RetentionService) ComputeRetentionInventory(ctx context.Context, projectID string) ([]*models.Deployment, []*models.Deployment, map[string]struct{}, error) {
	project, err := s.projectRepo.Get(ctx, projectID)
	if err != nil || project == nil {
		return nil, nil, nil, err
	}
	deps, err := s.deploymentRepo.ListForRetention(ctx, projectID)
	if err != nil {
		return nil, nil, nil, err
	}
	window, _ := s.ResolveRollbackWindow(project)

	retained := make([]*models.Deployment, 0)
	overflow := make([]*models.Deployment, 0)
	retainedImages := make(map[string]struct{})

	pastReleases := 0
	for i, dep := range deps {
		isFirstOrActive := i == 0 || strings.ToUpper(string(dep.Status)) == "ACTIVE"
		if isFirstOrActive || dep.Pinned {
			retained = append(retained, dep)
		} else if isRestorableStatus(dep.Status) && pastReleases < window {
			retained = append(retained, dep)
			pastReleases++
		} else if dep.ArtifactRetainedAt != nil {
			overflow = append(overflow, dep)
		}
	}

	for _, dep := range retained {
		ref := strings.TrimSpace(dep.ImageRef)
		if ref != "" {
			retainedImages[ref] = struct{}{}
		}
	}

	return retained, overflow, retainedImages, nil
}

func (s *RetentionService) PruneProject(ctx context.Context, projectID string) (*models.RetentionReclaimSummary, error) {
	summary := &models.RetentionReclaimSummary{}
	_, overflow, retainedImages, err := s.ComputeRetentionInventory(ctx, projectID)
	if err != nil {
		return summary, err
	}

	var dockerCli *client.Client
	if s.dockerProvider != nil {
		if cli, rel, err := s.dockerProvider(ctx, projectID); err == nil {
			dockerCli = cli
			defer rel()
		}
	}

	for _, candidate := range overflow {
		ref := strings.TrimSpace(candidate.ImageRef)
		if ref != "" && dockerCli != nil {
			if _, inUse := retainedImages[ref]; !inUse {
				if _, err := dockerCli.ImageRemove(ctx, ref, image.RemoveOptions{Force: false, PruneChildren: true}); err != nil {
					if !client.IsErrNotFound(err) {
						slog.Warn("retention prune could not remove image", "image", ref, "error", err)
						summary.Failed++
						continue
					}
				}
			}
		}
		if err := s.deploymentRepo.SetArtifactRetainedAt(ctx, candidate.ID, nil); err != nil {
			summary.Failed++
		} else {
			summary.Purged++
		}
	}

	pressureSummary, err := s.ReclaimStoragePressure(ctx, projectID)
	if err == nil && pressureSummary != nil {
		summary.Purged += pressureSummary.Purged
		summary.Failed += pressureSummary.Failed
	}
	return summary, nil
}

func (s *RetentionService) ReclaimStoragePressure(ctx context.Context, projectID string) (*models.RetentionReclaimSummary, error) {
	summary := &models.RetentionReclaimSummary{}
	project, err := s.projectRepo.Get(ctx, projectID)
	if err != nil || project == nil {
		return summary, err
	}
	if project.RollbackWindow != nil {
		return summary, nil
	}

	total, free, err := s.diskReader(ctx, project)
	if err != nil || !rollbackStoragePressure(total, free) {
		return summary, nil
	}

	retained, _, _, err := s.ComputeRetentionInventory(ctx, projectID)
	if err != nil {
		return summary, err
	}

	var dockerCli *client.Client
	if s.dockerProvider != nil {
		if cli, rel, err := s.dockerProvider(ctx, projectID); err == nil {
			dockerCli = cli
			defer rel()
		}
	}

	for i := len(retained) - 1; i >= 0; i-- {
		candidate := retained[i]
		if candidate.ArtifactRetainedAt == nil || candidate.Pinned || i == 0 {
			continue
		}
		if !isReproducible(candidate) {
			continue
		}

		ref := strings.TrimSpace(candidate.ImageRef)
		if ref != "" && dockerCli != nil {
			otherRetainedUses := false
			for j, other := range retained {
				if j != i && strings.TrimSpace(other.ImageRef) == ref {
					otherRetainedUses = true
					break
				}
			}
			if !otherRetainedUses {
				if _, err := dockerCli.ImageRemove(ctx, ref, image.RemoveOptions{Force: false, PruneChildren: true}); err != nil {
					if !client.IsErrNotFound(err) {
						slog.Warn("pressure reclaim could not remove image", "image", ref, "error", err)
						summary.Failed++
						continue
					}
				}
			}
		}

		if err := s.deploymentRepo.SetArtifactRetainedAt(ctx, candidate.ID, nil); err != nil {
			summary.Failed++
			continue
		}
		candidate.ArtifactRetainedAt = nil
		summary.Purged++

		total, free, err = s.diskReader(ctx, project)
		if err != nil || !rollbackStoragePressure(total, free) {
			break
		}
	}

	return summary, nil
}
