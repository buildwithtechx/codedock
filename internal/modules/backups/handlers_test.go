package backups

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterRoutes_NilSafety(t *testing.T) {
	e := echo.New()
	authGroup := e.Group("/api")

	assert.NotPanics(t, func() {
		RegisterRoutes(authGroup, nil, RouteMiddleware{})
	})
}

func TestBackupHandler_List(t *testing.T) {
	repo := newMockBackupConfigRepo()
	svc := NewService(repo, nil, nil)
	h := NewHandler(svc).SetSecurity(nil, func(ctx context.Context, projectID, userID, role, permission string) bool {
		return true
	}, nil)

	_ = svc.CreateConfig(context.Background(), &BackupConfig{
		ID:            "b-1",
		Name:          "Daily",
		DbPassword:    "secret",
		BackupEnabled: true,
	})

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/backups", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	type mockUser struct {
		UserID string
		Role   string
	}
	c.Set("user", &mockUser{UserID: "u-1", Role: "admin"})

	err := h.List(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Daily")
	assert.NotContains(t, rec.Body.String(), "secret")
}
