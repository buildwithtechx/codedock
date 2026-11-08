package environments

import (
	"context"
	"testing"
	"time"
)

type mockEnvRepo struct {
	envs []EnvironmentConfig
}

func (m *mockEnvRepo) Get(ctx context.Context, id string) (*EnvironmentConfig, error) {
	for _, e := range m.envs {
		if e.ID == id {
			return &e, nil
		}
	}
	return nil, ErrEnvironmentNotFound
}

func (m *mockEnvRepo) ListByProject(ctx context.Context, projectID string) ([]EnvironmentConfig, error) {
	var res []EnvironmentConfig
	for _, e := range m.envs {
		if e.ProjectID == projectID {
			res = append(res, e)
		}
	}
	return res, nil
}

func (m *mockEnvRepo) Create(ctx context.Context, env *EnvironmentConfig) error {
	m.envs = append(m.envs, *env)
	return nil
}

func (m *mockEnvRepo) Delete(ctx context.Context, id string) error {
	for i, e := range m.envs {
		if e.ID == id {
			m.envs = append(m.envs[:i], m.envs[i+1:]...)
			return nil
		}
	}
	return nil
}

func TestEnvironmentServiceCRUD(t *testing.T) {
	repo := &mockEnvRepo{}
	svc := NewEnvironmentService(repo, nil)
	ctx := context.Background()

	// Validation
	if _, err := svc.CreateEnvironment(ctx, nil); err == nil {
		t.Fatal("expected error on nil env")
	}

	created, err := svc.CreateEnvironment(ctx, &EnvironmentConfig{
		ProjectID: "proj-123",
		Name:      "production",
		IsDefault: true,
	})
	if err != nil {
		t.Fatalf("unexpected error creating env: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected generated UUID")
	}

	fetched, err := svc.GetEnvironment(ctx, created.ID)
	if err != nil || fetched.Name != "production" {
		t.Fatalf("failed to get env: %v", err)
	}

	list, err := svc.ListByProject(ctx, "proj-123")
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 env in list, got %d", len(list))
	}

	if err := svc.DeleteEnvironment(ctx, created.ID); err != nil {
		t.Fatalf("failed to delete env: %v", err)
	}
	listAfter, _ := svc.ListByProject(ctx, "proj-123")
	if len(listAfter) != 0 {
		t.Fatalf("expected 0 envs after delete, got %d", len(listAfter))
	}
}

func TestRouteProtectionPolicyAndTokens(t *testing.T) {
	password := "SecretPassword123!"
	policy, err := CreatePolicy(password, 3600)
	if err != nil {
		t.Fatalf("failed to create policy: %v", err)
	}

	if !VerifyPassword(password, policy) {
		t.Fatal("expected password verification to succeed")
	}
	if VerifyPassword("WrongPassword123!", policy) {
		t.Fatal("expected wrong password to fail")
	}

	hostname := "app.example.com"
	expires := time.Now().Unix() + 500 // future
	token := SignSessionToken(policy, hostname, expires)

	if !VerifySessionToken(policy, hostname, token) {
		t.Fatal("session token should be valid")
	}
	if VerifySessionToken(policy, "other.example.com", token) {
		t.Fatal("session token should fail on mismatched hostname")
	}

	path := "/dashboard"
	csrf := SignCSRFToken(policy, hostname, path, expires)
	if !VerifyCSRFToken(policy, hostname, path, csrf) {
		t.Fatal("csrf token should be valid")
	}
	if VerifyCSRFToken(policy, hostname, "/settings", csrf) {
		t.Fatal("csrf token should fail on mismatched path")
	}
}

func TestEnvSuggestionParser(t *testing.T) {
	svc := NewEnvSuggestionService(nil)
	input := `
# Configuration
PORT=8080
DATABASE_URL="postgres://user:pass@localhost:5432/db"
NODE_ENV='production'
export API_KEY=secret_key
INVALID-VAR=123
`
	suggestions := svc.parseEnvExampleText(input, ".env.example")
	if len(suggestions) != 4 {
		t.Fatalf("expected 4 valid suggestions, got %d", len(suggestions))
	}
	expectedKeys := map[string]string{
		"PORT":         "8080",
		"DATABASE_URL": "postgres://user:pass@localhost:5432/db",
		"NODE_ENV":     "production",
		"API_KEY":      "secret_key",
	}
	for _, s := range suggestions {
		expVal, ok := expectedKeys[s.Key]
		if !ok {
			t.Fatalf("unexpected suggested key: %s", s.Key)
		}
		if s.Value != expVal {
			t.Fatalf("expected value %s for %s, got %s", expVal, s.Key, s.Value)
		}
	}
}
