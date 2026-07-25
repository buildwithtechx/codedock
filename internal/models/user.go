package models

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

type UserClaims struct {
	UserID      string   `json:"sub"`
	Email       string   `json:"email"`
	Role        UserRole `json:"role"`
	TOTPEnabled bool     `json:"totpEnabled"`
	PlanType    string   `json:"planType"`
}

type PersonalAccessToken struct {
	ID              string    `json:"id" db:"id"`
	UserID          string    `json:"userId" db:"user_id"`
	Name            string    `json:"name" db:"name"`
	TokenHash       string    `json:"-" db:"token_hash"`
	Prefix          string    `json:"prefix" db:"prefix"`
	AccessLevel     string    `json:"accessLevel" db:"access_level"`
	ProjectScope    string    `json:"projectScope" db:"project_scope"`
	AllowedProjects *string   `json:"allowedProjects" db:"allowed_projects"`
	ExpiresAt       time.Time `json:"expiresAt" db:"expires_at"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
}

type UpdateProfileRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type CreatePATRequest struct {
	Name            string     `json:"name"`
	AccessLevel     string     `json:"accessLevel"`
	ProjectScope    string     `json:"projectScope"`
	AllowedProjects []string   `json:"allowedProjects"`
	ExpiresAt       *time.Time `json:"expiresAt"`
}

type CreatePATResponse struct {
	Token *PersonalAccessToken `json:"token"`
	Plain string               `json:"plain"`
}
