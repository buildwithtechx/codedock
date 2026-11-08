package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"codedock/internal/config"
	"codedock/internal/modules/audit"
	"codedock/internal/modules/users"
)

type memoryPasskeyRepo struct {
	creds      map[string]*PasskeyCredential
	challenges map[string]*PasskeyChallenge
}

func newMemoryPasskeyRepo() *memoryPasskeyRepo {
	return &memoryPasskeyRepo{
		creds:      make(map[string]*PasskeyCredential),
		challenges: make(map[string]*PasskeyChallenge),
	}
}

func (m *memoryPasskeyRepo) CreateCredential(_ context.Context, cred *PasskeyCredential) error {
	m.creds[cred.ID] = cred
	return nil
}

func (m *memoryPasskeyRepo) GetCredentialByID(_ context.Context, id string) (*PasskeyCredential, error) {
	if c, ok := m.creds[id]; ok {
		return c, nil
	}
	return nil, errors.New("not found")
}

func (m *memoryPasskeyRepo) GetCredentialByCredentialID(_ context.Context, id string) (*PasskeyCredential, error) {
	for _, c := range m.creds {
		if c.CredentialID == id {
			return c, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *memoryPasskeyRepo) ListCredentialsByUserID(_ context.Context, userID string) ([]*PasskeyCredential, error) {
	var list []*PasskeyCredential
	for _, c := range m.creds {
		if c.UserID == userID {
			list = append(list, c)
		}
	}
	return list, nil
}

func (m *memoryPasskeyRepo) UpdateCredentialCounter(_ context.Context, id string, counter int64, lastUsedAt time.Time) error {
	if c, ok := m.creds[id]; ok {
		c.Counter = counter
		c.LastUsedAt = &lastUsedAt
		return nil
	}
	return errors.New("not found")
}

func (m *memoryPasskeyRepo) RenameCredential(_ context.Context, id, userID, name string) error {
	if c, ok := m.creds[id]; ok && c.UserID == userID {
		c.Name = name
		return nil
	}
	return errors.New("not found or unauthorized")
}

func (m *memoryPasskeyRepo) DeleteCredential(_ context.Context, id, userID string) error {
	if c, ok := m.creds[id]; ok && c.UserID == userID {
		delete(m.creds, id)
		return nil
	}
	return errors.New("not found or unauthorized")
}

func (m *memoryPasskeyRepo) SaveChallenge(_ context.Context, ch *PasskeyChallenge) error {
	m.challenges[ch.Challenge] = ch
	return nil
}

func (m *memoryPasskeyRepo) GetAndConsumeChallenge(_ context.Context, challenge, ceremonyType string) (*PasskeyChallenge, error) {
	ch, ok := m.challenges[challenge]
	if !ok || ch.CeremonyType != ceremonyType {
		return nil, errors.New("challenge expired or not found")
	}
	if time.Now().After(ch.ExpiresAt) {
		delete(m.challenges, challenge)
		return nil, errors.New("challenge expired")
	}
	delete(m.challenges, challenge)
	return ch, nil
}

type memoryUserRepo struct {
	users.UserRepository
	users map[string]*users.User
}

func (m *memoryUserRepo) CreateUser(_ context.Context, u *users.User) error {
	m.users[u.ID] = u
	return nil
}

func (m *memoryUserRepo) GetUserByID(_ context.Context, id string) (*users.User, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return nil, errors.New("not found")
}

func (m *memoryUserRepo) UpdateUser(_ context.Context, u *users.User) error {
	m.users[u.ID] = u
	return nil
}

func (m *memoryUserRepo) GetUserTOTPSecret(_ context.Context, _ string) (string, []string, error) {
	return "JBSWY3DPEHPK3PXP", []string{"1111-2222"}, nil
}

func setupTestPasskeyService(t *testing.T) (*PasskeyService, *memoryPasskeyRepo, *memoryUserRepo) {
	t.Helper()
	cfg := config.Get()
	cfg.Server.DashboardURL = "http://localhost:3000"
	cfg.Security.JWTSecret = "test-jwt-secret-for-passkeys-unit-testing"

	passkeyRepo := newMemoryPasskeyRepo()
	userRepo := &memoryUserRepo{users: make(map[string]*users.User)}
	tokenService, err := NewTokenService()
	if err != nil {
		t.Fatalf("token service: %v", err)
	}
	svc, err := NewPasskeyService(passkeyRepo, userRepo, tokenService, nil, audit.NewAuditService(nil, nil))
	if err != nil {
		t.Fatalf("create passkey service: %v", err)
	}
	return svc, passkeyRepo, userRepo
}

func TestPasskeyRegistrationCeremony(t *testing.T) {
	svc, repo, userRepo := setupTestPasskeyService(t)
	ctx := context.Background()

	user := &users.User{
		ID:       uuid.NewString(),
		Email:    "passkey-test@example.com",
		Name:     "Passkey User",
		IsActive: true,
	}
	_ = userRepo.CreateUser(ctx, user)

	options, err := svc.BeginRegistration(ctx, user.ID)
	if err != nil {
		t.Fatalf("BeginRegistration failed: %v", err)
	}
	if options == nil || len(options.Response.Challenge) == 0 {
		t.Fatal("expected non-empty challenge in creation options")
	}
	if options.Response.RelyingParty.Name != "Codedock" {
		t.Errorf("expected RP name Codedock, got %s", options.Response.RelyingParty.Name)
	}

	if len(repo.challenges) != 1 {
		t.Fatalf("expected 1 challenge in repo, got %d", len(repo.challenges))
	}

	consumed, err := repo.GetAndConsumeChallenge(ctx, options.Response.Challenge.String(), "registration")
	if err != nil {
		t.Fatalf("failed to consume challenge: %v", err)
	}
	if consumed.UserID != user.ID {
		t.Errorf("expected user ID %s, got %s", user.ID, consumed.UserID)
	}

	_, err = repo.GetAndConsumeChallenge(ctx, options.Response.Challenge.String(), "registration")
	if err == nil {
		t.Fatal("expected replayed challenge to fail")
	}
}

func TestPasskeyAuthenticationCeremonyAnd2FA(t *testing.T) {
	svc, repo, userRepo := setupTestPasskeyService(t)
	ctx := context.Background()

	user := &users.User{
		ID:          uuid.NewString(),
		Email:       "totp-user@example.com",
		Name:        "TOTP User",
		IsActive:    true,
		TOTPEnabled: true,
	}
	_ = userRepo.CreateUser(ctx, user)

	cred := &PasskeyCredential{
		ID:           uuid.NewString(),
		UserID:       user.ID,
		Name:         "YubiKey",
		CredentialID: "test-credential-id",
		PublicKey:    "test-public-key",
		CreatedAt:    time.Now().UTC(),
	}
	_ = repo.CreateCredential(ctx, cred)

	options, err := svc.BeginAuthentication(ctx)
	if err != nil {
		t.Fatalf("BeginAuthentication failed: %v", err)
	}
	if options == nil || len(options.Response.Challenge) == 0 {
		t.Fatal("expected non-empty assertion challenge")
	}

	list, err := svc.ListUserPasskeys(ctx, user.ID)
	if err != nil {
		t.Fatalf("ListUserPasskeys failed: %v", err)
	}
	if len(list) != 1 || list[0].Name != "YubiKey" {
		t.Fatalf("unexpected list response: %+v", list)
	}

	err = svc.RenamePasskey(ctx, cred.ID, user.ID, "Office YubiKey", "127.0.0.1")
	if err != nil {
		t.Fatalf("RenamePasskey failed: %v", err)
	}
	if repo.creds[cred.ID].Name != "Office YubiKey" {
		t.Errorf("expected renamed key, got %s", repo.creds[cred.ID].Name)
	}

	err = svc.DeletePasskey(ctx, cred.ID, user.ID, "127.0.0.1")
	if err != nil {
		t.Fatalf("DeletePasskey failed: %v", err)
	}
	if _, ok := repo.creds[cred.ID]; ok {
		t.Error("expected key to be deleted")
	}
}
