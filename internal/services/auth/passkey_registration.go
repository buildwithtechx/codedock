package auth

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"

	"codedock/internal/models"
)

func (s *PasskeyService) BeginRegistration(ctx context.Context, userID string) (*protocol.CredentialCreation, error) {
	u, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil || u == nil {
		return nil, errors.New("user not found")
	}
	creds, err := s.passkeyRepo.ListCredentialsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	webUser := toWebAuthnUser(u, creds)

	options, sessionData, err := s.webauthn.BeginRegistration(webUser,
		webauthn.WithAuthenticatorSelection(protocol.AuthenticatorSelection{
			ResidentKey:      protocol.ResidentKeyRequirementRequired,
			UserVerification: protocol.VerificationPreferred,
		}),
		webauthn.WithConveyancePreference(protocol.PreferNoAttestation),
	)
	if err != nil {
		return nil, fmt.Errorf("begin registration: %w", err)
	}

	sessionBytes, err := json.Marshal(sessionData)
	if err != nil {
		return nil, err
	}
	challengeStr := sessionData.Challenge
	err = s.passkeyRepo.SaveChallenge(ctx, &models.PasskeyChallenge{
		Challenge:    challengeStr,
		UserID:       userID,
		CeremonyType: "registration",
		SessionData:  string(sessionBytes),
		ExpiresAt:    time.Now().UTC().Add(5 * time.Minute),
	})
	if err != nil {
		return nil, err
	}
	return options, nil
}

func (s *PasskeyService) FinishRegistration(ctx context.Context, userID, name string, parsed *protocol.ParsedCredentialCreationData, ip string) (*models.PasskeyResponse, error) {
	if parsed == nil {
		return nil, errors.New("missing registration data")
	}
	challengeStr := parsed.Response.CollectedClientData.Challenge
	challengeRecord, err := s.passkeyRepo.GetAndConsumeChallenge(ctx, challengeStr, "registration")
	if err != nil {
		return nil, errors.New("challenge expired or already used")
	}
	if challengeRecord.UserID != userID {
		return nil, errors.New("unauthorized challenge")
	}

	var sessionData webauthn.SessionData
	if err := json.Unmarshal([]byte(challengeRecord.SessionData), &sessionData); err != nil {
		return nil, fmt.Errorf("parse session data: %w", err)
	}

	u, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil || u == nil {
		return nil, errors.New("user not found")
	}
	creds, err := s.passkeyRepo.ListCredentialsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	webUser := toWebAuthnUser(u, creds)

	cred, err := s.webauthn.CreateCredential(webUser, sessionData, parsed)
	if err != nil {
		return nil, fmt.Errorf("verify registration: %w", err)
	}

	deviceType := "singleDevice"
	if cred.Flags.BackupEligible || cred.Flags.BackupState {
		deviceType = "multiDevice"
	}
	var transportStrs []string
	for _, t := range cred.Transport {
		transportStrs = append(transportStrs, string(t))
	}

	passkeyName := strings.TrimSpace(name)
	if passkeyName == "" {
		passkeyName = "Passkey"
	}

	record := &models.PasskeyCredential{
		ID:           uuid.NewString(),
		UserID:       userID,
		Name:         passkeyName,
		CredentialID: base64.RawURLEncoding.EncodeToString(cred.ID),
		PublicKey:    base64.RawURLEncoding.EncodeToString(cred.PublicKey),
		Counter:      int64(cred.Authenticator.SignCount),
		DeviceType:   deviceType,
		BackedUp:     cred.Flags.BackupState,
		Transports:   strings.Join(transportStrs, ","),
		AAGUID:       hex.EncodeToString(cred.Authenticator.AAGUID),
		CreatedAt:    time.Now().UTC(),
	}

	if err := s.passkeyRepo.CreateCredential(ctx, record); err != nil {
		return nil, err
	}

	s.auditService.LogAction(ctx, AuditActionOpts{
		UserID:    userID,
		Action:    "auth.passkey.create",
		Resource:  record.ID,
		IPAddress: ip,
	})

	return toPasskeyResponse(record), nil
}
