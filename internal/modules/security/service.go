package security

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"codedock/internal/models"
	"codedock/internal/repositories"
)

type DeploymentTrigger interface {
	CreateDeployment(ctx context.Context, d *models.Deployment) (*models.Deployment, error)
}

type SecurityService interface {
	GetAdvisories(ctx context.Context) (*AdvisoriesResponse, error)
	CheckProject(ctx context.Context, projectID string, force bool) (*ProjectSecurityResult, error)
	GetProjectSecurity(ctx context.Context, projectID string) (*ProjectSecurityResult, error)
	SetProjectPolicy(ctx context.Context, projectID string, autoRefresh bool) error
	ApplyRemediation(ctx context.Context, projectID, advisoryID, deploymentID string, serviceID, runtimePatchID *string) (*RefreshAttempt, error)
}

type securityService struct {
	db            *sql.DB
	feedManager   *FeedManager
	projectRepo   repositories.ProjectRepository
	appRepo       repositories.AppServiceRepository
	deployRepo    repositories.DeploymentRepository
	deployTrigger DeploymentTrigger
}

func NewSecurityService(
	db *sql.DB,
	feedManager *FeedManager,
	projectRepo repositories.ProjectRepository,
	appRepo repositories.AppServiceRepository,
	deployRepo repositories.DeploymentRepository,
	deployTrigger DeploymentTrigger,
) SecurityService {
	return &securityService{
		db:            db,
		feedManager:   feedManager,
		projectRepo:   projectRepo,
		appRepo:       appRepo,
		deployRepo:    deployRepo,
		deployTrigger: deployTrigger,
	}
}

func (s *securityService) GetAdvisories(ctx context.Context) (*AdvisoriesResponse, error) {
	manifest, err := s.feedManager.GetManifest(ctx)
	if err != nil {
		return nil, fmt.Errorf("load security manifest: %w", err)
	}

	var affectedProjects []ProjectSecurityResult
	if s.db != nil {
		rows, err := s.db.QueryContext(ctx, `
			SELECT project_id, deployment_id, inventory, checked_at, auto_refresh
			FROM project_security
			ORDER BY checked_at DESC
		`)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var p ProjectSecurityResult
				var invRaw []byte
				if err := rows.Scan(&p.ProjectID, &p.DeploymentID, &invRaw, &p.CheckedAt, &p.AutoRefresh); err == nil {
					if len(invRaw) > 0 {
						_ = json.Unmarshal(invRaw, &p.Findings)
					}
					if len(p.Findings) > 0 {
						affectedProjects = append(affectedProjects, p)
					}
				}
			}
		}
	}

	now := time.Now().UTC()
	stale := now.After(manifest.ExpiresAt)

	return &AdvisoriesResponse{
		Revision:         manifest.Revision,
		CheckedAt:        manifest.IssuedAt,
		ExpiresAt:        manifest.ExpiresAt,
		Stale:            stale,
		Advisories:       manifest.Advisories,
		AffectedProjects: affectedProjects,
	}, nil
}

func (s *securityService) GetProjectSecurity(ctx context.Context, projectID string) (*ProjectSecurityResult, error) {
	if s.db == nil {
		return s.CheckProject(ctx, projectID, false)
	}

	var p ProjectSecurityResult
	var invRaw []byte
	err := s.db.QueryRowContext(ctx, `
		SELECT project_id, deployment_id, inventory, checked_at, auto_refresh
		FROM project_security
		WHERE project_id = $1
	`, projectID).Scan(&p.ProjectID, &p.DeploymentID, &invRaw, &p.CheckedAt, &p.AutoRefresh)

	if err == sql.ErrNoRows || time.Since(p.CheckedAt) > 6*time.Hour {
		return s.CheckProject(ctx, projectID, false)
	}
	if err != nil {
		return nil, fmt.Errorf("load project security: %w", err)
	}

	if len(invRaw) > 0 {
		_ = json.Unmarshal(invRaw, &p.Findings)
	}

	// Check for active refresh attempt
	var att RefreshAttempt
	err = s.db.QueryRowContext(ctx, `
		SELECT id, project_id, deployment_id, advisory_id, status, error, created_at
		FROM security_refresh_attempts
		WHERE project_id = $1 AND status IN ('pending', 'queued')
		ORDER BY created_at DESC LIMIT 1
	`, projectID).Scan(&att.ID, &att.ProjectID, &att.DeploymentID, &att.AdvisoryID, &att.Status, &att.Error, &att.CreatedAt)
	if err == nil {
		p.ActiveRefreshAttempt = &att
	}

	return &p, nil
}

func (s *securityService) CheckProject(ctx context.Context, projectID string, force bool) (*ProjectSecurityResult, error) {
	manifest, err := s.feedManager.GetManifest(ctx)
	if err != nil {
		return nil, fmt.Errorf("load security manifest: %w", err)
	}

	services, err := s.appRepo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("list project services: %w", err)
	}

	var serviceEvidences []ServiceEvidence
	var activeDeploymentID string

	for _, svc := range services {
		image := svc.ImageRef
		isPinned := strings.Contains(image, "@sha256:")
		digest := ""
		if idx := strings.Index(image, "@sha256:"); idx != -1 {
			digest = "sha256:" + image[idx+8:]
		}

		// Check if deployment exists for service to get active deployment ID
		if activeDeploymentID == "" && s.deployRepo != nil {
			if deps, _, err := s.deployRepo.ListByService(ctx, svc.ID, 1, 0); err == nil && len(deps) > 0 {
				activeDeploymentID = deps[0].ID
			}
		}

		runtimeVersions := make(map[string]string)
		if svc.Stack != "" && svc.RuntimeVersion != "" {
			runtimeVersions[strings.ToLower(svc.Stack)] = svc.RuntimeVersion
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
							if err != nil {
								fmt.Printf("DEBUG patch err: %v\n", err)
							}
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

func (s *securityService) SetProjectPolicy(ctx context.Context, projectID string, autoRefresh bool) error {
	if s.db == nil {
		return nil
	}
	query := `
		INSERT INTO project_security (project_id, auto_refresh, checked_at, updated_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (project_id) DO UPDATE SET
			auto_refresh = EXCLUDED.auto_refresh,
			updated_at = CURRENT_TIMESTAMP
	`
	_, err := s.db.ExecContext(ctx, query, projectID, autoRefresh)
	if err != nil {
		return fmt.Errorf("set project security policy: %w", err)
	}
	return nil
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
	if targetFinding.RuntimePatch != nil {
		services, err := s.appRepo.ListByProject(ctx, projectID)
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
					_ = s.appRepo.Update(ctx, svc)
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
	if s.deployTrigger != nil {
		services, _ := s.appRepo.ListByProject(ctx, projectID)
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
				_, _ = s.deployTrigger.CreateDeployment(ctx, &models.Deployment{
					ID:        uuid.New().String(),
					ServiceID: svc.ID,
					Status:    models.DeploymentStatusPending,
				})
			}
		}
	}

	return attempt, nil
}
