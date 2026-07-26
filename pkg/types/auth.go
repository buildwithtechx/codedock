package types

import "time"

type UserRole string

const (
	UserRoleOwner  UserRole = "owner"
	UserRoleAdmin  UserRole = "admin"
	UserRoleMember UserRole = "member"
)

type User struct {
	ID                   string     `json:"id" db:"id"`
	Email                string     `json:"email" db:"email"`
	Name                 string     `json:"name" db:"name"`
	PasswordHash         string     `json:"-" db:"password_hash"`
	Role                 UserRole   `json:"role" db:"role"`
	IsActive             bool       `json:"isActive" db:"is_active"`
	EmailVerified        bool       `json:"emailVerified" db:"email_verified"`
	TOTPEnabled          bool       `json:"totpEnabled" db:"totp_enabled"`
	OAuthProvider        string     `json:"oauthProvider,omitempty" db:"oauth_provider"`
	CreatedAt            time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt            time.Time  `json:"updatedAt" db:"updated_at"`
	LastLogin            *time.Time `json:"lastLogin,omitempty" db:"last_login"`
	ProjectsCount        int        `json:"projectsCount" db:"projects_count"`
	ServicesCount        int        `json:"servicesCount" db:"services_count"`
	APIKeysCount         int        `json:"apiKeysCount" db:"api_keys_count"`
	PlanType             string     `json:"planType" db:"plan_type"`
	StripeCustomerID     *string    `json:"stripeCustomerId,omitempty" db:"stripe_customer_id"`
	StripeSubscriptionID *string    `json:"stripeSubscriptionId,omitempty" db:"stripe_subscription_id"`
	StripePriceID        *string    `json:"stripePriceId,omitempty" db:"stripe_price_id"`
}

type AuthResult struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type SigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupRequest struct {
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Role     UserRole `json:"role"`
}
