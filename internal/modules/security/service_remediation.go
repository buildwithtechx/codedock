package security

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (s *securityService) CheckProject(ctx context.Context, projectID string, force bool) (*ProjectSecurityResult, error) {
	manifest, err := s.feedManager.GetManifest(ctx)
	if err != nil {
		return nil, fmt.Errorf("load security manifest: %w", err)
	}

	var serviceEvidences []ServiceEvidence
	var activeDeploymentID string

	var services []*SecurityAppService
	if s.appRepo != nil {
		services, err = s.appRepo.ListAppServices(ctx, projectID)
		if err != nil {
			return nil, fmt.Errorf("list project services: %w", err)
		}
	}

	for _, svc := range services {
		image := ""
		if svc.BuildImage != "" {
			image = svc.BuildImage
		} else if svc.RuntimeImage != "" {
			image = svc.RuntimeImage
		}
		isPinned := strings.Contains(image, "@sha256:")
		digest := ""
		if idx := strings.Index(image, "@sha256:"); idx != -1 {
			digest = "sha256:" + image[idx+8:]
		}

		runtimeVersions := make(map[string]string)
		if svc.RuntimeVersion != "" {
			runtimeVersions["node"] = svc.RuntimeVersion
		}

		serviceEvidences = append(serviceEvidences, ServiceEvidence{
			ServiceName:     svc.Name,
			ServiceID:       svc.ID,
			Image:           image,
			ImageDigest:     digest,
			BuildImage:      svc.BuildImage,
			RuntimeImage:    svc.RuntimeImage,
			IsPinnedDigest:  isPinned,
			RuntimeVersions: runtimeVersions,
		})
	}

	if s.projectRepo != nil {
		if proj, err := s.projectRepo.GetProject(ctx, projectID); err == nil && proj != nil && proj.ProductionDeploymentID != nil {
			activeDeploymentID = *proj.ProductionDeploymentID
		}
	}

	findings := MatchAdvisories(manifest.Advisories, serviceEvidences)

	// Enrich runtime findings with precision semver runtime patches
	for i := range findings {
		f := &findings[i]
		if f.TargetKind == TargetKindRuntime {
			for _, adv := range manifest.Advisories {
				if adv.ID != f.AdvisoryID {
					continue
				}
				for _, target := range adv.Targets {
					if target.Kind != TargetKindRuntime || target.Fixed == nil || *target.Fixed == "" {
						continue
					}
					for _, svc := range services {
						if (f.ServiceID != "" && svc.ID == f.ServiceID) || (f.ServiceID == "" && svc.Name == f.ServiceName) {
							ver := svc.RuntimeVersion
							if ver == "" {
								continue
							}
							buildImg := svc.BuildImage
							runtimeImg := svc.RuntimeImage
							if buildImg == "" {
								if target.Runtime == "bun" {
									buildImg = fmt.Sprintf("oven/bun:%s-alpine", ver)
								} else {
									buildImg = fmt.Sprintf("node:%s-bookworm-slim", ver)
								}
							}
							if runtimeImg == "" {
								if target.Runtime == "bun" {
									runtimeImg = fmt.Sprintf("oven/bun:%s-alpine", ver)
								} else {
									runtimeImg = fmt.Sprintf("node:%s-bookworm-slim", ver)
								}
							}

							patch, err := PlanSecurityRuntimePatch(
								projectID,
								activeDeploymentID,
								adv.ID,
								adv.Revision,
								target.Runtime,
								ver,
								*target.Fixed,
								buildImg,
								runtimeImg,
							)
							if err == nil && patch != nil {
								f.RuntimePatch = patch
								f.CanRefresh = true
								f.RuntimeUpgradeRequired = false
							}
						}
					}
				}
			}
		}
	}
	now := time.Now().UTC()

	result := &ProjectSecurityResult{
		ProjectID:    projectID,
		DeploymentID: activeDeploymentID,
		CheckedAt:    now,
		Findings:     findings,
	}

	if s.db != nil {
		invBytes, _ := json.Marshal(findings)
		query := `
			INSERT INTO project_security (project_id, deployment_id, inventory, checked_at, updated_at)
			VALUES ($1, $2, $3, $4, $4)
			ON CONFLICT (project_id) DO UPDATE SET
				deployment_id = EXCLUDED.deployment_id,
				inventory = EXCLUDED.inventory,
				checked_at = EXCLUDED.checked_at,
				updated_at = EXCLUDED.updated_at
		`
		_, _ = s.db.ExecContext(ctx, query, projectID, activeDeploymentID, invBytes, now)
	}

	return result, nil
}

func (s *securityService) ApplyRemediation(ctx context.Context, projectID, advisoryID, deploymentID string, serviceID, runtimePatchID *string) (*RefreshAttempt, error) {
	sec, err := s.GetProjectSecurity(ctx, projectID)
	if err != nil {
		return nil, err
	}

	var targetFinding *AdvisoryFinding
	for _, f := range sec.Findings {
		if f.AdvisoryID == advisoryID {
			targetFinding = &f
			break
		}
	}

	if targetFinding == nil {
		return nil, fmt.Errorf("finding %q not found on project %q", advisoryID, projectID)
	}
	if !targetFinding.CanRefresh {
		return nil, errors.New("this advisory requires manual configuration review and cannot be refreshed automatically")
	}

	// Validate runtime patch review ID if provided
	if runtimePatchID != nil && *runtimePatchID != "" {
		if targetFinding.RuntimePatch == nil {
			return nil, errors.New("no planned runtime patch available for this advisory")
		}
		if targetFinding.RuntimePatch.ReviewID != *runtimePatchID {
			return nil, errors.New("the runtime patch has changed or expired; please review the updated patch")
		}
	}

	// Persist patched images and version to project service if runtime patch is present
	if targetFinding.RuntimePatch != nil && s.appRepo != nil {
		services, err := s.appRepo.ListAppServices(ctx, projectID)
		if err == nil {
			for _, svc := range services {
				matchesService := false
				if serviceID != nil && *serviceID != "" && svc.ID == *serviceID {
					matchesService = true
				} else if targetFinding.ServiceID != "" && svc.ID == targetFinding.ServiceID {
					matchesService = true
				} else if svc.Name == targetFinding.ServiceName || targetFinding.ServiceName == "" {
					matchesService = true
				}

				if matchesService {
					svc.BuildImage = targetFinding.RuntimePatch.BuildImage
					svc.RuntimeImage = targetFinding.RuntimePatch.RuntimeImage
					svc.RuntimeVersion = targetFinding.RuntimePatch.TargetVersion
					_ = s.appRepo.UpdateAppService(ctx, svc)
				}
			}
		}
	}

	attempt := &RefreshAttempt{
		ID:           uuid.New().String(),
		ProjectID:    projectID,
		DeploymentID: deploymentID,
		AdvisoryID:   advisoryID,
		Status:       "queued",
		CreatedAt:    time.Now().UTC(),
	}

	if s.db != nil {
		query := `
			INSERT INTO security_refresh_attempts (id, project_id, deployment_id, advisory_id, status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $6)
		`
		if _, err := s.db.ExecContext(ctx, query, attempt.ID, projectID, deploymentID, advisoryID, attempt.Status, attempt.CreatedAt); err != nil {
			return nil, fmt.Errorf("create refresh attempt: %w", err)
		}
	}

	// Trigger deployment if trigger is provided
	if s.deployTrigger != nil && s.appRepo != nil {
		services, _ := s.appRepo.ListAppServices(ctx, projectID)
		for _, svc := range services {
			matchesService := false
			if serviceID != nil && *serviceID != "" && svc.ID == *serviceID {
				matchesService = true
			} else if targetFinding.ServiceID != "" && svc.ID == targetFinding.ServiceID {
				matchesService = true
			} else if svc.Name == targetFinding.ServiceName || targetFinding.ServiceName == "" {
				matchesService = true
			}

			if matchesService {
				_ = s.deployTrigger.TriggerRemediation(ctx, svc.ID)
			}
		}
	}

	return attempt, nil
}
