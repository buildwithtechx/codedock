package auth

import (
	"encoding/json"
	"time"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	TOTPCode string `json:"totpCode,omitempty"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type AdminInviteUserRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type ResendVerificationEmailRequest struct {
	Email string `json:"email"`
}

type VerifyEmailRequest struct {
	Token string `json:"token"`
}

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

type Verify2FARequest struct {
	Passcode string `json:"passcode"`
}

type PasskeyRenameRequest struct {
	Name string `json:"name"`
}

type verifyRegistrationPayload struct {
	Name     string          `json:"name"`
	Response json.RawMessage `json:"response"`
}

type verifyAuthenticationPayload struct {
	TOTPCode string          `json:"totpCode"`
	Response json.RawMessage `json:"response"`
}

type CreateDeviceCodeRequest struct {
	ClientName string `json:"clientName"`
}

type PollDeviceTokenRequest struct {
	DeviceCode string `json:"deviceCode"`
}

type AuthorizeDeviceRequest struct {
	UserCode string `json:"userCode"`
	Action   string `json:"action"`
}

type CreateOrganizationRequest struct {
	Name string `json:"name"`
}

type UpdateOrganizationRequest struct {
	Name string `json:"name"`
}

type InviteOrganizationMemberRequest struct {
	Email      string           `json:"email"`
	Permission MemberPermission `json:"permission"`
	Delivery   string           `json:"delivery,omitempty"`
	Projects   []string         `json:"projects,omitempty"`
}
