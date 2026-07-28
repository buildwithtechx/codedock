package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	"codedock.run/codedock/internal/utils"
)

type Mailer interface {
	SendSystemEmail(ctx context.Context, templateName string, toAddress string, subject string, data any) error
}

type AuthService struct {
	userRepo         repositories.UserRepository
	settingsRepo     repositories.SettingsRepository
	notifRepo        repositories.NotificationSettingsRepository
	projectSettings  repositories.ProjectSettingsRepository
	tokenService     *TokenService
	mailer           Mailer
	refreshTokenRepo repositories.RefreshTokenRepository
}

func NewAuthService(
	ur repositories.UserRepository,
	sr repositories.SettingsRepository,
	nr repositories.NotificationSettingsRepository,
	psr repositories.ProjectSettingsRepository,
	ts *TokenService,
	mailer Mailer,
	rtr repositories.RefreshTokenRepository,
) *AuthService {
	return &AuthService{
		userRepo:         ur,
		settingsRepo:     sr,
		notifRepo:        nr,
		projectSettings:  psr,
		tokenService:     ts,
		mailer:           mailer,
		refreshTokenRepo: rtr,
	}
}

func (a *AuthService) getBaseUrl(ctx context.Context, fallback string) string {
	if s, err := a.settingsRepo.GetServerSettings(ctx); err == nil && s != nil && s.PanelDomain != "" {
		if strings.HasPrefix(s.PanelDomain, "http") {
			return s.PanelDomain
		}
		return "https://" + s.PanelDomain
	}
	return fallback
}

func (a *AuthService) Register(ctx context.Context, name, email, password, originUrl string) (*models.User, string, string, error) {
	if email == "" || password == "" || name == "" {
		return nil, "", "", errors.New("name, email and password are required")
	}
	_, total, _ := a.userRepo.ListUsers(ctx, 1, 0)
	isInitial := total == 0
	cfg, _ := a.settingsRepo.GetServerSettings(ctx)
	if cfg != nil && !cfg.RegistrationEnabled && !isInitial {
		return nil, "", "", errors.New("user registration is disabled on this server")
	}
	if cfg != nil && !isInitial && strings.TrimSpace(cfg.RegistrationDomainAllowlist) != "" {
		allowed := false
		for d := range strings.SplitSeq(cfg.RegistrationDomainAllowlist, ",") {
			d = strings.TrimSpace(d)
			if d != "" && strings.HasSuffix(strings.ToLower(email), "@"+strings.ToLower(d)) {
				allowed = true
				break
			}
		}
		if !allowed {
			return nil, "", "", errors.New("email domain is not allowed on this server")
		}
	}
	if err := utils.ValidatePassword(password); err != nil {
		return nil, "", "", err
	}
	existing, _ := a.userRepo.GetUserByEmail(ctx, email)
	if existing != nil {
		return nil, "", "", errors.New("user already exists with that email")
	}
	hashed, err := utils.HashPassword(password)
	if err != nil {
		return nil, "", "", err
	}
	role := models.UserRoleMember
	if total == 0 {
		role = models.UserRoleOwner
	}
	u := &models.User{
		ID:            uuid.New().String(),
		Email:         email,
		Name:          name,
		PasswordHash:  string(hashed),
		Role:          role,
		IsActive:      true,
		EmailVerified: isInitial,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := a.userRepo.CreateUser(ctx, u); err != nil {
		return nil, "", "", err
	}

	if !u.EmailVerified {
		_ = a.SendEmailVerification(ctx, u.Email, originUrl)
	}

	token, err := a.tokenService.GenerateToken(u)
	if err != nil {
		return nil, "", "", err
	}
	refreshToken, err := a.tokenService.GenerateRefreshToken(u)
	if err != nil {
		return nil, "", "", err
	}
	uCopy := *u
	uCopy.PasswordHash = ""
	return &uCopy, token, refreshToken, nil
}

func (a *AuthService) Login(ctx context.Context, email, password, totpCode string) (*models.User, string, string, error) {
	if email == "" || password == "" {
		return nil, "", "", errors.New("email and password are required")
	}
	u, err := a.userRepo.GetUserByEmail(ctx, email)
	if err != nil || u == nil {
		return nil, "", "", errors.New("invalid email or password")
	}
	if !u.IsActive {
		return nil, "", "", errors.New("account is disabled")
	}
	if !utils.CheckPasswordHash(password, u.PasswordHash) {
		return nil, "", "", errors.New("invalid email or password")
	}

	if u.TOTPEnabled {
		if totpCode == "" {
			return nil, "", "", errors.New("2FA code required")
		}
		secret, recoveryCodes, err := a.userRepo.GetUserTOTPSecret(ctx, u.ID)
		if err != nil || secret == "" {
			return nil, "", "", errors.New("failed to retrieve TOTP configuration")
		}
		if !ValidateTOTP(secret, totpCode) {
			validRecovery := false
			for _, rc := range recoveryCodes {
				if rc == totpCode {
					validRecovery = true
					break
				}
			}
			if !validRecovery {
				return nil, "", "", errors.New("invalid 2FA verification code")
			}
		}
	}

	token, err := a.tokenService.GenerateToken(u)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := a.tokenService.GenerateRefreshToken(u)
	if err != nil {
		return nil, "", "", err
	}

	if a.refreshTokenRepo != nil {
		tokenHash := repositories.HashToken(refreshToken)
		_ = a.refreshTokenRepo.StoreToken(ctx, u.ID, tokenHash, time.Now().Add(7*24*time.Hour))
	}

	now := time.Now()
	u.LastLogin = &now
	_ = a.userRepo.UpdateUser(ctx, u)

	uCopy := *u
	uCopy.PasswordHash = ""
	return &uCopy, token, refreshToken, nil
}

func (a *AuthService) RefreshToken(ctx context.Context, refreshTokenStr string) (*models.User, string, string, error) {
	if refreshTokenStr == "" {
		return nil, "", "", errors.New("refresh token is required")
	}

	userID, tokenHash, err := a.tokenService.ValidateRefreshToken(refreshTokenStr)
	if err != nil {
		return nil, "", "", errors.New("invalid or expired refresh token")
	}

	if a.refreshTokenRepo != nil {
		inboundHash := repositories.HashToken(refreshTokenStr)
		revoked, err := a.refreshTokenRepo.IsRevoked(ctx, inboundHash)
		if err != nil || revoked {
			return nil, "", "", errors.New("refresh token has been revoked")
		}
		_ = a.refreshTokenRepo.RevokeToken(ctx, inboundHash)
	}

	u, err := a.userRepo.GetUserByID(ctx, userID)
	if err != nil || u == nil {
		return nil, "", "", errors.New("user not found")
	}
	if !u.IsActive {
		return nil, "", "", errors.New("account is disabled")
	}
	expectedHash := u.PasswordHash
	if len(expectedHash) > 10 {
		expectedHash = expectedHash[:10]
	}
	if expectedHash != tokenHash {
		return nil, "", "", errors.New("session revoked: password changed")
	}

	token, err := a.tokenService.GenerateToken(u)
	if err != nil {
		return nil, "", "", err
	}

	newRefreshToken, err := a.tokenService.GenerateRefreshToken(u)
	if err != nil {
		return nil, "", "", err
	}

	if a.refreshTokenRepo != nil {
		newHash := repositories.HashToken(newRefreshToken)
		_ = a.refreshTokenRepo.StoreToken(ctx, u.ID, newHash, time.Now().Add(7*24*time.Hour))
	}

	uCopy := *u
	uCopy.PasswordHash = ""
	return &uCopy, token, newRefreshToken, nil
}
