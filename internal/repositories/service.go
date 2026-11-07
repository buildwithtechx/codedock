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
	ListByOrganization(ctx context.Context, organizationID string, search string) ([]*models.AppService, error)
	ListAll(ctx context.Context, search string) ([]*models.AppService, error)
	SearchIDsByName(ctx context.Context, term string, limit int) ([]string, error)
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
		`INSERT INTO app_services (id, project_id, environment_id, name, repository_url, git_user_id, image_ref, branch, root_directory, icon, runtime_mode, install_command, build_command, start_command, dockerfile_path, build_engine, internal_port, domain, static_output, health_check_path, container_id, status, replicas, cpu_limit, memory_limit, swap_disabled, deploy_token, enable_pr_previews, maintenance_mode, registry_id, created_at, updated_at, optimize, stack, stop_grace_period_seconds, runtime_version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36)`,
		svc.ID, svc.ProjectID, svc.EnvironmentID, svc.Name, svc.RepositoryURL, svc.GitUserID, svc.ImageRef, svc.Branch,
		svc.RootDirectory, svc.Icon, svc.RuntimeMode, svc.InstallCommand, svc.BuildCommand, svc.StartCommand, svc.DockerfilePath, svc.BuildEngine,
		svc.InternalPort, svc.Domain, svc.StaticOutput, svc.HealthCheckPath, svc.ContainerID, svc.Status, svc.Replicas, svc.CPULimit, svc.MemoryLimit, svc.SwapDisabled, svc.DeployToken, svc.EnablePRPreviews, svc.MaintenanceMode, svc.RegistryID, svc.CreatedAt, svc.UpdatedAt, svc.Optimize, svc.Stack, svc.StopGracePeriodSeconds, svc.RuntimeVersion,
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
		`SELECT COALESCE((SELECT app_id FROM projects WHERE projects.id = app_services.project_id), '') AS app_id, id, project_id, environment_id, name, repository_url, git_user_id, COALESCE(image_ref,'') AS image_ref, branch, root_directory, COALESCE(icon,'git') AS icon, runtime_mode, COALESCE(install_command,'') AS install_command, build_command, start_command, dockerfile_path, build_engine, internal_port, domain, COALESCE(static_output,'') AS static_output, health_check_path, container_id, status, replicas, COALESCE(cpu_limit, 0) AS cpu_limit, COALESCE(memory_limit, 0) AS memory_limit, COALESCE(swap_disabled, FALSE) AS swap_disabled, COALESCE(deploy_token,'') AS deploy_token, COALESCE(enable_pr_previews, FALSE) AS enable_pr_previews, COALESCE(maintenance_mode, FALSE) AS maintenance_mode, registry_id, created_at, updated_at, COALESCE(optimize, FALSE) AS optimize, COALESCE(stack, '') AS stack, COALESCE(stop_grace_period_seconds, 0) AS stop_grace_period_seconds, COALESCE(runtime_version, '') AS runtime_version
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
		`SELECT COALESCE((SELECT app_id FROM projects WHERE projects.id = app_services.project_id), '') AS app_id, id, project_id, environment_id, name, repository_url, git_user_id, COALESCE(image_ref,'') AS image_ref, branch, root_directory, COALESCE(icon,'git') AS icon, runtime_mode, COALESCE(install_command,'') AS install_command, build_command, start_command, dockerfile_path, build_engine, internal_port, domain, COALESCE(static_output,'') AS static_output, health_check_path, container_id, status, replicas, COALESCE(cpu_limit, 0) AS cpu_limit, COALESCE(memory_limit, 0) AS memory_limit, COALESCE(swap_disabled, FALSE) AS swap_disabled, COALESCE(deploy_token,'') AS deploy_token, COALESCE(enable_pr_previews, FALSE) AS enable_pr_previews, COALESCE(maintenance_mode, FALSE) AS maintenance_mode, registry_id, created_at, updated_at, COALESCE(optimize, FALSE) AS optimize, COALESCE(stack, '') AS stack, COALESCE(stop_grace_period_seconds, 0) AS stop_grace_period_seconds, COALESCE(runtime_version, '') AS runtime_version
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
		`SELECT COALESCE((SELECT app_id FROM projects WHERE projects.id = app_services.project_id), '') AS app_id, id, project_id, environment_id, name, repository_url, git_user_id, COALESCE(image_ref,'') AS image_ref, branch, root_directory, COALESCE(icon,'git') AS icon, runtime_mode, COALESCE(install_command,'') AS install_command, build_command, start_command, dockerfile_path, build_engine, internal_port, domain, COALESCE(static_output,'') AS static_output, health_check_path, container_id, status, replicas, COALESCE(cpu_limit, 0) AS cpu_limit, COALESCE(memory_limit, 0) AS memory_limit, COALESCE(swap_disabled, FALSE) AS swap_disabled, COALESCE(deploy_token,'') AS deploy_token, COALESCE(enable_pr_previews, FALSE) AS enable_pr_previews, COALESCE(maintenance_mode, FALSE) AS maintenance_mode, registry_id, created_at, updated_at, COALESCE(optimize, FALSE) AS optimize, COALESCE(stack, '') AS stack, COALESCE(stop_grace_period_seconds, 0) AS stop_grace_period_seconds, COALESCE(runtime_version, '') AS runtime_version
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

func (r *AppServiceRepo) ListByOrganization(ctx context.Context, organizationID string, search string) ([]*models.AppService, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*models.AppService
	search = strings.TrimSpace(search)
	pattern := ContainsLikePattern(search)
	err := r.db.SelectContext(ctx, &list,
		`SELECT COALESCE(projects.app_id, '') AS app_id, app_services.id, app_services.project_id, app_services.environment_id, app_services.name, app_services.repository_url, app_services.git_user_id, COALESCE(app_services.image_ref,'') AS image_ref, app_services.branch, app_services.root_directory, COALESCE(app_services.icon,'git') AS icon, app_services.runtime_mode, COALESCE(app_services.install_command,'') AS install_command, app_services.build_command, app_services.start_command, app_services.dockerfile_path, app_services.build_engine, app_services.internal_port, app_services.domain, COALESCE(app_services.static_output,'') AS static_output, app_services.health_check_path, app_services.container_id, app_services.status, app_services.replicas, COALESCE(app_services.cpu_limit, 0) AS cpu_limit, COALESCE(app_services.memory_limit, 0) AS memory_limit, COALESCE(app_services.swap_disabled, FALSE) AS swap_disabled, COALESCE(app_services.deploy_token,'') AS deploy_token, COALESCE(app_services.enable_pr_previews, FALSE) AS enable_pr_previews, COALESCE(app_services.maintenance_mode, FALSE) AS maintenance_mode, app_services.registry_id, app_services.created_at, app_services.updated_at, COALESCE(app_services.optimize, FALSE) AS optimize, COALESCE(app_services.stack, '') AS stack, COALESCE(app_services.stop_grace_period_seconds, 0) AS stop_grace_period_seconds, COALESCE(app_services.runtime_version, '') AS runtime_version
		FROM app_services JOIN projects ON projects.id = app_services.project_id
		WHERE projects.organization_id = $1 AND ($2 = '' OR app_services.name ILIKE $3 OR app_services.domain ILIKE $3 OR app_services.repository_url ILIKE $3) ORDER BY app_services.created_at DESC`, organizationID, search, pattern,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list app services by organization: %w", err)
	}
	if list == nil {
		list = make([]*models.AppService, 0)
	}
	return list, nil
}

func (r *AppServiceRepo) ListAll(ctx context.Context, search string) ([]*models.AppService, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*models.AppService
	search = strings.TrimSpace(search)
	pattern := ContainsLikePattern(search)
	err := r.db.SelectContext(ctx, &list,
		`SELECT COALESCE((SELECT app_id FROM projects WHERE projects.id = app_services.project_id), '') AS app_id, id, project_id, environment_id, name, repository_url, git_user_id, COALESCE(image_ref,'') AS image_ref, branch, root_directory, COALESCE(icon,'git') AS icon, runtime_mode, COALESCE(install_command,'') AS install_command, build_command, start_command, dockerfile_path, build_engine, internal_port, domain, COALESCE(static_output,'') AS static_output, health_check_path, container_id, status, replicas, COALESCE(cpu_limit, 0) AS cpu_limit, COALESCE(memory_limit, 0) AS memory_limit, COALESCE(swap_disabled, FALSE) AS swap_disabled, COALESCE(deploy_token,'') AS deploy_token, COALESCE(enable_pr_previews, FALSE) AS enable_pr_previews, COALESCE(maintenance_mode, FALSE) AS maintenance_mode, registry_id, created_at, updated_at, COALESCE(optimize, FALSE) AS optimize, COALESCE(stack, '') AS stack, COALESCE(stop_grace_period_seconds, 0) AS stop_grace_period_seconds, COALESCE(runtime_version, '') AS runtime_version
		FROM app_services WHERE ($1 = '' OR name ILIKE $2 OR domain ILIKE $2 OR repository_url ILIKE $2) ORDER BY created_at ASC`, search, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to list all app services: %w", err)
	}
	if list == nil {
		list = make([]*models.AppService, 0)
	}
	return list, nil
}

func (r *AppServiceRepo) SearchIDsByName(ctx context.Context, term string, limit int) ([]string, error) {
	if limit <= 0 {
		limit = 200
	}
	pattern := ContainsLikePattern(term)
	ids := make([]string, 0)
	err := r.db.SelectContext(ctx, &ids, `SELECT id FROM app_services WHERE name ILIKE $1 LIMIT $2`, pattern, limit)
	if err != nil {
		return nil, fmt.Errorf("search service ids by name: %w", err)
	}
	return ids, nil
}

func (r *AppServiceRepo) Update(_ context.Context, svc *models.AppService) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	svc.UpdatedAt = time.Now().UTC()
	_, err := r.db.Exec(
		`UPDATE app_services SET
		name = $1, repository_url = $2, git_user_id = $3, image_ref = $4, branch = $5, root_directory = $6, icon = $7, runtime_mode = $8,
		install_command = $9, build_command = $10, start_command = $11, dockerfile_path = $12, build_engine = $13,
		internal_port = $14, domain = $15, static_output = $16, health_check_path = $17, container_id = $18, status = $19, replicas = $20, cpu_limit = $21, memory_limit = $22, swap_disabled = $23, deploy_token = $24, enable_pr_previews = $25, maintenance_mode = $26, registry_id = $27, optimize = $28, stack = $29, stop_grace_period_seconds = $30, runtime_version = $31, updated_at = $32
		WHERE id = $33`,
		svc.Name, svc.RepositoryURL, svc.GitUserID, svc.ImageRef, svc.Branch, svc.RootDirectory, svc.Icon, svc.RuntimeMode,
		svc.InstallCommand, svc.BuildCommand, svc.StartCommand, svc.DockerfilePath, svc.BuildEngine,
		svc.InternalPort, svc.Domain, svc.StaticOutput, svc.HealthCheckPath, svc.ContainerID, svc.Status, svc.Replicas, svc.CPULimit, svc.MemoryLimit, svc.SwapDisabled, svc.DeployToken, svc.EnablePRPreviews, svc.MaintenanceMode, svc.RegistryID, svc.Optimize, svc.Stack, svc.StopGracePeriodSeconds, svc.RuntimeVersion, svc.UpdatedAt, svc.ID,
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
