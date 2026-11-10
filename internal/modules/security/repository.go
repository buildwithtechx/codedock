package security

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type SecurityAppService struct {
	ID             string
	Name           string
	RuntimeVersion string
	BuildImage     string
	RuntimeImage   string
}

type SecurityProject struct {
	ID                     string
	ProductionDeploymentID *string
}

type SecurityDeployment struct {
	ID          string
	RawManifest string
}

type sqlProvider struct {
	db *sql.DB
}

func newSQLProvider(db *sql.DB) *sqlProvider {
	return &sqlProvider{db: db}
}

func (p *sqlProvider) GetProject(ctx context.Context, id string) (*SecurityProject, error) {
	if p.db == nil {
		return nil, nil
	}
	var proj SecurityProject
	var prodDep sql.NullString
	err := p.db.QueryRowContext(ctx, `SELECT id, production_deployment_id FROM projects WHERE id = $1`, id).Scan(&proj.ID, &prodDep)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if prodDep.Valid && prodDep.String != "" {
		proj.ProductionDeploymentID = &prodDep.String
	}
	return &proj, nil
}

func (p *sqlProvider) ListAppServices(ctx context.Context, projectID string) ([]*SecurityAppService, error) {
	if p.db == nil {
		return nil, nil
	}
	rows, err := p.db.QueryContext(ctx, `SELECT id, name, COALESCE(runtime_version, ''), COALESCE(build_image, ''), COALESCE(runtime_image, '') FROM app_services WHERE project_id = $1`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*SecurityAppService
	for rows.Next() {
		var app SecurityAppService
		if err := rows.Scan(&app.ID, &app.Name, &app.RuntimeVersion, &app.BuildImage, &app.RuntimeImage); err != nil {
			return nil, err
		}
		list = append(list, &app)
	}
	return list, rows.Err()
}

func (p *sqlProvider) UpdateAppService(ctx context.Context, svc *SecurityAppService) error {
	if p.db == nil || svc == nil {
		return nil
	}
	_, err := p.db.ExecContext(ctx, `UPDATE app_services SET runtime_version = $1, build_image = $2, runtime_image = $3 WHERE id = $4`, svc.RuntimeVersion, svc.BuildImage, svc.RuntimeImage, svc.ID)
	return err
}

func (p *sqlProvider) GetDeployment(ctx context.Context, id string) (*SecurityDeployment, error) {
	if p.db == nil {
		return nil, nil
	}
	var dep SecurityDeployment
	err := p.db.QueryRowContext(ctx, `SELECT id, COALESCE(raw_manifest, '') FROM deployments WHERE id = $1`, id).Scan(&dep.ID, &dep.RawManifest)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &dep, nil
}

func (p *sqlProvider) TriggerRemediation(ctx context.Context, serviceID string) error {
	if p.db == nil {
		return nil
	}
	_, err := p.db.ExecContext(ctx, `INSERT INTO deployments (id, service_id, status) VALUES ($1, $2, 'PENDING')`, uuid.NewString(), serviceID)
	return err
}
