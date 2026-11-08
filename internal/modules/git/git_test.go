package git

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockEncryptor struct{}

func (m *mockEncryptor) Encrypt(plain string) (string, error) {
	return "enc_" + plain, nil
}

func (m *mockEncryptor) Decrypt(cipher string) (string, error) {
	if len(cipher) > 4 && cipher[:4] == "enc_" {
		return cipher[4:], nil
	}
	return cipher, nil
}

type mockProviderRepo struct {
	providers map[string]*GitProviderConfig
}

func newMockProviderRepo() *mockProviderRepo {
	return &mockProviderRepo{providers: make(map[string]*GitProviderConfig)}
}

func (m *mockProviderRepo) SaveProvider(ctx context.Context, gp *GitProviderConfig) error {
	m.providers[gp.UserID+":"+gp.Provider] = gp
	return nil
}

func (m *mockProviderRepo) GetProvider(ctx context.Context, userID, provider string) (*GitProviderConfig, error) {
	gp, ok := m.providers[userID+":"+provider]
	if !ok {
		return nil, nil
	}
	return gp, nil
}

func (m *mockProviderRepo) GetAnyProviderByType(ctx context.Context, provider string) (*GitProviderConfig, error) {
	for _, gp := range m.providers {
		if gp.Provider == provider {
			return gp, nil
		}
	}
	return nil, ErrProviderNotFound
}

func (m *mockProviderRepo) ListProvidersByUser(ctx context.Context, userID string) ([]*GitProviderConfig, error) {
	var list []*GitProviderConfig
	for _, gp := range m.providers {
		if gp.UserID == userID {
			list = append(list, gp)
		}
	}
	return list, nil
}

func (m *mockProviderRepo) DeleteProvider(ctx context.Context, userID, provider string) error {
	delete(m.providers, userID+":"+provider)
	return nil
}

func TestProviderService_ConnectAndStatus(t *testing.T) {
	repo := newMockProviderRepo()
	svc := NewProviderService(repo)

	// Invalid provider
	_, err := svc.ConnectProvider(context.Background(), "u1", &GitConnectRequest{
		Provider:    "invalid",
		AccessToken: "tok123",
	})
	assert.ErrorIs(t, err, ErrUnsupportedProvider)

	// Missing access token
	_, err = svc.ConnectProvider(context.Background(), "u1", &GitConnectRequest{
		Provider: ProviderGitHub,
	})
	assert.ErrorIs(t, err, ErrAccessTokenRequired)

	// Valid connect
	gp, err := svc.ConnectProvider(context.Background(), "u1", &GitConnectRequest{
		Provider:    ProviderGitHub,
		AccessToken: "ghp_secret",
		AccountName: "octocat",
	})
	require.NoError(t, err)
	assert.Equal(t, ProviderGitHub, gp.Provider)
	assert.Equal(t, "octocat", gp.AccountName)

	// Check status
	statuses, err := svc.GetConnectedProviders(context.Background(), "u1")
	require.NoError(t, err)
	assert.Len(t, statuses, 4)

	var ghStatus *ConnectedProviderStatus
	for i := range statuses {
		if statuses[i].Provider == ProviderGitHub {
			ghStatus = &statuses[i]
		}
	}
	require.NotNil(t, ghStatus)
	assert.True(t, ghStatus.Connected)
	assert.Equal(t, "octocat", ghStatus.AccountName)

	// Disconnect
	err = svc.DisconnectProvider(context.Background(), "u1", ProviderGitHub)
	require.NoError(t, err)

	statusesAfter, err := svc.GetConnectedProviders(context.Background(), "u1")
	require.NoError(t, err)
	for _, s := range statusesAfter {
		if s.Provider == ProviderGitHub {
			assert.False(t, s.Connected)
		}
	}
}

func TestSubmodules_ParseSameOwner(t *testing.T) {
	content := `
[submodule "libs/common"]
	path = libs/common
	url = git@github.com:myorg/common.git
[submodule "libs/external"]
	path = libs/external
	url = https://github.com/otherorg/ext.git
[submodule "libs/relative"]
	path = libs/relative
	url = ../core.git
`
	repos := ParseSameOwnerSubmoduleRepos(content, "myorg", "parent-app", "github.com")
	assert.Contains(t, repos, "common")
	assert.Contains(t, repos, "core")
	assert.NotContains(t, repos, "ext")
}

func TestChangeImpact_Resolve(t *testing.T) {
	// Force deploy in commit msg
	imp := ResolveProjectChangeImpact("apps/web", "", "", []string{"docs/readme.md"}, "[force-deploy] update docs")
	assert.True(t, imp.Affected)
	assert.Equal(t, "force-deploy-token", imp.Reason)

	// Root app always affected
	imp = ResolveProjectChangeImpact("", "", "", []string{"docs/readme.md"}, "normal commit")
	assert.True(t, imp.Affected)
	assert.Equal(t, "repository-root", imp.Reason)

	// Global manifest affected
	imp = ResolveProjectChangeImpact("apps/web", "", "", []string{"pnpm-lock.yaml"}, "lock update")
	assert.True(t, imp.Affected)
	assert.Equal(t, "global-manifest", imp.Reason)

	// Affected files in subfolder
	imp = ResolveProjectChangeImpact("apps/web", "", "", []string{"apps/web/src/index.ts"}, "code change")
	assert.True(t, imp.Affected)
	assert.Equal(t, "project-files", imp.Reason)

	// Unrelated workspace changes
	imp = ResolveProjectChangeImpact("apps/web", "", "", []string{"apps/backend/main.go"}, "backend update")
	assert.False(t, imp.Affected)
	assert.Equal(t, "unrelated-workspaces", imp.Reason)
}

func TestPRPreview_Eligibility(t *testing.T) {
	validSHA := "0123456789abcdef0123456789abcdef01234567"
	pr := PullRequestInfo{
		Number:   12,
		State:    "open",
		Draft:    false,
		HeadSHA:  validSHA,
		HeadRef:  "feature/foo",
		HeadRepo: "octocat/hello-world",
		BaseRef:  "main",
		BaseRepo: "octocat/hello-world",
	}

	ok, reason := IsPreviewEligible(pr, "octocat", "hello-world", "main")
	assert.True(t, ok, reason)

	// Closed PR
	closedPR := pr
	closedPR.State = "closed"
	ok, _ = IsPreviewEligible(closedPR, "octocat", "hello-world", "main")
	assert.False(t, ok)

	// Draft PR
	draftPR := pr
	draftPR.Draft = true
	ok, _ = IsPreviewEligible(draftPR, "octocat", "hello-world", "main")
	assert.False(t, ok)

	// Wrong base ref
	wrongBase := pr
	wrongBase.BaseRef = "develop"
	ok, _ = IsPreviewEligible(wrongBase, "octocat", "hello-world", "main")
	assert.False(t, ok)

	// Invalid SHA
	badSHA := pr
	badSHA.HeadSHA = "invalid-sha"
	ok, _ = IsPreviewEligible(badSHA, "octocat", "hello-world", "main")
	assert.False(t, ok)
}

func TestPRPreview_GeneratePreviewDomain(t *testing.T) {
	domain := GeneratePreviewDomain("preview.codedock.io", "", "web-app", "proj-123", 42)
	assert.Contains(t, domain, "web-app-pr-42-")
	assert.Contains(t, domain, ".preview.codedock.io")

	domainMagic := GeneratePreviewDomain("", "sslip.io", "web-app", "proj-123", 42)
	assert.Contains(t, domainMagic, "web-app-pr-42-")
	assert.Contains(t, domainMagic, ".sslip.io")
}

func TestWebhook_HMAC(t *testing.T) {
	secret := "test-secret"
	payload := []byte(`{"action":"ping"}`)

	// Compute valid HMAC
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	sig := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	valid := VerifyHMAC(payload, secret, sig)
	assert.True(t, valid)

	invalid := VerifyHMAC(payload, "wrong-secret", sig)
	assert.False(t, invalid)
}

func TestRoutes_Register(t *testing.T) {
	e := echo.New()
	apiGroup := e.Group("/api")
	authGroup := e.Group("/auth")

	repo := newMockProviderRepo()
	svc := NewProviderService(repo)
	pHandler := NewProvidersHandler(svc)
	iHandler := NewInspectionHandler(svc)
	aHandler := NewAppsHandler(NewGitAppsService(nil))
	wHandler := NewWebhookHandler(nil, nil, nil, nil)

	RegisterRoutes(apiGroup, authGroup, pHandler, iHandler, aHandler, wHandler, RouteMiddleware{})

	routes := e.Routes()
	assert.NotEmpty(t, routes)
}
