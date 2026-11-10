package databases

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"

	"codedock/internal/utils"
)

type Service struct {
	repo     Repository
	deployer Deployer
}

type DatabaseService = Service

func NewService(repo Repository, deployer Deployer) *Service {
	return &Service{repo: repo, deployer: deployer}
}

func NewDatabaseService(repo Repository, deployer Deployer) *Service {
	return NewService(repo, deployer)
}

func generateRandomPassword(length int) string {
	b := make([]byte, length/2)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func defaultPortForEngine(engine DatabaseEngine) int {
	switch strings.ToLower(string(engine)) {
	case "postgres", "postgresql":
		return 5432
	case "mysql", "mariadb":
		return 3306
	case "redis":
		return 6379
	case "mongodb":
		return 27017
	case "clickhouse":
		return 8123
	default:
		return 5432
	}
}

func (s *Service) VerifyProjectOwnership(ctx context.Context, projectID, userID, userRole string) error {
	if s.repo == nil {
		return nil
	}
	exists, err := s.repo.ProjectExists(ctx, projectID)
	if err != nil || !exists {
		return ErrProjectNotFound
	}
	if !s.repo.IsProjectMemberOrOwner(ctx, projectID, userID, userRole) {
		return ErrAccessDenied
	}
	return nil
}

func (s *Service) IsProjectMemberOrOwner(ctx context.Context, projectID, userID, userRole string) bool {
	if s.repo == nil {
		return true
	}
	return s.repo.IsProjectMemberOrOwner(ctx, projectID, userID, userRole)
}

func (s *Service) LogAction(ctx context.Context, entry AuditLogEntry) {
	if s.repo != nil {
		s.repo.LogAction(ctx, entry)
	}
}

func (s *Service) CreateDatabaseFromRequest(ctx context.Context, req *CreateDatabaseRequest) (*Database, error) {
	if req.Name == "" {
		req.Name = utils.GenerateRandomName()
	}
	db := &Database{
		ProjectID:          req.ProjectID,
		EnvironmentID:      req.EnvironmentID,
		Name:               req.Name,
		Engine:             req.Engine,
		Version:            req.Version,
		Username:           req.Username,
		Password:           req.Password,
		DatabaseName:       req.DatabaseName,
		Port:               req.Port,
		VolumePath:         req.VolumePath,
		CustomArgs:         req.CustomArgs,
		LogicalReplication: req.LogicalReplication,
	}
	return s.CreateDatabase(ctx, db)
}

func (s *Service) CreateDatabase(ctx context.Context, db *Database) (*Database, error) {
	if db.Name == "" {
		db.Name = utils.GenerateRandomName()
	}
	db.ID = uuid.NewString()
	db.Status = DatabaseStatusCreated
	db.CreatedAt = time.Now()
	db.UpdatedAt = time.Now()

	if db.Username == "" {
		db.Username = "codedock"
	}
	if db.Password == "" {
		db.Password = generateRandomPassword(24)
	}
	if db.DatabaseName == "" {
		db.DatabaseName = strings.ReplaceAll(strings.ToLower(db.Name), "-", "_")
	}
	if db.Port == 0 {
		db.Port = defaultPortForEngine(db.Engine)
	}
	if db.VolumePath == "" {
		home, _ := os.UserHomeDir()
		if home == "" {
			home = "/tmp"
		}
		db.VolumePath = fmt.Sprintf("%s/.codedock/volumes/db_%s", home, db.ID)
	}
	if db.InternalHost == "" {
		db.InternalHost = fmt.Sprintf("codedock_db_%s", strings.ToLower(db.Name))
	}

	if err := s.repo.Create(ctx, db); err != nil {
		return nil, err
	}

	if s.deployer != nil {
		dbCopy := *db
		go func(target Database) {
			bgCtx := context.Background()
			containerID, err := s.deployer.SpinUp(bgCtx, &target)
			if err != nil {
				slog.Error("failed to spin up database container", "error", err)
				target.Status = DatabaseStatusError
				_ = s.repo.Update(bgCtx, &target)
				return
			}
			target.ContainerID = containerID
			target.Status = DatabaseStatusRunning
			target.UpdatedAt = time.Now()
			_ = s.repo.Update(bgCtx, &target)
		}(dbCopy)
	}

	return db, nil
}

func (s *Service) GetDatabase(ctx context.Context, id string) (*Database, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) ListDatabases(ctx context.Context, search string) ([]*Database, error) {
	return s.repo.ListAll(ctx, search)
}

func (s *Service) ListByProject(ctx context.Context, projectID string) ([]*Database, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *Service) UpdateDatabase(ctx context.Context, db *Database) error {
	db.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, db); err != nil {
		return err
	}
	if s.deployer != nil {
		go func() {
			bgCtx := context.Background()
			if _, err := s.deployer.SpinUp(bgCtx, db); err != nil {
				slog.Error("failed to restart database container after update", "error", err)
			}
		}()
	}
	return nil
}

func (s *Service) DeleteDatabase(ctx context.Context, id string) error {
	if s.deployer != nil {
		return s.deployer.Remove(ctx, id, func() error {
			return s.repo.Delete(ctx, id)
		})
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) StartDatabase(ctx context.Context, id string) (*Database, error) {
	db, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if s.deployer != nil {
		if _, err := s.deployer.SpinUp(ctx, db); err != nil {
			return nil, err
		}
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) StopDatabase(ctx context.Context, id string) error {
	if s.deployer == nil {
		return ErrDeployerNotConfigured
	}
	return s.deployer.Stop(ctx, id)
}

func (s *Service) ImportData(ctx context.Context, id, sourceURL string) error {
	db, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if s.deployer == nil {
		return ErrDeployerNotConfigured
	}
	return s.deployer.ImportData(ctx, db, sourceURL)
}
