package auth

import (
	"database/sql"
	"fmt"
)

type Module struct {
	TokenService      *TokenService
	SettingsService   *SettingsService
	AuthService       *AuthService
	UserService       *UserService
	OrgService        *OrganizationService
	AuditService      *AuditService
	SettingsHandler   *SettingsHandler
	AppearanceHandler *AppearanceHandler
	UserHandler       *UserHandler
	AuthHandler       *AuthHandler
	OAuthHandler      *OAuthHandler
	DeviceFlowHandler *DeviceFlowHandler
	PasskeyHandler    *PasskeyHandler
	OrgHandler        *OrganizationHandler
	OnboardingHandler *OnboardingHandler
	AuditLogHandler   *AuditLogHandler
}

func NewModule(db *sql.DB, v Vault, mailer Mailer) (*Module, error) {
	userRepo := NewUserRepo(db)
	settingsRepo := NewSettingsRepo(db, v)
	oauthRepo := NewOAuthRepo(db)
	auditRepo := NewAuditLogRepo(db)
	orgRepo := NewOrganizationRepository(db)
	refreshTokenRepo := NewRefreshTokenRepo(db)
	passkeyRepo := NewPasskeyRepo(db)

	tokenService, err := NewTokenService()
	if err != nil {
		return nil, fmt.Errorf("token service: %w", err)
	}
	settingsService := NewSettingsService(settingsRepo)
	authService := NewAuthService(userRepo, settingsRepo, tokenService, mailer, refreshTokenRepo)
	userService := NewUserService(userRepo)
	orgService := NewOrganizationService(orgRepo, userRepo)
	oAuthService := NewOAuthService(oauthRepo, userRepo, tokenService)
	auditService := NewAuditService(auditRepo, NewSQLAuditResourceSearcher(db))
	deviceFlowService := NewDeviceFlowService()
	passkeyService, err := NewPasskeyService(passkeyRepo, userRepo, tokenService, refreshTokenRepo, auditService)
	if err != nil {
		return nil, fmt.Errorf("initialize passkey service: %w", err)
	}

	return &Module{
		TokenService:      tokenService,
		SettingsService:   settingsService,
		AuthService:       authService,
		UserService:       userService,
		OrgService:        orgService,
		AuditService:      auditService,
		SettingsHandler:   NewSettingsHandler(settingsService),
		AppearanceHandler: NewAppearanceHandler(userRepo),
		UserHandler:       NewUserHandler(userService, mailer, auditService),
		AuthHandler:       NewAuthHandler(authService, auditService),
		OAuthHandler:      NewOAuthHandler(oAuthService),
		DeviceFlowHandler: NewDeviceFlowHandler(deviceFlowService, userService, auditService),
		PasskeyHandler:    NewPasskeyHandler(passkeyService),
		OrgHandler:        NewOrganizationHandler(orgService, auditService),
		OnboardingHandler: NewOnboardingHandler(userService),
		AuditLogHandler:   NewAuditLogHandler(auditService),
	}, nil
}
