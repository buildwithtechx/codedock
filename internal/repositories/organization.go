package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"codedock.run/codedock/internal/models"
)

type OrganizationRepository interface {
	Create(ctx context.Context, org *models.Organization) error
	GetByID(ctx context.Context, id string) (*models.Organization, error)
	ListByUser(ctx context.Context, userID string) ([]*models.Organization, error)
	Update(ctx context.Context, org *models.Organization) error
	Delete(ctx context.Context, id string) error

	AddMember(ctx context.Context, member *models.OrganizationMember) error
	GetMember(ctx context.Context, orgID, userID string) (*models.OrganizationMember, error)
	GetMemberByID(ctx context.Context, id string) (*models.OrganizationMember, error)
	GetMemberByEmail(ctx context.Context, orgID, email string) (*models.OrganizationMember, error)
	ListMembers(ctx context.Context, orgID string) ([]*models.OrganizationMember, error)
	UpdateMember(ctx context.Context, member *models.OrganizationMember) error
	RemoveMember(ctx context.Context, id string) error
	ListInvitesByEmail(ctx context.Context, email string) ([]*models.OrganizationMember, error)
}

type organizationRepository struct {
	db *sqlx.DB
}

const organizationColumns = `id, name, created_at, updated_at`
const organizationMemberColumns = `id, organization_id, user_id, email, permission, status, invited_at, accepted_at`

func NewOrganizationRepository(db *sql.DB) OrganizationRepository {
	return &organizationRepository{db: sqlx.NewDb(db, "sqlite")}
}

func (r *organizationRepository) Create(ctx context.Context, org *models.Organization) error {
	query := `
		INSERT INTO organizations (id, name, created_at, updated_at)
		VALUES (:id, :name, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, org)
	return err
}

func (r *organizationRepository) GetByID(ctx context.Context, id string) (*models.Organization, error) {
	query := `SELECT ` + organizationColumns + ` FROM organizations WHERE id = ?`
	var org models.Organization
	err := r.db.GetContext(ctx, &org, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &org, err
}

func (r *organizationRepository) ListByUser(ctx context.Context, userID string) ([]*models.Organization, error) {
	query := `
		SELECT o.id, o.name, o.created_at, o.updated_at FROM organizations o
		JOIN organization_members om ON o.id = om.organization_id
		WHERE om.user_id = ?
		ORDER BY o.name ASC
	`
	var orgs []*models.Organization
	err := r.db.SelectContext(ctx, &orgs, query, userID)
	if orgs == nil {
		orgs = make([]*models.Organization, 0)
	}
	return orgs, err
}

func (r *organizationRepository) Update(ctx context.Context, org *models.Organization) error {
	org.UpdatedAt = time.Now()
	query := `
		UPDATE organizations
		SET name = :name, updated_at = :updated_at
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, org)
	return err
}

func (r *organizationRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM organizations WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *organizationRepository) AddMember(ctx context.Context, member *models.OrganizationMember) error {
	query := `
		INSERT INTO organization_members (id, organization_id, user_id, email, permission, status, invited_at, accepted_at)
		VALUES (:id, :organization_id, :user_id, :email, :permission, :status, :invited_at, :accepted_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, member)
	return err
}

func (r *organizationRepository) GetMember(ctx context.Context, orgID, userID string) (*models.OrganizationMember, error) {
	query := `SELECT ` + organizationMemberColumns + ` FROM organization_members WHERE organization_id = ? AND user_id = ?`
	var member models.OrganizationMember
	err := r.db.GetContext(ctx, &member, query, orgID, userID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &member, err
}

func (r *organizationRepository) GetMemberByEmail(ctx context.Context, orgID, email string) (*models.OrganizationMember, error) {
	query := `SELECT ` + organizationMemberColumns + ` FROM organization_members WHERE organization_id = ? AND email = ?`
	var member models.OrganizationMember
	err := r.db.GetContext(ctx, &member, query, orgID, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &member, err
}

func (r *organizationRepository) ListMembers(ctx context.Context, orgID string) ([]*models.OrganizationMember, error) {
	query := `SELECT ` + organizationMemberColumns + ` FROM organization_members WHERE organization_id = ? ORDER BY invited_at DESC`
	var members []*models.OrganizationMember
	err := r.db.SelectContext(ctx, &members, query, orgID)
	if members == nil {
		members = make([]*models.OrganizationMember, 0)
	}
	return members, err
}

func (r *organizationRepository) UpdateMember(ctx context.Context, member *models.OrganizationMember) error {
	query := `
		UPDATE organization_members
		SET user_id = :user_id, permission = :permission, status = :status, accepted_at = :accepted_at
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, member)
	return err
}

func (r *organizationRepository) GetMemberByID(ctx context.Context, id string) (*models.OrganizationMember, error) {
	query := `SELECT ` + organizationMemberColumns + ` FROM organization_members WHERE id = ?`
	var member models.OrganizationMember
	err := r.db.GetContext(ctx, &member, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &member, err
}

func (r *organizationRepository) RemoveMember(ctx context.Context, id string) error {
	query := `DELETE FROM organization_members WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *organizationRepository) ListInvitesByEmail(ctx context.Context, email string) ([]*models.OrganizationMember, error) {
	query := `
		SELECT ` + organizationMemberColumns + ` FROM organization_members
		WHERE email = ? AND status = ?
	`
	var members []*models.OrganizationMember
	err := r.db.SelectContext(ctx, &members, query, email, models.MemberStatusPending)
	if members == nil {
		members = make([]*models.OrganizationMember, 0)
	}
	return members, err
}
