package apps

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

	"codedock/internal/utils"
)

type AppServiceRepo struct {
	mu sync.RWMutex
	db *sqlx.DB
}

func NewAppServiceRepo(db *sql.DB) *AppServiceRepo {
	return &AppServiceRepo{db: sqlx.NewDb(db, "pgx")}
}

const appServiceSelectColumns = `
	COALESCE((SELECT app_id FROM projects WHERE projects.id = app_services.project_id), '') AS app_id,
	id, project_id, environment_id, name, repository_url, git_user_id, COALESCE(image_ref,'') AS image_ref,
	branch, root_directory, COALESCE(icon,'git') AS icon, runtime_mode,
	COALESCE(install_command,'') AS install_command, build_command, start_command, dockerfile_path,
	build_engine, internal_port, domain, COALESCE(static_output,'') AS static_output,
	health_check_path, container_id, status, replicas, COALESCE(cpu_limit, 0) AS cpu_limit,
	COALESCE(memory_limit, 0) AS memory_limit, COALESCE(swap_disabled, FALSE) AS swap_disabled,
	COALESCE(deploy_token,'') AS deploy_token, COALESCE(enable_pr_previews, FALSE) AS enable_pr_previews,
	COALESCE(maintenance_mode, FALSE) AS maintenance_mode, registry_id, created_at, updated_at,
	COALESCE(optimize, FALSE) AS optimize, COALESCE(stack, '') AS stack,
	COALESCE(stop_grace_period_seconds, 0) AS stop_grace_period_seconds,
	COALESCE(runtime_version, '') AS runtime_version, COALESCE(spa_fallback, TRUE) AS spa_fallback,
	COALESCE(build_image, '') AS build_image, COALESCE(runtime_image, '') AS runtime_image`

func (r *AppServiceRepo) Create(_ context.Context, svc *AppService) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if svc.ID == "" {
		svc.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	svc.CreatedAt = now
	svc.UpdatedAt = now
	if svc.Status == "" {
		svc.Status = AppServiceStatusBuilding
	}
	if svc.InternalPort == 0 {
		svc.InternalPort = 3000
	}
	spaFallback := true
	if svc.SPAFallback != nil {
		spaFallback = *svc.SPAFallback
	}
	_, err := r.db.Exec(
		`INSERT INTO app_services (id, project_id, environment_id, name, repository_url, git_user_id, image_ref, branch, root_directory, icon, runtime_mode, install_command, build_command, start_command, dockerfile_path, build_engine, internal_port, domain, static_output, health_check_path, container_id, status, replicas, cpu_limit, memory_limit, swap_disabled, deploy_token, enable_pr_previews, maintenance_mode, registry_id, created_at, updated_at, optimize, stack, stop_grace_period_seconds, runtime_version, spa_fallback, build_image, runtime_image)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39)`,
		svc.ID, svc.ProjectID, svc.EnvironmentID, svc.Name, svc.RepositoryURL, svc.GitUserID, svc.ImageRef, svc.Branch,
		svc.RootDirectory, svc.Icon, svc.RuntimeMode, svc.InstallCommand, svc.BuildCommand, svc.StartCommand, svc.DockerfilePath, svc.BuildEngine,
		svc.InternalPort, svc.Domain, svc.StaticOutput, svc.HealthCheckPath, svc.ContainerID, svc.Status, svc.Replicas, svc.CPULimit, svc.MemoryLimit, svc.SwapDisabled, svc.DeployToken, svc.EnablePRPreviews, svc.MaintenanceMode, svc.RegistryID, svc.CreatedAt, svc.UpdatedAt, svc.Optimize, svc.Stack, svc.StopGracePeriodSeconds, svc.RuntimeVersion, spaFallback, svc.BuildImage, svc.RuntimeImage,
	)
	if err != nil {
		return fmt.Errorf("failed to create app service: %w", err)
	}
	return nil
}

func (r *AppServiceRepo) GetByID(ctx context.Context, id string) (*AppService, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var svc AppService
	err := r.db.GetContext(ctx, &svc,
		`SELECT `+appServiceSelectColumns+` FROM app_services WHERE id = $1`, id,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("AppService", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get app service: %w", err)
	}
	return &svc, nil
}

func (r *AppServiceRepo) ListByEnvironment(ctx context.Context, environmentID string) ([]*AppService, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*AppService
	err := r.db.SelectContext(ctx, &list,
		`SELECT `+appServiceSelectColumns+` FROM app_services WHERE environment_id = $1 ORDER BY created_at ASC`, environmentID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list app services by environment: %w", err)
	}
	if list == nil {
		list = make([]*AppService, 0)
	}
	return list, nil
}

func (r *AppServiceRepo) ListByProject(ctx context.Context, projectID string) ([]*AppService, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*AppService
	err := r.db.SelectContext(ctx, &list,
		`SELECT `+appServiceSelectColumns+` FROM app_services WHERE project_id = $1 ORDER BY created_at ASC`, projectID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list app services by project: %w", err)
	}
	if list == nil {
		list = make([]*AppService, 0)
	}
	return list, nil
}

func (r *AppServiceRepo) ListByOrganization(ctx context.Context, organizationID string, search string) ([]*AppService, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*AppService
	search = strings.TrimSpace(search)
	pattern := "%" + search + "%"
	err := r.db.SelectContext(ctx, &list,
		`SELECT COALESCE(projects.app_id, '') AS app_id, app_services.id, app_services.project_id, app_services.environment_id, app_services.name, app_services.repository_url, app_services.git_user_id, COALESCE(app_services.image_ref,'') AS image_ref, app_services.branch, app_services.root_directory, COALESCE(app_services.icon,'git') AS icon, app_services.runtime_mode, COALESCE(app_services.install_command,'') AS install_command, app_services.build_command, app_services.start_command, app_services.dockerfile_path, app_services.build_engine, app_services.internal_port, app_services.domain, COALESCE(app_services.static_output,'') AS static_output, app_services.health_check_path, app_services.container_id, app_services.status, app_services.replicas, COALESCE(app_services.cpu_limit, 0) AS cpu_limit, COALESCE(app_services.memory_limit, 0) AS memory_limit, COALESCE(app_services.swap_disabled, FALSE) AS swap_disabled, COALESCE(app_services.deploy_token,'') AS deploy_token, COALESCE(app_services.enable_pr_previews, FALSE) AS enable_pr_previews, COALESCE(app_services.maintenance_mode, FALSE) AS maintenance_mode, app_services.registry_id, app_services.created_at, app_services.updated_at, COALESCE(app_services.optimize, FALSE) AS optimize, COALESCE(app_services.stack, '') AS stack, COALESCE(app_services.stop_grace_period_seconds, 0) AS stop_grace_period_seconds, COALESCE(app_services.runtime_version, '') AS runtime_version, COALESCE(app_services.spa_fallback, TRUE) AS spa_fallback, COALESCE(app_services.build_image, '') AS build_image, COALESCE(app_services.runtime_image, '') AS runtime_image
		FROM app_services JOIN projects ON projects.id = app_services.project_id
		WHERE projects.organization_id = $1 AND ($2 = '' OR app_services.name ILIKE $3 OR app_services.domain ILIKE $3 OR app_services.repository_url ILIKE $3) ORDER BY app_services.created_at DESC`, organizationID, search, pattern,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list app services by organization: %w", err)
	}
	if list == nil {
		list = make([]*AppService, 0)
	}
	return list, nil
}

func (r *AppServiceRepo) ListAll(ctx context.Context, search string) ([]*AppService, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*AppService
	search = strings.TrimSpace(search)
	pattern := "%" + search + "%"
	err := r.db.SelectContext(ctx, &list,
		`SELECT `+appServiceSelectColumns+` FROM app_services WHERE ($1 = '' OR name ILIKE $2 OR domain ILIKE $2 OR repository_url ILIKE $2) ORDER BY created_at ASC`, search, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to list all app services: %w", err)
	}
	if list == nil {
		list = make([]*AppService, 0)
	}
	return list, nil
}

func (r *AppServiceRepo) SearchIDsByName(ctx context.Context, term string, limit int) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var ids []string
	if limit <= 0 {
		limit = 10
	}
	pattern := "%" + strings.TrimSpace(term) + "%"
	err := r.db.SelectContext(ctx, &ids, `SELECT id FROM app_services WHERE name ILIKE $1 LIMIT $2`, pattern, limit)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *AppServiceRepo) Update(_ context.Context, svc *AppService) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	svc.UpdatedAt = time.Now().UTC()
	spaFallback := true
	if svc.SPAFallback != nil {
		spaFallback = *svc.SPAFallback
	}
	_, err := r.db.Exec(
		`UPDATE app_services SET name = $1, repository_url = $2, image_ref = $3, branch = $4, root_directory = $5, icon = $6, runtime_mode = $7, install_command = $8, build_command = $9, start_command = $10, dockerfile_path = $11, build_engine = $12, internal_port = $13, domain = $14, static_output = $15, health_check_path = $16, container_id = $17, status = $18, replicas = $19, cpu_limit = $20, memory_limit = $21, swap_disabled = $22, deploy_token = $23, enable_pr_previews = $24, maintenance_mode = $25, registry_id = $26, updated_at = $27, optimize = $28, stack = $29, stop_grace_period_seconds = $30, runtime_version = $31, spa_fallback = $32, build_image = $33, runtime_image = $34 WHERE id = $35`,
		svc.Name, svc.RepositoryURL, svc.ImageRef, svc.Branch, svc.RootDirectory, svc.Icon, svc.RuntimeMode, svc.InstallCommand, svc.BuildCommand, svc.StartCommand, svc.DockerfilePath, svc.BuildEngine, svc.InternalPort, svc.Domain, svc.StaticOutput, svc.HealthCheckPath, svc.ContainerID, svc.Status, svc.Replicas, svc.CPULimit, svc.MemoryLimit, svc.SwapDisabled, svc.DeployToken, svc.EnablePRPreviews, svc.MaintenanceMode, svc.RegistryID, svc.UpdatedAt, svc.Optimize, svc.Stack, svc.StopGracePeriodSeconds, svc.RuntimeVersion, spaFallback, svc.BuildImage, svc.RuntimeImage, svc.ID,
	)
	return err
}

func (r *AppServiceRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.Exec(`DELETE FROM app_services WHERE id = $1`, id)
	return err
}

func (r *AppServiceRepo) CreateWebhook(_ context.Context, w *Webhook) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if w.ID == "" {
		w.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	w.CreatedAt = now
	w.UpdatedAt = now
	_, err := r.db.Exec(`INSERT INTO service_webhooks (id, service_id, url, secret, events, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		w.ID, w.ServiceID, w.URL, w.Secret, w.Events, w.CreatedAt, w.UpdatedAt)
	return err
}

func (r *AppServiceRepo) ListWebhooksByService(_ context.Context, serviceID string) ([]*Webhook, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*Webhook
	err := r.db.Select(&list, `SELECT id, service_id, url, secret, events, created_at, updated_at FROM service_webhooks WHERE service_id = $1 ORDER BY created_at DESC`, serviceID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *AppServiceRepo) DeleteWebhook(_ context.Context, id, serviceID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.Exec(`DELETE FROM service_webhooks WHERE id = $1 AND service_id = $2`, id, serviceID)
	return err
}

func (r *AppServiceRepo) CreateLogDrain(ctx context.Context, d *LogDrain) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	d.CreatedAt = now
	d.UpdatedAt = now
	if d.ProjectID == "" {
		_ = r.db.QueryRowContext(ctx, `SELECT project_id FROM app_services WHERE id = $1`, d.ServiceID).Scan(&d.ProjectID)
	}
	_, err := r.db.ExecContext(ctx, `INSERT INTO log_drains (id, project_id, service_id, drain_type, endpoint_url, auth_token, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		d.ID, d.ProjectID, d.ServiceID, d.DrainType, d.EndpointURL, d.APIKey, d.CreatedAt, d.UpdatedAt)
	return err
}

func (r *AppServiceRepo) ListLogDrainsByService(ctx context.Context, serviceID string) ([]*LogDrain, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*LogDrain
	err := r.db.SelectContext(ctx, &list, `SELECT id, project_id, service_id, drain_type, endpoint_url, COALESCE(auth_token, '') as api_key, created_at, updated_at FROM log_drains WHERE service_id = $1 ORDER BY created_at DESC`, serviceID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *AppServiceRepo) DeleteLogDrain(ctx context.Context, id, serviceID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `DELETE FROM log_drains WHERE id = $1 AND service_id = $2`, id, serviceID)
	return err
}
