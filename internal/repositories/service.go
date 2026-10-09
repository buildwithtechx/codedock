package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"codedock/internal/models"
	"codedock/internal/utils"
)

type AppServiceRepository interface {
	Create(ctx context.Context, svc *models.AppService) error
	GetByID(ctx context.Context, id string) (*models.AppService, error)
	ListByEnvironment(ctx context.Context, environmentID string) ([]*models.AppService, error)
	ListByProject(ctx context.Context, projectID string) ([]*models.AppService, error)
	ListByOrganization(ctx context.Context, organizationID string) ([]*models.AppService, error)
	ListAll(ctx context.Context) ([]*models.AppService, error)
	Update(ctx context.Context, svc *models.AppService) error
	Delete(ctx context.Context, id string) error
	CreateWebhook(ctx context.Context, w *models.Webhook) error
	ListWebhooksByService(ctx context.Context, serviceID string) ([]*models.Webhook, error)
	DeleteWebhook(ctx context.Context, id, serviceID string) error
	CreateLogDrain(ctx context.Context, d *models.LogDrain) error
	ListLogDrainsByService(ctx context.Context, serviceID string) ([]*models.LogDrain, error)
	DeleteLogDrain(ctx context.Context, id, serviceID string) error
}

type AppServiceRepo struct {
	mu sync.RWMutex
	db *sqlx.DB
}

func NewAppServiceRepo(db *sql.DB) *AppServiceRepo {
	return &AppServiceRepo{db: sqlx.NewDb(db, "pgx")}
}

func (r *AppServiceRepo) Create(_ context.Context, svc *models.AppService) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if svc.ID == "" {
		svc.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	svc.CreatedAt = now
	svc.UpdatedAt = now
	if svc.Status == "" {
		svc.Status = "building"
	}
	if svc.InternalPort == 0 {
		svc.InternalPort = 3000
	}
	_, err := r.db.Exec(
		`INSERT INTO app_services (id, project_id, environment_id, name, repository_url, git_user_id, image_ref, branch, root_directory, icon, runtime_mode, install_command, build_command, start_command, dockerfile_path, build_engine, internal_port, domain, static_output, health_check_path, container_id, status, replicas, cpu_limit, memory_limit, deploy_token, enable_pr_previews, maintenance_mode, registry_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31)`,
		svc.ID, svc.ProjectID, svc.EnvironmentID, svc.Name, svc.RepositoryURL, svc.GitUserID, svc.ImageRef, svc.Branch,
		svc.RootDirectory, svc.Icon, svc.RuntimeMode, svc.InstallCommand, svc.BuildCommand, svc.StartCommand, svc.DockerfilePath, svc.BuildEngine,
		svc.InternalPort, svc.Domain, svc.StaticOutput, svc.HealthCheckPath, svc.ContainerID, svc.Status, svc.Replicas, svc.CPULimit, svc.MemoryLimit, svc.DeployToken, svc.EnablePRPreviews, svc.MaintenanceMode, svc.RegistryID, svc.CreatedAt, svc.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create app service: %w", err)
	}
	return nil
}

func (r *AppServiceRepo) GetByID(ctx context.Context, id string) (*models.AppService, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var svc models.AppService
	err := r.db.GetContext(ctx, &svc,
		`SELECT COALESCE((SELECT app_id FROM projects WHERE projects.id = app_services.project_id), '') AS app_id, id, project_id, environment_id, name, repository_url, git_user_id, COALESCE(image_ref,'') AS image_ref, branch, root_directory, COALESCE(icon,'git') AS icon, runtime_mode, COALESCE(install_command,'') AS install_command, build_command, start_command, dockerfile_path, build_engine, internal_port, domain, COALESCE(static_output,'') AS static_output, health_check_path, container_id, status, replicas, COALESCE(cpu_limit, 0) AS cpu_limit, COALESCE(memory_limit, 0) AS memory_limit, COALESCE(deploy_token,'') AS deploy_token, COALESCE(enable_pr_previews, FALSE) AS enable_pr_previews, COALESCE(maintenance_mode, FALSE) AS maintenance_mode, registry_id, created_at, updated_at
		FROM app_services WHERE id = $1`, id,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("AppService", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get app service: %w", err)
	}
	return &svc, nil
}

func (r *AppServiceRepo) ListByEnvironment(ctx context.Context, environmentID string) ([]*models.AppService, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*models.AppService
	err := r.db.SelectContext(ctx, &list,
		`SELECT COALESCE((SELECT app_id FROM projects WHERE projects.id = app_services.project_id), '') AS app_id, id, project_id, environment_id, name, repository_url, git_user_id, COALESCE(image_ref,'') AS image_ref, branch, root_directory, COALESCE(icon,'git') AS icon, runtime_mode, COALESCE(install_command,'') AS install_command, build_command, start_command, dockerfile_path, build_engine, internal_port, domain, COALESCE(static_output,'') AS static_output, health_check_path, container_id, status, replicas, COALESCE(cpu_limit, 0) AS cpu_limit, COALESCE(memory_limit, 0) AS memory_limit, COALESCE(deploy_token,'') AS deploy_token, COALESCE(enable_pr_previews, FALSE) AS enable_pr_previews, COALESCE(maintenance_mode, FALSE) AS maintenance_mode, registry_id, created_at, updated_at
		FROM app_services WHERE environment_id = $1 ORDER BY created_at ASC`, environmentID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list app services by environment: %w", err)
	}
	if list == nil {
		list = make([]*models.AppService, 0)
	}
	return list, nil
}

func (r *AppServiceRepo) ListByProject(ctx context.Context, projectID string) ([]*models.AppService, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*models.AppService
	err := r.db.SelectContext(ctx, &list,
		`SELECT COALESCE((SELECT app_id FROM projects WHERE projects.id = app_services.project_id), '') AS app_id, id, project_id, environment_id, name, repository_url, git_user_id, COALESCE(image_ref,'') AS image_ref, branch, root_directory, COALESCE(icon,'git') AS icon, runtime_mode, COALESCE(install_command,'') AS install_command, build_command, start_command, dockerfile_path, build_engine, internal_port, domain, COALESCE(static_output,'') AS static_output, health_check_path, container_id, status, replicas, COALESCE(cpu_limit, 0) AS cpu_limit, COALESCE(memory_limit, 0) AS memory_limit, COALESCE(deploy_token,'') AS deploy_token, COALESCE(enable_pr_previews, FALSE) AS enable_pr_previews, COALESCE(maintenance_mode, FALSE) AS maintenance_mode, registry_id, created_at, updated_at
		FROM app_services WHERE project_id = $1 ORDER BY created_at ASC`, projectID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list app services by project: %w", err)
	}
	if list == nil {
		list = make([]*models.AppService, 0)
	}
	return list, nil
}

func (r *AppServiceRepo) ListByOrganization(ctx context.Context, organizationID string) ([]*models.AppService, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*models.AppService
	err := r.db.SelectContext(ctx, &list,
		`SELECT COALESCE(projects.app_id, '') AS app_id, app_services.id, app_services.project_id, app_services.environment_id, app_services.name, app_services.repository_url, app_services.git_user_id, COALESCE(app_services.image_ref,'') AS image_ref, app_services.branch, app_services.root_directory, COALESCE(app_services.icon,'git') AS icon, app_services.runtime_mode, COALESCE(app_services.install_command,'') AS install_command, app_services.build_command, app_services.start_command, app_services.dockerfile_path, app_services.build_engine, app_services.internal_port, app_services.domain, COALESCE(app_services.static_output,'') AS static_output, app_services.health_check_path, app_services.container_id, app_services.status, app_services.replicas, COALESCE(app_services.cpu_limit, 0) AS cpu_limit, COALESCE(app_services.memory_limit, 0) AS memory_limit, COALESCE(app_services.deploy_token,'') AS deploy_token, COALESCE(app_services.enable_pr_previews, FALSE) AS enable_pr_previews, COALESCE(app_services.maintenance_mode, FALSE) AS maintenance_mode, app_services.registry_id, app_services.created_at, app_services.updated_at
		FROM app_services JOIN projects ON projects.id = app_services.project_id
		WHERE projects.organization_id = $1 ORDER BY app_services.created_at DESC`, organizationID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list app services by organization: %w", err)
	}
	if list == nil {
		list = make([]*models.AppService, 0)
	}
	return list, nil
}

func (r *AppServiceRepo) ListAll(ctx context.Context) ([]*models.AppService, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*models.AppService
	err := r.db.SelectContext(ctx, &list,
		`SELECT COALESCE((SELECT app_id FROM projects WHERE projects.id = app_services.project_id), '') AS app_id, id, project_id, environment_id, name, repository_url, git_user_id, COALESCE(image_ref,'') AS image_ref, branch, root_directory, COALESCE(icon,'git') AS icon, runtime_mode, COALESCE(install_command,'') AS install_command, build_command, start_command, dockerfile_path, build_engine, internal_port, domain, COALESCE(static_output,'') AS static_output, health_check_path, container_id, status, replicas, COALESCE(cpu_limit, 0) AS cpu_limit, COALESCE(memory_limit, 0) AS memory_limit, COALESCE(deploy_token,'') AS deploy_token, COALESCE(enable_pr_previews, FALSE) AS enable_pr_previews, COALESCE(maintenance_mode, FALSE) AS maintenance_mode, registry_id, created_at, updated_at
		FROM app_services ORDER BY created_at ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list all app services: %w", err)
	}
	if list == nil {
		list = make([]*models.AppService, 0)
	}
	return list, nil
}

func (r *AppServiceRepo) Update(_ context.Context, svc *models.AppService) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	svc.UpdatedAt = time.Now().UTC()
	_, err := r.db.Exec(
		`UPDATE app_services SET
		name = $1, repository_url = $2, git_user_id = $3, image_ref = $4, branch = $5, root_directory = $6, icon = $7, runtime_mode = $8,
		install_command = $9, build_command = $10, start_command = $11, dockerfile_path = $12, build_engine = $13,
		internal_port = $14, domain = $15, static_output = $16, health_check_path = $17, container_id = $18, status = $19, replicas = $20, cpu_limit = $21, memory_limit = $22, deploy_token = $23, enable_pr_previews = $24, maintenance_mode = $25, registry_id = $26, updated_at = $27
		WHERE id = $28`,
		svc.Name, svc.RepositoryURL, svc.GitUserID, svc.ImageRef, svc.Branch, svc.RootDirectory, svc.Icon, svc.RuntimeMode,
		svc.InstallCommand, svc.BuildCommand, svc.StartCommand, svc.DockerfilePath, svc.BuildEngine,
		svc.InternalPort, svc.Domain, svc.StaticOutput, svc.HealthCheckPath, svc.ContainerID, svc.Status, svc.Replicas, svc.CPULimit, svc.MemoryLimit, svc.DeployToken, svc.EnablePRPreviews, svc.MaintenanceMode, svc.RegistryID, svc.UpdatedAt, svc.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update app service: %w", err)
	}
	return nil
}

func (r *AppServiceRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.Exec(`DELETE FROM app_services WHERE id = $1`, id)
	return err
}

func (r *AppServiceRepo) CreateWebhook(ctx context.Context, w *models.Webhook) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if w.ID == "" {
		w.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	w.CreatedAt = now
	w.UpdatedAt = now
	eventTypesStr := strings.Join(w.EventTypes, ",")
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO service_webhooks (id, service_id, url, event_types, include_pr_environments, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		w.ID, w.ServiceID, w.URL, eventTypesStr, w.IncludePREnvironments, w.CreatedAt, w.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create webhook: %w", err)
	}
	return nil
}

func (r *AppServiceRepo) ListWebhooksByService(ctx context.Context, serviceID string) ([]*models.Webhook, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, service_id, url, event_types, include_pr_environments, created_at, updated_at
		 FROM service_webhooks WHERE service_id = $1 ORDER BY created_at DESC`, serviceID)
	if err != nil {
		return nil, fmt.Errorf("list webhooks: %w", err)
	}
	defer rows.Close()
	var out []*models.Webhook
	for rows.Next() {
		var w models.Webhook
		var eventsStr string
		if err := rows.Scan(&w.ID, &w.ServiceID, &w.URL, &eventsStr, &w.IncludePREnvironments, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan webhook: %w", err)
		}
		if eventsStr != "" {
			w.EventTypes = strings.Split(eventsStr, ",")
		} else {
			w.EventTypes = []string{}
		}
		out = append(out, &w)
	}
	return out, rows.Err()
}

func (r *AppServiceRepo) DeleteWebhook(ctx context.Context, id, serviceID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	query := "DELETE FROM service_webhooks WHERE id = $1 AND service_id = $2"
	res, err := r.db.ExecContext(ctx, query, id, serviceID)
	if err != nil {
		return fmt.Errorf("delete webhook: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return utils.NewNotFoundError("Webhook", id)
	}
	return nil
}

func (r *AppServiceRepo) CreateLogDrain(ctx context.Context, d *models.LogDrain) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	d.CreatedAt = now
	d.UpdatedAt = now
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO log_drains (id, service_id, project_id, drain_type, endpoint_url, auth_token, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		d.ID, d.ServiceID, d.ProjectID, d.DrainType, d.EndpointURL, d.AuthToken, d.CreatedAt, d.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create log drain: %w", err)
	}
	return nil
}

func (r *AppServiceRepo) ListLogDrainsByService(ctx context.Context, serviceID string) ([]*models.LogDrain, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, service_id, project_id, drain_type, endpoint_url, COALESCE(auth_token, '') AS auth_token, created_at, updated_at
		 FROM log_drains WHERE service_id = $1 ORDER BY created_at DESC`, serviceID)
	if err != nil {
		return nil, fmt.Errorf("list log drains: %w", err)
	}
	defer rows.Close()
	var out []*models.LogDrain
	for rows.Next() {
		var d models.LogDrain
		if err := rows.Scan(&d.ID, &d.ServiceID, &d.ProjectID, &d.DrainType, &d.EndpointURL, &d.AuthToken, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan log drain: %w", err)
		}
		out = append(out, &d)
	}
	return out, rows.Err()
}

func (r *AppServiceRepo) DeleteLogDrain(ctx context.Context, id, serviceID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	query := "DELETE FROM log_drains WHERE id = $1 AND service_id = $2"
	res, err := r.db.ExecContext(ctx, query, id, serviceID)
	if err != nil {
		return fmt.Errorf("delete log drain: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return utils.NewNotFoundError("LogDrain", id)
	}
	return nil
}
