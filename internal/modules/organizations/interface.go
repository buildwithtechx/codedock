package organizations

import "context"

// OrganizationRepository defines persistence for organizations, team members, and roles.
type OrganizationRepository interface {
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
