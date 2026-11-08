package auth

import (
	"encoding/json"
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
