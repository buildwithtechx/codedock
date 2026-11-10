package adapters

import (
	"context"
	"errors"
	"fmt"
	"time"

	"codedock/internal/modules/apps"
	"codedock/internal/modules/apps/tasks"
	"codedock/internal/modules/backups"
	"codedock/internal/modules/databases"
	"codedock/internal/modules/projects"
	"codedock/internal/platform/scheduler"
	"codedock/internal/platform/storage"
)

type EngineAdapter struct {
	appRepo           apps.AppServiceRepository
	dbRepo            databases.Repository
	projectRepo       projects.ProjectRepository
	scheduledTaskRepo tasks.ScheduledTaskRepository
	backupRepo        backups.Repository
	sftpRepo          backups.SFTPRepository
}

func NewEngineAdapter(
	appRepo apps.AppServiceRepository,
	dbRepo databases.Repository,
	projectRepo projects.ProjectRepository,
	scheduledTaskRepo tasks.ScheduledTaskRepository,
	backupRepo backups.Repository,
) *EngineAdapter {
	return &EngineAdapter{
		appRepo:           appRepo,
		dbRepo:            dbRepo,
		projectRepo:       projectRepo,
		scheduledTaskRepo: scheduledTaskRepo,
		backupRepo:        backupRepo,
	}
}

func (a *EngineAdapter) SetSFTPRepo(repo backups.SFTPRepository) {
	a.sftpRepo = repo
}

func (a *EngineAdapter) ListScheduledTasks() ([]scheduler.ScheduledTask, error) {
	tasksList, err := a.scheduledTaskRepo.ListAll(context.Background())
	if err != nil {
		return nil, err
	}
	res := make([]scheduler.ScheduledTask, len(tasksList))
	for i, t := range tasksList {
		var lastRun *time.Time
		if !t.LastRunAt.IsZero() {
			lr := t.LastRunAt
			lastRun = &lr
		}
		res[i] = scheduler.ScheduledTask{
			ID:         t.ID,
			ServiceID:  t.ServiceID,
			Name:       t.Name,
			Schedule:   t.Schedule,
			Command:    t.Command,
			Status:     scheduler.ScheduledTaskStatus(t.Status),
			LastRunAt:  lastRun,
			LastOutput: t.LastOutput,
		}
	}
	return res, nil
}

func (a *EngineAdapter) GetScheduledTask(id string) (*scheduler.ScheduledTask, error) {
	t, err := a.scheduledTaskRepo.GetByID(context.Background(), id)
	if err != nil || t == nil {
		return nil, err
	}
	var lastRun *time.Time
	if !t.LastRunAt.IsZero() {
		lr := t.LastRunAt
		lastRun = &lr
	}
	return &scheduler.ScheduledTask{
		ID:         t.ID,
		ServiceID:  t.ServiceID,
		Name:       t.Name,
		Schedule:   t.Schedule,
		Command:    t.Command,
		Status:     scheduler.ScheduledTaskStatus(t.Status),
		LastRunAt:  lastRun,
		LastOutput: t.LastOutput,
	}, nil
}

func (a *EngineAdapter) GetProject(id string) (*scheduler.ProjectConfig, error) {
	p, err := a.projectRepo.Get(context.Background(), id)
	if err != nil || p == nil {
		return nil, err
	}
	return &scheduler.ProjectConfig{
		ID:   p.ID,
		Name: p.Name,
	}, nil
}

func (a *EngineAdapter) GetAppService(id string) (*scheduler.AppService, error) {
	app, err := a.appRepo.GetByID(context.Background(), id)
	if err != nil || app == nil {
		return nil, err
	}
	return &scheduler.AppService{
		ID:          app.ID,
		ContainerID: app.ContainerID,
		Name:        app.Name,
	}, nil
}

func (a *EngineAdapter) UpdateScheduledTaskStatusAndOutput(id string, status scheduler.ScheduledTaskStatus, lastRunAt *time.Time, output string) error {
	return a.scheduledTaskRepo.UpdateStatus(context.Background(), id, tasks.ScheduledTaskStatus(status), lastRunAt, output)
}

func (a *EngineAdapter) ListAllActiveBackupConfigs() ([]*storage.BackupConfig, error) {
	configs, err := a.backupRepo.ListAllActiveConfigs(context.Background())
	if err != nil {
		return nil, err
	}
	res := make([]*storage.BackupConfig, len(configs))
	for i, c := range configs {
		res[i] = toStorageBackupConfig(c)
	}
	return res, nil
}

func (a *EngineAdapter) GetBackupConfig(id string) (*storage.BackupConfig, error) {
	c, err := a.backupRepo.GetConfigByID(context.Background(), id)
	if err != nil || c == nil {
		return nil, err
	}
	return toStorageBackupConfig(c), nil
}

func (a *EngineAdapter) CreateBackupRecord(rec *storage.BackupRecord) error {
	if rec == nil {
		return errors.New("record is nil")
	}
	modRec := toModuleBackupRecord(rec)
	return a.backupRepo.CreateRecord(context.Background(), modRec)
}

func (a *EngineAdapter) UpdateBackupRecord(opts storage.UpdateBackupRecordOpts) error {
	rec := &backups.BackupRecord{
		ID:                opts.ID,
		Status:            backups.RecordStatus(opts.Status),
		FilePath:          opts.FilePath,
		FileSizeBytes:     opts.FileSizeBytes,
		S3URL:             opts.S3URL,
		SFTPURL:           opts.SFTPURL,
		ParentRecordID:    opts.ParentRecordID,
		S3DestinationID:   opts.S3DestinationID,
		SFTPDestinationID: opts.SFTPDestinationID,
		Logs:              opts.Logs,
		SHA256:            opts.SHA256,
		VerifiedAt:        opts.VerifiedAt,
		CompletedAt:       opts.CompletedAt,
	}
	return a.backupRepo.UpdateRecord(context.Background(), rec)
}

func (a *EngineAdapter) GetBackupRecord(id string) (*storage.BackupRecord, error) {
	rec, err := a.backupRepo.GetRecordByID(context.Background(), id)
	if err != nil || rec == nil {
		return nil, err
	}
	return toStorageBackupRecord(rec), nil
}

func (a *EngineAdapter) GetDatabase(id string) (*storage.Database, error) {
	db, err := a.dbRepo.GetByID(context.Background(), id)
	if err != nil || db == nil {
		return nil, err
	}
	return &storage.Database{
		ID:           db.ID,
		ProjectID:    db.ProjectID,
		Name:         db.Name,
		Engine:       storage.DatabaseEngine(db.Engine),
		Version:      db.Version,
		ContainerID:  db.ContainerID,
		Status:       storage.DatabaseStatus(db.Status),
		InternalPort: db.Port,
		DatabaseName: db.DatabaseName,
		User:         db.Username,
		Username:     db.Username,
		Password:     db.Password,
	}, nil
}

func (a *EngineAdapter) GetS3Destination(id string) (*storage.S3Destination, error) {
	dest, err := a.backupRepo.GetS3Destination(context.Background(), id)
	if err != nil || dest == nil {
		return nil, err
	}
	return &storage.S3Destination{
		ID:              dest.ID,
		Name:            dest.Name,
		Description:     dest.Description,
		Provider:        dest.Provider,
		Endpoint:        dest.Endpoint,
		Bucket:          dest.Bucket,
		Region:          dest.Region,
		PathPrefix:      dest.PathPrefix,
		IsDefault:       dest.IsDefault,
		LastVerifiedAt:  dest.LastVerifiedAt,
		LastVerifyError: dest.LastVerifyError,
		AccessKeyID:     dest.AccessKeyID,
		SecretAccessKey: dest.SecretAccessKey,
		CreatedAt:       dest.CreatedAt,
	}, nil
}

func (a *EngineAdapter) GetSFTPDestination(id string) (*storage.SFTPDestination, error) {
	if a.sftpRepo == nil {
		return nil, fmt.Errorf("sftp storage unavailable")
	}
	dest, err := a.sftpRepo.Get(context.Background(), id)
	if err != nil || dest == nil {
		return nil, err
	}
	return &storage.SFTPDestination{
		ID:              dest.ID,
		OrganizationID:  dest.OrganizationID,
		ProjectID:       dest.ProjectID,
		Name:            dest.Name,
		Description:     dest.Description,
		Host:            dest.Host,
		Port:            dest.Port,
		Username:        dest.Username,
		Password:        dest.Password,
		PrivateKey:      dest.PrivateKey,
		PathPrefix:      dest.PathPrefix,
		LastVerifiedAt:  dest.LastVerifiedAt,
		LastVerifyError: dest.LastVerifyError,
		CreatedAt:       dest.CreatedAt,
	}, nil
}

func (a *EngineAdapter) ListBackupRecords(backupConfigID string) ([]*storage.BackupRecord, error) {
	records, err := a.backupRepo.ListRecordsByConfig(context.Background(), backupConfigID)
	if err != nil {
		return nil, err
	}
	res := make([]*storage.BackupRecord, len(records))
	for i, r := range records {
		res[i] = toStorageBackupRecord(r)
	}
	return res, nil
}
