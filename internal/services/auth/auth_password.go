package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"codedock.run/codedock/internal/config"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
)

func isEmailEnabled(cfg *models.NotificationSettings) bool {
	if cfg != nil && (cfg.SMTPEnabled || cfg.ResendEnabled) {
		return true
	}
	appCfg := config.Get()
	if appCfg.Cloud.Enabled {
		return true
	}
	return appCfg.SMTP.Host != "" || appCfg.Resend.APIKey != ""
}

func (a *AuthService) ForgotPassword(ctx context.Context, email string, originUrl string) error {
	originUrl = a.getBaseUrl(ctx, originUrl)
	if email == "" {
		return errors.New("email is required")
	}

	cfg, _ := a.notifRepo.GetNotificationSettings(ctx)
	if !isEmailEnabled(cfg) {
		return errors.New("your team is yet to set or enable email")
	}

	u, err := a.userRepo.GetUserByEmail(ctx, email)
	if err != nil || u == nil {
		return nil
	}

	token, err := a.tokenService.GeneratePasswordResetToken(u)
	if err != nil {
		return err
	}

	data := map[string]any{
		"ResetUrl": originUrl + "/reset-password?token=" + token,
	}

	err = a.mailer.SendSystemEmail(ctx, "password_reset", u.Email, "Reset Your Password", data)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) ResetPassword(ctx context.Context, tokenStr, newPassword string) error {
	if err := utils.ValidatePassword(newPassword); err != nil {
		return err
	}

	email, tokenHash, err := a.tokenService.ValidatePasswordResetToken(tokenStr)
	if err != nil {
		return errors.New("invalid or expired reset token")
	}

	u, err := a.userRepo.GetUserByEmail(ctx, email)
	if err != nil || u == nil {
		return utils.NewNotFoundError("User", email)
	}

	expectedHash := u.PasswordHash
	if len(expectedHash) > 10 {
		expectedHash = expectedHash[:10]
	}
	if expectedHash != tokenHash {
		return errors.New("reset token revoked: password has already been changed")
	}

	hashed, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hashed)
	if err := a.userRepo.UpdateUser(ctx, u); err != nil {
		return err
	}

	return nil
}

func (a *AuthService) InviteUser(ctx context.Context, email string, role models.UserRole, originUrl string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	existing, _ := a.userRepo.GetUserByEmail(ctx, email)
	if existing != nil {
		return existing, nil
	}
	u := &models.User{
		ID:           uuid.New().String(),
		Email:        email,
		Name:         strings.Split(email, "@")[0],
		PasswordHash: "INVITED_NO_LOGIN_ALLOWED_MUST_RESET",
		Role:         role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := a.userRepo.CreateUser(ctx, u); err != nil {
		return nil, err
	}

	_ = a.ForgotPassword(ctx, u.Email, originUrl)

	return u, nil
}

func (a *AuthService) SendEmailVerification(ctx context.Context, email string, originUrl string) error {
	originUrl = a.getBaseUrl(ctx, originUrl)
	if email == "" {
		return errors.New("email is required")
	}

	cfg, _ := a.notifRepo.GetNotificationSettings(ctx)
	if !isEmailEnabled(cfg) {
		return errors.New("email is not configured on this server")
	}

	u, err := a.userRepo.GetUserByEmail(ctx, email)
	if err != nil || u == nil {
		return nil
	}
	if u.EmailVerified {
		return errors.New("email is already verified")
	}

	token, err := a.tokenService.GenerateEmailVerificationToken(u.Email)
	if err != nil {
		return err
	}

	data := map[string]any{
		"VerifyUrl": originUrl + "/verify-email?token=" + token,
	}

	err = a.mailer.SendSystemEmail(ctx, "email_verification", u.Email, "Verify Your Email", data)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) VerifyEmail(ctx context.Context, tokenStr string) error {
	email, err := a.tokenService.ValidateEmailVerificationToken(tokenStr)
	if err != nil {
		return errors.New("invalid or expired verification token")
	}

	u, err := a.userRepo.GetUserByEmail(ctx, email)
	if err != nil || u == nil {
		return utils.NewNotFoundError("User", email)
	}

	if u.EmailVerified {
		return nil
	}

	u.EmailVerified = true
	if err := a.userRepo.UpdateUser(ctx, u); err != nil {
		return err
	}

	return nil
}
