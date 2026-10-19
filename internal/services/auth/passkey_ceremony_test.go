package auth

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/descope/virtualwebauthn"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/google/uuid"

	"codedock/internal/models"
)

func TestPasskeyVirtualCeremonyAnd2FAGating(t *testing.T) {
	svc, repo, userRepo := setupTestPasskeyService(t)
	ctx := context.Background()

	rp := virtualwebauthn.RelyingParty{
		ID:     "localhost",
		Origin: "http://localhost:3000",
	}

	user := &models.User{
		ID:          uuid.NewString(),
		Email:       "virtual-passkey@example.com",
		Name:        "Virtual WebAuthn User",
		IsActive:    true,
		TOTPEnabled: true,
	}
	_ = userRepo.CreateUser(ctx, user)

	// 1. Begin registration
	creationOptions, err := svc.BeginRegistration(ctx, user.ID)
	if err != nil {
		t.Fatalf("BeginRegistration failed: %v", err)
	}

	optionsJSON, err := json.Marshal(creationOptions.Response)
	if err != nil {
		t.Fatalf("marshal creation options failed: %v", err)
	}

	attestationOptions, err := virtualwebauthn.ParseAttestationOptions(string(optionsJSON))
	if err != nil {
		t.Fatalf("parse attestation options: %v", err)
	}

	authenticator := virtualwebauthn.NewAuthenticator()
	credential := virtualwebauthn.NewCredential(virtualwebauthn.KeyTypeEC2)

	attestationJSON := virtualwebauthn.CreateAttestationResponse(rp, authenticator, credential, *attestationOptions)
	parsedCreation, err := protocol.ParseCredentialCreationResponseBody(strings.NewReader(attestationJSON))
	if err != nil {
		t.Fatalf("parse credential creation body: %v", err)
	}

	// 2. Finish registration
	regResult, err := svc.FinishRegistration(ctx, user.ID, "Virtual FIDO2 Key", parsedCreation, "127.0.0.1")
	if err != nil {
		t.Fatalf("FinishRegistration failed: %v", err)
	}
	if regResult.Name != "Virtual FIDO2 Key" {
		t.Errorf("expected key name 'Virtual FIDO2 Key', got %s", regResult.Name)
	}

	// 3. Begin authentication
	assertionOptionsWrapped, err := svc.BeginAuthentication(ctx)
	if err != nil {
		t.Fatalf("BeginAuthentication failed: %v", err)
	}

	assertionJSON, err := json.Marshal(assertionOptionsWrapped.Response)
	if err != nil {
		t.Fatalf("marshal assertion options failed: %v", err)
	}

	parsedAssertionOpts, err := virtualwebauthn.ParseAssertionOptions(string(assertionJSON))
	if err != nil {
		t.Fatalf("parse assertion options: %v", err)
	}

	assertionPayload := virtualwebauthn.CreateAssertionResponse(rp, authenticator, credential, *parsedAssertionOpts)
	parsedAssertion, err := protocol.ParseCredentialRequestResponseBody(strings.NewReader(assertionPayload))
	if err != nil {
		t.Fatalf("parse assertion payload: %v", err)
	}

	// 4. Test 2FA gate (First attempt without code)
	_, _, _, err = svc.FinishAuthentication(ctx, parsedAssertion, "", "127.0.0.1")
	if err == nil || !strings.Contains(err.Error(), "2FA code required") {
		t.Fatalf("expected '2FA code required', got %v", err)
	}

	// Challenge must NOT have been consumed
	if len(repo.challenges) == 0 {
		t.Fatal("challenge was consumed prematurely on missing 2FA code")
	}

	// 5. Test 2FA gate with invalid TOTP code
	_, _, _, err = svc.FinishAuthentication(ctx, parsedAssertion, "999999", "127.0.0.1")
	if err == nil || !strings.Contains(err.Error(), "invalid 2FA verification code") {
		t.Fatalf("expected 'invalid 2FA verification code', got %v", err)
	}

	// Challenge must still NOT have been consumed
	if len(repo.challenges) == 0 {
		t.Fatal("challenge was consumed prematurely on mistyped 2FA code")
	}

	// 6. Test 2FA gate with valid recovery code (from mock user repo: "1111-2222")
	authedUser, token, refresh, err := svc.FinishAuthentication(ctx, parsedAssertion, "1111-2222", "127.0.0.1")
	if err != nil {
		t.Fatalf("FinishAuthentication with 2FA code failed: %v", err)
	}
	if authedUser == nil || token == "" || refresh == "" {
		t.Fatal("expected non-empty user, token, and refresh token")
	}

	// 7. Verify replay fails (challenge is now atomically consumed)
	_, _, _, err = svc.FinishAuthentication(ctx, parsedAssertion, "1111-2222", "127.0.0.1")
	if err == nil || !strings.Contains(err.Error(), "challenge expired or already used") {
		t.Fatalf("expected replay rejection 'challenge expired or already used', got %v", err)
	}
}

func TestPasskeyVirtualCeremonyDirectLogin(t *testing.T) {
	svc, _, userRepo := setupTestPasskeyService(t)
	ctx := context.Background()

	rp := virtualwebauthn.RelyingParty{
		ID:     "localhost",
		Origin: "http://localhost:3000",
	}

	user := &models.User{
		ID:          uuid.NewString(),
		Email:       "direct-passkey@example.com",
		Name:        "Direct Passkey User",
		IsActive:    true,
		TOTPEnabled: false,
	}
	_ = userRepo.CreateUser(ctx, user)

	// Register
	creationOptions, err := svc.BeginRegistration(ctx, user.ID)
	if err != nil {
		t.Fatalf("BeginRegistration failed: %v", err)
	}

	optionsJSON, _ := json.Marshal(creationOptions.Response)
	attestationOptions, _ := virtualwebauthn.ParseAttestationOptions(string(optionsJSON))
	authenticator := virtualwebauthn.NewAuthenticator()
	credential := virtualwebauthn.NewCredential(virtualwebauthn.KeyTypeEC2)

	attestationJSON := virtualwebauthn.CreateAttestationResponse(rp, authenticator, credential, *attestationOptions)
	parsedCreation, _ := protocol.ParseCredentialCreationResponseBody(strings.NewReader(attestationJSON))
	_, err = svc.FinishRegistration(ctx, user.ID, "MacBook TouchID", parsedCreation, "127.0.0.1")
	if err != nil {
		t.Fatalf("FinishRegistration failed: %v", err)
	}

	// Authenticate directly (no 2FA required)
	assertionOptionsWrapped, _ := svc.BeginAuthentication(ctx)
	assertionJSON, _ := json.Marshal(assertionOptionsWrapped.Response)
	parsedAssertionOpts, _ := virtualwebauthn.ParseAssertionOptions(string(assertionJSON))
	assertionPayload := virtualwebauthn.CreateAssertionResponse(rp, authenticator, credential, *parsedAssertionOpts)
	parsedAssertion, _ := protocol.ParseCredentialRequestResponseBody(strings.NewReader(assertionPayload))

	authedUser, token, refresh, err := svc.FinishAuthentication(ctx, parsedAssertion, "", "127.0.0.1")
	if err != nil {
		t.Fatalf("FinishAuthentication failed: %v", err)
	}
	if authedUser == nil || token == "" || refresh == "" {
		t.Fatal("expected non-empty user, token, and refresh token")
	}
	if authedUser.Email != user.Email {
		t.Errorf("expected user email %s, got %s", user.Email, authedUser.Email)
	}
}
