package organizations

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type OrganizationRepository interface {
	Create(ctx context.Context, org *Organization) error
	CreateWithOwner(ctx context.Context, org *Organization, owner *OrganizationMember) error
	GetByID(ctx context.Context, id string) (*Organization, error)
	ListByUser(ctx context.Context, userID string) ([]*Organization, error)
	Update(ctx context.Context, org *Organization) error
	Delete(ctx context.Context, id string) error

	AddMember(ctx context.Context, member *OrganizationMember) error
	GetMember(ctx context.Context, orgID, userID string) (*OrganizationMember, error)
	GetMemberByID(ctx context.Context, id string) (*OrganizationMember, error)
	GetMemberByEmail(ctx context.Context, orgID, email string) (*OrganizationMember, error)
	ListMembers(ctx context.Context, orgID string) ([]*OrganizationMember, error)
	UpdateMember(ctx context.Context, member *OrganizationMember) error
	RemoveMember(ctx context.Context, id string) error
	ListInvitesByEmail(ctx context.Context, email string) ([]*OrganizationMember, error)
	GetUserRole(ctx context.Context, userID string) (string, error)
}

type organizationRepository struct {
	db *sqlx.DB
}

const organizationColumns = `id, name, created_at, updated_at`
const organizationMemberColumns = `id, organization_id, user_id, email, permission, status, invited_at, accepted_at`

func NewOrganizationRepository(db *sql.DB) OrganizationRepository {
	return &organizationRepository{db: sqlx.NewDb(db, "pgx")}
}

func (r *organizationRepository) Create(ctx context.Context, org *Organization) error {
	query := `
		INSERT INTO organizations (id, name, created_at, updated_at)
		VALUES (:id, :name, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, org)
	return err
}

func (r *organizationRepository) CreateWithOwner(ctx context.Context, org *Organization, owner *OrganizationMember) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.NamedExecContext(ctx, `
		INSERT INTO organizations (id, name, created_at, updated_at)
		VALUES (:id, :name, :created_at, :updated_at)
	`, org); err != nil {
		return err
	}
	if _, err := tx.NamedExecContext(ctx, `
		INSERT INTO organization_members (id, organization_id, user_id, email, permission, status, invited_at, accepted_at)
		VALUES (:id, :organization_id, :user_id, :email, :permission, :status, :invited_at, :accepted_at)
	`, owner); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *organizationRepository) GetByID(ctx context.Context, id string) (*Organization, error) {
	query := `SELECT ` + organizationColumns + ` FROM organizations WHERE id = $1`
	var org Organization
	err := r.db.GetContext(ctx, &org, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &org, err
}

func (r *organizationRepository) ListByUser(ctx context.Context, userID string) ([]*Organization, error) {
	query := `
		SELECT o.id, o.name, o.created_at, o.updated_at FROM organizations o
		JOIN organization_members om ON o.id = om.organization_id
		WHERE om.user_id = $1
		ORDER BY o.name ASC
	`
	var orgs []*Organization
	err := r.db.SelectContext(ctx, &orgs, query, userID)
	if orgs == nil {
		orgs = make([]*Organization, 0)
	}
	return orgs, err
}

func (r *organizationRepository) Update(ctx context.Context, org *Organization) error {
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
	query := `DELETE FROM organizations WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *organizationRepository) AddMember(ctx context.Context, member *OrganizationMember) error {
	query := `
		INSERT INTO organization_members (id, organization_id, user_id, email, permission, status, invited_at, accepted_at)
		VALUES (:id, :organization_id, :user_id, :email, :permission, :status, :invited_at, :accepted_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, member)
	return err
}

func (r *organizationRepository) GetMember(ctx context.Context, orgID, userID string) (*OrganizationMember, error) {
	query := `SELECT ` + organizationMemberColumns + ` FROM organization_members WHERE organization_id = $1 AND user_id = $2`
	var member OrganizationMember
	err := r.db.GetContext(ctx, &member, query, orgID, userID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &member, err
}

func (r *organizationRepository) GetMemberByEmail(ctx context.Context, orgID, email string) (*OrganizationMember, error) {
	query := `SELECT ` + organizationMemberColumns + ` FROM organization_members WHERE organization_id = $1 AND email = $2`
	var member OrganizationMember
	err := r.db.GetContext(ctx, &member, query, orgID, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &member, err
}

func (r *organizationRepository) ListMembers(ctx context.Context, orgID string) ([]*OrganizationMember, error) {
	query := `SELECT ` + organizationMemberColumns + ` FROM organization_members WHERE organization_id = $1 ORDER BY invited_at DESC`
	var members []*OrganizationMember
	err := r.db.SelectContext(ctx, &members, query, orgID)
	if members == nil {
		members = make([]*OrganizationMember, 0)
	}
	return members, err
}

func (r *organizationRepository) UpdateMember(ctx context.Context, member *OrganizationMember) error {
	query := `
		UPDATE organization_members
		SET user_id = :user_id, permission = :permission, status = :status, accepted_at = :accepted_at
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, member)
	return err
}

func (r *organizationRepository) GetMemberByID(ctx context.Context, id string) (*OrganizationMember, error) {
	query := `SELECT ` + organizationMemberColumns + ` FROM organization_members WHERE id = $1`
	var member OrganizationMember
	err := r.db.GetContext(ctx, &member, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &member, err
}

func (r *organizationRepository) RemoveMember(ctx context.Context, id string) error {
	query := `DELETE FROM organization_members WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *organizationRepository) ListInvitesByEmail(ctx context.Context, email string) ([]*OrganizationMember, error) {
	query := `
		SELECT ` + organizationMemberColumns + ` FROM organization_members
		WHERE email = $1 AND status = $2
	`
	var members []*OrganizationMember
	err := r.db.SelectContext(ctx, &members, query, email, MemberStatusPending)
	if members == nil {
		members = make([]*OrganizationMember, 0)
	}
	return members, err
}

func (r *organizationRepository) GetUserRole(ctx context.Context, userID string) (string, error) {
	var role string
	err := r.db.GetContext(ctx, &role, `SELECT role FROM users WHERE id = $1`, userID)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return role, err
}
