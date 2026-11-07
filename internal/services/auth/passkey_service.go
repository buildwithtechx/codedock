package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"

	"codedock/internal/config"
	"codedock/internal/models"
	"codedock/internal/repositories"
)

type PasskeyService struct {
	webauthn         *webauthn.WebAuthn
	passkeyRepo      repositories.PasskeyRepository
	userRepo         repositories.UserRepository
	tokenService     *TokenService
	refreshTokenRepo repositories.RefreshTokenRepository
	auditService     *AuditService
}

func NewPasskeyService(
	passkeyRepo repositories.PasskeyRepository,
	userRepo repositories.UserRepository,
	tokenService *TokenService,
	refreshTokenRepo repositories.RefreshTokenRepository,
	auditService *AuditService,
) (*PasskeyService, error) {
	cfg := config.Get()
	dashboardURL := cfg.Server.DashboardURL
	if dashboardURL == "" {
		dashboardURL = "http://localhost:3000"
	}
	u, err := url.Parse(dashboardURL)
	if err != nil {
		return nil, fmt.Errorf("invalid dashboard url: %w", err)
	}
	rpID := u.Hostname()
	origin := fmt.Sprintf("%s://%s", u.Scheme, u.Host)

	origins := []string{origin}
	for _, o := range cfg.Server.ExtraCORSOrigins {
		trimmed := strings.TrimSpace(o)
		if trimmed != "" {
			origins = append(origins, trimmed)
		}
	}

	w, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "Codedock",
		RPID:          rpID,
		RPOrigins:     origins,
	})
	if err != nil {
		return nil, fmt.Errorf("initialize webauthn: %w", err)
	}

	return &PasskeyService{
		webauthn:         w,
		passkeyRepo:      passkeyRepo,
		userRepo:         userRepo,
		tokenService:     tokenService,
		refreshTokenRepo: refreshTokenRepo,
		auditService:     auditService,
	}, nil
}

func (s *PasskeyService) BeginAuthentication(ctx context.Context) (*protocol.CredentialAssertion, error) {
	options, sessionData, err := s.webauthn.BeginDiscoverableLogin(
		webauthn.WithUserVerification(protocol.VerificationPreferred),
	)
	if err != nil {
		return nil, fmt.Errorf("begin authentication: %w", err)
	}

	sessionBytes, err := json.Marshal(sessionData)
	if err != nil {
		return nil, err
	}
	challengeStr := sessionData.Challenge
	err = s.passkeyRepo.SaveChallenge(ctx, &models.PasskeyChallenge{
		Challenge:    challengeStr,
		UserID:       "",
		CeremonyType: "authentication",
		SessionData:  string(sessionBytes),
		ExpiresAt:    time.Now().UTC().Add(5 * time.Minute),
	})
	if err != nil {
		return nil, err
	}
	return options, nil
}

func (s *PasskeyService) FinishAuthentication(ctx context.Context, parsed *protocol.ParsedCredentialAssertionData, totpCode, ip string) (*models.User, string, string, error) {
	if parsed == nil {
		return nil, "", "", errors.New("missing assertion data")
	}

	credIDStr := base64.RawURLEncoding.EncodeToString(parsed.RawID)
	storedCred, err := s.passkeyRepo.GetCredentialByCredentialID(ctx, credIDStr)
	if err != nil || storedCred == nil {
		return nil, "", "", errors.New("credential not recognized")
	}

	u, err := s.userRepo.GetUserByID(ctx, storedCred.UserID)
	if err != nil || u == nil {
		return nil, "", "", errors.New("user not found")
	}
	if !u.IsActive {
		return nil, "", "", errors.New("account is disabled")
	}

	if u.TOTPEnabled {
		if totpCode == "" {
			return nil, "", "", errors.New("2FA code required")
		}
		secret, recoveryCodes, err := s.userRepo.GetUserTOTPSecret(ctx, u.ID)
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

	challengeStr := parsed.Response.CollectedClientData.Challenge
	challengeRecord, err := s.passkeyRepo.GetAndConsumeChallenge(ctx, challengeStr, "authentication")
	if err != nil {
		return nil, "", "", errors.New("challenge expired or already used")
	}

	var sessionData webauthn.SessionData
	if err := json.Unmarshal([]byte(challengeRecord.SessionData), &sessionData); err != nil {
		return nil, "", "", fmt.Errorf("parse session data: %w", err)
	}

	webUser := toWebAuthnUser(u, []*models.PasskeyCredential{storedCred})
	if len(sessionData.UserID) == 0 {
		sessionData.UserID = webUser.WebAuthnID()
	}
	verifiedCred, err := s.webauthn.ValidateLogin(webUser, sessionData, parsed)
	if err != nil {
		return nil, "", "", fmt.Errorf("passkey verification failed: %w", err)
	}

	_ = s.passkeyRepo.UpdateCredentialCounter(ctx, storedCred.ID, int64(verifiedCred.Authenticator.SignCount), time.Now().UTC())

	token, err := s.tokenService.GenerateToken(u)
	if err != nil {
		return nil, "", "", err
	}
	refreshToken, err := s.tokenService.GenerateRefreshToken(u)
	if err != nil {
		return nil, "", "", err
	}

	if s.refreshTokenRepo != nil {
		tokenHash := repositories.HashToken(refreshToken)
		_ = s.refreshTokenRepo.StoreToken(ctx, u.ID, tokenHash, time.Now().Add(7*24*time.Hour))
	}

	now := time.Now()
	u.LastLogin = &now
	_ = s.userRepo.UpdateUser(ctx, u)

	s.auditService.LogAction(ctx, AuditActionOpts{
		UserID:    u.ID,
		Action:    "auth.login.passkey",
		Resource:  storedCred.ID,
		IPAddress: ip,
	})

	uCopy := *u
	uCopy.PasswordHash = ""
	return &uCopy, token, refreshToken, nil
}

func (s *PasskeyService) ListUserPasskeys(ctx context.Context, userID string) ([]*models.PasskeyResponse, error) {
	creds, err := s.passkeyRepo.ListCredentialsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	var resp []*models.PasskeyResponse
	for _, c := range creds {
		resp = append(resp, toPasskeyResponse(c))
	}
	return resp, nil
}

func (s *PasskeyService) RenamePasskey(ctx context.Context, id, userID, name, ip string) error {
	if err := s.passkeyRepo.RenameCredential(ctx, id, userID, name); err != nil {
		return err
	}
	s.auditService.LogAction(ctx, AuditActionOpts{
		UserID:    userID,
		Action:    "auth.passkey.rename",
		Resource:  id,
		IPAddress: ip,
	})
	return nil
}

func (s *PasskeyService) DeletePasskey(ctx context.Context, id, userID, ip string) error {
	if err := s.passkeyRepo.DeleteCredential(ctx, id, userID); err != nil {
		return err
	}
	s.auditService.LogAction(ctx, AuditActionOpts{
		UserID:    userID,
		Action:    "auth.passkey.delete",
		Resource:  id,
		IPAddress: ip,
	})
	return nil
}
