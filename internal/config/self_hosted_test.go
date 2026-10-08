package config

import (
	"context"
	"encoding/hex"
	"errors"
	"testing"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/pkg/types"
)

type memorySelfHostedStore struct {
	stored models.SelfHostedConfig
	saves  int
}

func (s *memorySelfHostedStore) Load(context.Context) (*models.SelfHostedConfig, error) {
	stored := s.stored
	return &stored, nil
}

func (s *memorySelfHostedStore) Save(_ context.Context, stored *models.SelfHostedConfig) error {
	s.stored = *stored
	s.saves++
	return nil
}

type failingSelfHostedStore struct {
	saved bool
}

func (s *failingSelfHostedStore) Load(context.Context) (*models.SelfHostedConfig, error) {
	return nil, errors.New("storage unavailable")
}

func (s *failingSelfHostedStore) Save(_ context.Context, _ *models.SelfHostedConfig) error {
	s.saved = true
	return nil
}

func TestSelfHostedSecretsSurviveRestart(t *testing.T) {
	ctx := context.Background()
	store := &memorySelfHostedStore{}
	first := &types.Config{}
	if err := PrepareSelfHosted(ctx, first, store); err != nil {
		t.Fatal(err)
	}
	secrets := []string{first.Security.JWTSecret, first.Security.RefreshSecret, first.Telemetry.Salt}
	for _, secret := range secrets {
		decoded, err := hex.DecodeString(secret)
		if err != nil || len(decoded) != 32 {
			t.Fatal("expected a 256-bit secret")
		}
	}
	if secrets[0] == secrets[1] || secrets[0] == secrets[2] || secrets[1] == secrets[2] {
		t.Fatal("secrets must be independent")
	}
	if err := SaveSelfHostedOptions(ctx, first, store, "apps.example.com", "owner@example.com"); err != nil {
		t.Fatal(err)
	}
	second := &types.Config{}
	if err := PrepareSelfHosted(ctx, second, store); err != nil {
		t.Fatal(err)
	}
	if first.Security.JWTSecret != second.Security.JWTSecret || first.Security.RefreshSecret != second.Security.RefreshSecret || first.Telemetry.Salt != second.Telemetry.Salt {
		t.Fatal("restart changed authentication secrets")
	}
	if second.Domains.WildcardDomain != "apps.example.com" || second.Security.TLSEmail != "owner@example.com" {
		t.Fatal("setup options were not restored")
	}
}

func TestSelfHostedExplicitSecretsTakePrecedence(t *testing.T) {
	ctx := context.Background()
	store := &memorySelfHostedStore{}
	cfg := &types.Config{Security: types.SecurityConfig{JWTSecret: "existing-jwt", RefreshSecret: "existing-refresh"}}
	if err := PrepareSelfHosted(ctx, cfg, store); err != nil {
		t.Fatal(err)
	}
	if cfg.Security.JWTSecret != "existing-jwt" || cfg.Security.RefreshSecret != "existing-refresh" {
		t.Fatal("explicit secrets were replaced")
	}
	cfg.Security.JWTSecret = "override"
	if err := PrepareSelfHosted(ctx, cfg, store); err != nil {
		t.Fatal(err)
	}
	if cfg.Security.JWTSecret != "override" {
		t.Fatal("explicit override was ignored")
	}
}

func TestSelfHostedStorageFailureIsNotOverwritten(t *testing.T) {
	store := &failingSelfHostedStore{}
	if err := PrepareSelfHosted(context.Background(), &types.Config{}, store); err == nil {
		t.Fatal("storage failure must fail")
	}
	if store.saved {
		t.Fatal("failed load was overwritten")
	}
}

func TestCloudDoesNotGenerateSelfHostedSecrets(t *testing.T) {
	cfg := &types.Config{Cloud: types.CloudConfig{Enabled: true}}
	if err := PrepareSelfHosted(context.Background(), cfg, nil); err != nil {
		t.Fatal(err)
	}
	if cfg.Security.JWTSecret != "" {
		t.Fatal("cloud credentials must be configured by the operator")
	}
}

func TestSelfHostedOptionsCanBeCleared(t *testing.T) {
	ctx := context.Background()
	store := &memorySelfHostedStore{}
	cfg := &types.Config{}
	if err := PrepareSelfHosted(ctx, cfg, store); err != nil {
		t.Fatal(err)
	}
	if err := SaveSelfHostedOptions(ctx, cfg, store, "apps.example.com", "owner@example.com"); err != nil {
		t.Fatal(err)
	}
	if err := SaveSelfHostedOptions(ctx, cfg, store, "", ""); err != nil {
		t.Fatal(err)
	}
	restarted := &types.Config{}
	if err := PrepareSelfHosted(ctx, restarted, store); err != nil {
		t.Fatal(err)
	}
	if restarted.Domains.WildcardDomain != "" || restarted.Security.TLSEmail != "" {
		t.Fatal("cleared routing options restored on restart")
	}
	if restarted.Security.JWTSecret != cfg.Security.JWTSecret {
		t.Fatal("clearing options changed the auth secret")
	}
}
