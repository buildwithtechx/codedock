package auth

import (
	"database/sql"
	"fmt"

	"codedock/internal/modules/audit"
	"codedock/internal/modules/settings"
	"codedock/internal/modules/users"
)

type Module struct {
	TokenService      *TokenService
	AuthService       *AuthService
	AuthHandler       *AuthHandler
	OAuthHandler      *OAuthHandler
	DeviceFlowHandler *DeviceFlowHandler
	PasskeyHandler    *PasskeyHandler
}

func NewModule(
	db *sql.DB,
	userRepo users.UserRepository,
	userService *users.UserService,
	settingsRepo settings.SettingsRepository,
	auditService *audit.AuditService,
	mailer Mailer,
) (*Module, error) {
	oauthRepo := NewOAuthRepo(db)
	refreshTokenRepo := NewRefreshTokenRepo(db)
	passkeyRepo := NewPasskeyRepo(db)

	tokenService, err := NewTokenService()
	if err != nil {
		return nil, fmt.Errorf("token service: %w", err)
	}
	authService := NewAuthService(userRepo, settingsRepo, tokenService, mailer, refreshTokenRepo)
	oAuthService := NewOAuthService(oauthRepo, userRepo, tokenService)
	deviceFlowService := NewDeviceFlowService()
	passkeyService, err := NewPasskeyService(passkeyRepo, userRepo, tokenService, refreshTokenRepo, auditService)
	if err != nil {
		return nil, fmt.Errorf("initialize passkey service: %w", err)
	}

	return &Module{
		TokenService:      tokenService,
		AuthService:       authService,
		AuthHandler:       NewAuthHandler(authService, auditService),
		OAuthHandler:      NewOAuthHandler(oAuthService),
		DeviceFlowHandler: NewDeviceFlowHandler(deviceFlowService, userService, auditService),
		PasskeyHandler:    NewPasskeyHandler(passkeyService),
	}, nil
}
