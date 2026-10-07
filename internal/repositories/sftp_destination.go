package repositories

import (
	"codedock.run/codedock/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SFTPDestinationRepo struct {
	db *sqlx.DB
}

func NewSFTPDestinationRepo(db *sql.DB) *SFTPDestinationRepo {
	return &SFTPDestinationRepo{db: sqlx.NewDb(db, "sqlite")}
}

func (r *SFTPDestinationRepo) Create(ctx context.Context, dest *models.SFTPDestination) error {
	if dest.ID == "" {
		dest.ID = uuid.NewString()
	}
	if dest.Port <= 0 {
		dest.Port = 22
	}
	dest.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	_, err := r.db.ExecContext(ctx, `INSERT INTO sftp_destinations(id,organization_id,project_id,name,description,host,port,username,password,private_key,path_prefix,created_at) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)`, dest.ID, dest.OrganizationID, dest.ProjectID, dest.Name, dest.Description, dest.Host, dest.Port, dest.Username, dest.Password, dest.PrivateKey, dest.PathPrefix, dest.CreatedAt)
	return err
}

func (r *SFTPDestinationRepo) Get(ctx context.Context, id string) (*models.SFTPDestination, error) {
	var dest models.SFTPDestination
	if err := r.db.GetContext(ctx, &dest, `SELECT id,COALESCE(organization_id,'') AS organization_id,COALESCE(project_id,'') AS project_id,name,COALESCE(description,'') AS description,host,port,username,COALESCE(password,'') AS password,COALESCE(private_key,'') AS private_key,COALESCE(path_prefix,'') AS path_prefix,last_verified_at,COALESCE(last_verify_error,'') AS last_verify_error,created_at FROM sftp_destinations WHERE id=?`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("sftp destination not found")
		}
		return nil, err
	}
	return &dest, nil
}

func (r *SFTPDestinationRepo) ListByProject(ctx context.Context, projectID string) ([]models.SFTPDestination, error) {
	result := []models.SFTPDestination{}
	if err := r.db.SelectContext(ctx, &result, `SELECT id,COALESCE(organization_id,'') AS organization_id,COALESCE(project_id,'') AS project_id,name,COALESCE(description,'') AS description,host,port,username,'' AS password,'' AS private_key,COALESCE(path_prefix,'') AS path_prefix,last_verified_at,COALESCE(last_verify_error,'') AS last_verify_error,created_at FROM sftp_destinations WHERE project_id=? OR organization_id IN (SELECT organization_id FROM projects WHERE id=?) ORDER BY name`, projectID, projectID); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *SFTPDestinationRepo) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM sftp_destinations WHERE id=? AND NOT EXISTS(SELECT 1 FROM backup_configs WHERE sftp_destination_id=?)`, id, id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("destination is referenced or missing")
	}
	return nil
}
