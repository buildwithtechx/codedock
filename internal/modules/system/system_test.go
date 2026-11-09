package system

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"codedock/internal/modules/settings"
)

type mockSettingsRepo struct {
	settings *settings.ServerSettings
	err      error
}

func (m *mockSettingsRepo) GetServerSettings(ctx context.Context) (*settings.ServerSettings, error) {
	return m.settings, m.err
}

func (m *mockSettingsRepo) UpdateServerSettings(ctx context.Context, cfg *settings.ServerSettings) error {
	m.settings = cfg
	return m.err
}

type mockSystemOps struct {
	stats any
	err   error
}

func (m *mockSystemOps) GetStats(ctx context.Context) (any, error) {
	return m.stats, m.err
}

func (m *mockSystemOps) Restart(ctx context.Context) error {
	return m.err
}

func (m *mockSystemOps) Cleanup(ctx context.Context) error {
	return m.err
}

type mockUpdaterOps struct {
	status any
	err    error
}

func (m *mockUpdaterOps) GetStatus() any {
	return m.status
}

func (m *mockUpdaterOps) CheckUpdate(ctx context.Context) (any, error) {
	return m.status, m.err
}

func (m *mockUpdaterOps) ApplyUpdate(ctx context.Context) (any, error) {
	return m.status, m.err
}

func TestSystemService_GetStats(t *testing.T) {
	svc := NewSystemService()
	stats, err := svc.GetStats()
	assert.NoError(t, err)
	assert.NotNil(t, stats)
}

func TestSystemHandler_GetStats(t *testing.T) {
	ops := &mockSystemOps{stats: &SystemStats{Processes: 10}}
	h := NewSystemHandler(ops, nil, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/system/stats", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.GetStats(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUpdaterHandler_GetUpdateStatus(t *testing.T) {
	ops := &mockUpdaterOps{status: &UpdateInfo{CurrentVersion: "v1.0.0"}}
	h := NewUpdaterHandler(ops)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/settings/updates/status", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.GetUpdateStatus(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestIsNewerVersion(t *testing.T) {
	assert.True(t, isNewerVersion("v1.0.0", "v1.1.0"))
	assert.False(t, isNewerVersion("v1.0.0", "v1.0.0"))
}
