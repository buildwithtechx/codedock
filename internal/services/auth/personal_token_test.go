package auth

import (
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"
)

type tokenCreationStore struct {
	repositories.UserRepository
	created *models.PersonalAccessToken
}

func (s *tokenCreationStore) CreatePAT(_ context.Context, pat *models.PersonalAccessToken) error {
	s.created = pat
	return nil
}

func TestPersonalTokenCreationMatchesUIContract(t *testing.T) {
	store := &tokenCreationStore{}
	service := NewUserService(store)
	token, plain, err := service.CreatePAT(context.Background(), "user", "automation", "read", "specific", []string{"project"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if token.ExpiresAt != nil {
		t.Fatal("no expiration token acquired an expiry")
	}
	if len(plain) != 68 || plain[:4] != "vpt_" {
		t.Fatal("invalid bearer token")
	}
	digest := sha256.Sum256([]byte(plain))
	if store.created.TokenHash != hex.EncodeToString(digest[:]) {
		t.Fatal("credential was not hashed")
	}
	encoded, err := json.Marshal(models.CreatePATResponse{Token: token, Plain: plain})
	if err != nil {
		t.Fatal(err)
	}
	var result struct {
		Plain string `json:"plain"`
		Token struct {
			AllowedProjects []string `json:"allowedProjects"`
		} `json:"token"`
	}
	if err := json.Unmarshal(encoded, &result); err != nil {
		t.Fatal(err)
	}
	if result.Plain != plain || len(result.Token.AllowedProjects) != 1 || result.Token.AllowedProjects[0] != "project" {
		t.Fatal("token response does not match UI")
	}
}
func TestPersonalTokenCreationRejectsInvalidPolicy(t *testing.T) {
	expired := time.Now().Add(-time.Minute)
	for _, test := range []struct {
		access, scope string
		projects      []string
		expires       *time.Time
	}{
		{access: "admin", scope: "all"}, {access: "read", scope: "unknown"}, {access: "read", scope: "specific"}, {access: "read", scope: "specific", projects: []string{""}}, {access: "read", scope: "all", expires: &expired},
	} {
		store := &tokenCreationStore{}
		service := NewUserService(store)
		if _, _, err := service.CreatePAT(context.Background(), "user", "token", test.access, test.scope, test.projects, test.expires); err == nil || store.created != nil {
			t.Fatal("invalid policy was persisted")
		}
	}
}
