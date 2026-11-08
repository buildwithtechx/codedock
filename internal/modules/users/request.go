package users

import "time"

type UpdateProfileRequest struct {
	Name string   `json:"name"`
	Role UserRole `json:"role"`
}

type RequestEmailChangeRequest struct {
	NewEmail string `json:"newEmail" validate:"required"`
}

type VerifyEmailChangeRequest struct {
	OTP string `json:"otp" validate:"required"`
}

type CreatePATRequest struct {
	Name            string     `json:"name"`
	AccessLevel     string     `json:"accessLevel"`
	ProjectScope    string     `json:"projectScope"`
	AllowedProjects []string   `json:"allowedProjects"`
	ExpiresAt       *time.Time `json:"expiresAt"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
}

type SetAppearanceRequest struct {
	Appearance *AppearancePreferences `json:"appearance"`
}
