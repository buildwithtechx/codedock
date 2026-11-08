package http

import (
	"context"
	"errors"
	"fmt"
	"time"

	"codedock/internal/models"
	"codedock/internal/repositories"
	"codedock/internal/utils"
)

type engineAdapter struct {
	settingsRepo      repositories.SettingsRepository
	appRepo           repositories.AppServiceRepository
	envRepo           repositories.EnvRepository
	dbRepo            repositories.DatabaseRepository
	projectRepo       repositories.ProjectRepository
	scheduledTaskRepo repositories.ScheduledTaskRepository
	backupRepo        repositories.BackupRepository
	s3Repo            repositories.S3DestinationRepository
	sftpRepo          *repositories.SFTPDestinationRepo
	serviceVarRepo    repositories.ServiceVarRepository
	serverlessRepo    repositories.ServerlessRepository
}

func newEngineAdapter(
	settingsRepo repositories.SettingsRepository,
	appRepo repositories.AppServiceRepository,
	envRepo repositories.EnvRepository,
	dbRepo repositories.DatabaseRepository,
	projectRepo repositories.ProjectRepository,
	scheduledTaskRepo repositories.ScheduledTaskRepository,
	backupRepo repositories.BackupRepository,
	s3Repo repositories.S3DestinationRepository,
	serviceVarRepo repositories.ServiceVarRepository,
	serverlessRepo repositories.ServerlessRepository,
) *engineAdapter {
	return &engineAdapter{
		settingsRepo:      settingsRepo,
		appRepo:           appRepo,
		envRepo:           envRepo,
		dbRepo:            dbRepo,
		projectRepo:       projectRepo,
		scheduledTaskRepo: scheduledTaskRepo,
		backupRepo:        backupRepo,
		s3Repo:            s3Repo,
		serviceVarRepo:    serviceVarRepo,
		serverlessRepo:    serverlessRepo,
	}
}

func (a *engineAdapter) GetServerSettings() (*models.ServerSettings, error) {
	return a.settingsRepo.GetServerSettings(context.Background())
}

func (a *engineAdapter) ListAppServicesByProject(projectID string) ([]*models.AppService, error) {
	return a.appRepo.ListByProject(context.Background(), projectID)
}

func (a *engineAdapter) GetEnvVars(projectID string) (map[string]string, error) {
	return a.envRepo.GetVars(context.Background(), projectID)
}

func (a *engineAdapter) ListServiceVariables(serviceID string) ([]*models.Variable, error) {
	return a.serviceVarRepo.ListByService(context.Background(), serviceID)
}

func (a *engineAdapter) GetServerlessFunctionCode(serviceID string) (*models.ServerlessFunctionCode, error) {
	return a.serverlessRepo.GetCodeByServiceID(context.Background(), serviceID)
}

func (a *engineAdapter) GetDatabase(id string) (*models.Database, error) {
	return a.dbRepo.GetByID(context.Background(), id)
}

func (a *engineAdapter) ListScheduledTasks() ([]models.ScheduledTask, error) {
	return a.scheduledTaskRepo.ListAll(context.Background())
}

func (a *engineAdapter) GetScheduledTask(id string) (*models.ScheduledTask, error) {
	return a.scheduledTaskRepo.GetByID(context.Background(), id)
}

func (a *engineAdapter) GetProject(id string) (*models.ProjectConfig, error) {
	return a.projectRepo.Get(context.Background(), id)
}

func (a *engineAdapter) GetAppService(id string) (*models.AppService, error) {
	return a.appRepo.GetByID(context.Background(), id)
}

func (a *engineAdapter) UpdateScheduledTaskStatusAndOutput(id string, status models.ScheduledTaskStatus, lastRunAt *time.Time, output string) error {
	return a.scheduledTaskRepo.UpdateStatus(context.Background(), id, status, lastRunAt, output)
}

func (a *engineAdapter) ListAllActiveBackupConfigs() ([]*models.BackupConfig, error) {
	return a.backupRepo.ListAllActiveConfigs(context.Background())
}

func (a *engineAdapter) GetBackupConfig(id string) (*models.BackupConfig, error) {
	return a.backupRepo.GetConfigByID(context.Background(), id)
}

func (a *engineAdapter) CreateBackupRecord(rec *models.BackupRecord) error {
	return a.backupRepo.CreateRecord(context.Background(), rec)
}

func (a *engineAdapter) UpdateBackupRecord(opts models.UpdateBackupRecordOpts) error {
	rec, err := a.backupRepo.GetRecordByID(context.Background(), opts.ID)
	if err != nil {
		return err
	}
	if opts.SHA256 != "" {
		rec.SHA256 = opts.SHA256
		rec.VerifiedAt = opts.VerifiedAt
	}
	rec.Status = opts.Status
	rec.FilePath = opts.FilePath
	rec.S3URL = opts.S3URL
	rec.SFTPURL = opts.SFTPURL
	rec.ParentRecordID = opts.ParentRecordID
	if opts.S3DestinationID != "" {
		rec.S3DestinationID = opts.S3DestinationID
	}
	if opts.SFTPDestinationID != "" {
		rec.SFTPDestinationID = opts.SFTPDestinationID
	}
	rec.Logs = opts.Logs
	rec.FileSizeBytes = opts.FileSizeBytes
	rec.CompletedAt = opts.CompletedAt
	return a.backupRepo.UpdateRecord(context.Background(), rec)
}

func (a *engineAdapter) GetBackupRecord(id string) (*models.BackupRecord, error) {
	return a.backupRepo.GetRecordByID(context.Background(), id)
}

func (a *engineAdapter) GetS3Destination(id string) (*models.S3Destination, error) {
	return a.s3Repo.GetS3Destination(context.Background(), id)
}

func (a *engineAdapter) GetSFTPDestination(id string) (*models.SFTPDestination, error) {
	if a.sftpRepo == nil {
		return nil, fmt.Errorf("sftp storage unavailable")
	}
	return a.sftpRepo.Get(context.Background(), id)
}

func (a *engineAdapter) ListBackupRecords(backupConfigID string) ([]*models.BackupRecord, error) {
	return a.backupRepo.ListRecordsByConfig(context.Background(), backupConfigID)
}

func (a *engineAdapter) ClaimRecordExpiry(ctx context.Context, id string) (bool, error) {
	store, ok := a.backupRepo.(interface {
		ClaimRecordExpiry(context.Context, string) (bool, error)
	})
	if !ok {
		return false, fmt.Errorf("backup retention protection storage unavailable")
	}
	return store.ClaimRecordExpiry(ctx, id)
}

func (a *engineAdapter) ControlPlaneSnapshot(ctx context.Context) ([]byte, error) {
	store, ok := a.backupRepo.(interface {
		ControlPlaneSnapshot(context.Context) ([]byte, error)
	})
	if !ok {
		return nil, fmt.Errorf("consistent control-plane snapshots unavailable")
	}
	return store.ControlPlaneSnapshot(ctx)
}

func (a *engineAdapter) VolumeRestoreOwner(ctx context.Context, cfg *models.BackupConfig) (string, error) {
	return a.volumeOwner(ctx, cfg, false)
}

func (a *engineAdapter) VolumeBackupOwner(ctx context.Context, cfg *models.BackupConfig) (string, error) {
	return a.volumeOwner(ctx, cfg, true)
}

func (a *engineAdapter) volumeOwner(ctx context.Context, cfg *models.BackupConfig, allowRemote bool) (string, error) {
	var projectID, containerID string
	if cfg.ServiceID != "" && cfg.DatabaseID != "" {
		return "", errors.New("volume backup must have exactly one owner")
	}
	if cfg.ServiceID != "" {
		app, err := a.appRepo.GetByID(ctx, cfg.ServiceID)
		if err != nil {
			return "", fmt.Errorf("load volume service: %w", err)
		}
		projectID = app.ProjectID
		containerID = utils.NormalizeContainerName(app.ID)
		if app.ContainerID != "" {
			containerID = app.ContainerID
		}
	} else if cfg.DatabaseID != "" && cfg.DatabaseID != "global" {
		db, err := a.dbRepo.GetByID(ctx, cfg.DatabaseID)
		if err != nil {
			return "", fmt.Errorf("load volume database: %w", err)
		}
		projectID = db.ProjectID
		containerID = db.ContainerID
		if containerID == "" {
			containerID = utils.NormalizeContainerName("codedock-db-" + db.Name)
		}
		if cfg.VolumeName != "codedock-db-data-"+db.ID {
			return "", errors.New("volume is not the database data volume")
		}
	} else {
		return "", errors.New("volume backup requires a registered service or database owner")
	}
	project, err := a.projectRepo.Get(ctx, projectID)
	if err != nil {
		return "", fmt.Errorf("load volume project: %w", err)
	}
	if project.ServerID != "" && !allowRemote {
		return "", errors.New("volume restore supports local Docker targets only")
	}
	return containerID, nil
}
