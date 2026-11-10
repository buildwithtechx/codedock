package users

import "context"

// UserRepository defines persistence for user accounts, credentials, and personal access tokens.
type UserRepository interface {
	CreateUser(ctx context.Context, u *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByStripeCustomerID(ctx context.Context, stripeCustomerID string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUserTOTPSecret(ctx context.Context, userID string) (string, []string, error)
	ListUsers(ctx context.Context, limit, offset int) ([]User, int, error)
	CountUsers(ctx context.Context) (int, error)
	UpdateUser(ctx context.Context, u *User) error
	DeleteUser(ctx context.Context, id string) error
	CreatePAT(ctx context.Context, pat *PersonalAccessToken) error
	ListPATs(ctx context.Context, userID string) ([]*PersonalAccessToken, error)
	DeletePAT(ctx context.Context, id, userID string) error
	GetPATByHash(ctx context.Context, hash string) (*PersonalAccessToken, error)
	FindFoundingAdmin(ctx context.Context) (*User, error)
	GetAppearance(ctx context.Context, userID string) (*AppearancePreferences, error)
	SetAppearance(ctx context.Context, userID string, pref *AppearancePreferences) error
}
