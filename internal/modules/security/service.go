package security

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

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
	projectRepo   ProjectProvider
	appRepo       AppServiceProvider
	deployRepo    DeploymentProvider
	deployTrigger DeploymentTrigger
}

func NewSecurityService(
	db *sql.DB,
	feedManager *FeedManager,
	projectRepo ProjectProvider,
	appRepo AppServiceProvider,
	deployRepo DeploymentProvider,
	deployTrigger DeploymentTrigger,
) SecurityService {
	if db != nil {
		sqlProv := newSQLProvider(db)
		if projectRepo == nil {
			projectRepo = sqlProv
		}
		if appRepo == nil {
			appRepo = sqlProv
		}
		if deployRepo == nil {
			deployRepo = sqlProv
		}
		if deployTrigger == nil {
			deployTrigger = sqlProv
		}
	}
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
