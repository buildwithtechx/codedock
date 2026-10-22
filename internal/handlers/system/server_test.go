package system

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"codedock/internal/models"
)

type testServerService struct {
	servers []*models.Server
}

func (s testServerService) CreateServer(context.Context, string, models.CreateServerRequest) (*models.Server, error) {
	return nil, nil
}

func (s testServerService) TestSSH(context.Context, models.TestSSHRequest) error {
	return nil
}

func (s testServerService) ListServersByUser(context.Context, string, string) ([]*models.Server, error) {
	return s.servers, nil
}

func (s testServerService) GetServer(ctx context.Context, id string, userID string) (*models.Server, error) {
	for _, srv := range s.servers {
		if srv.ID == id {
			return srv, nil
		}
	}
	return nil, nil
}

func (s testServerService) UpdateServer(context.Context, string, string, models.UpdateServerRequest) (*models.Server, error) {
	return nil, nil
}

func (s testServerService) DeleteServer(context.Context, string, string) error {
	return nil
}

func (s testServerService) InspectSwarm(context.Context, string, string) (*models.SwarmInspection, error) {
	return &models.SwarmInspection{
		State: models.SwarmStateActive,
	}, nil
}

func TestServerListReturnsControlPlane(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/servers", nil)
	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)
	ctx.Set("user", &models.UserClaims{UserID: "owner", Role: models.UserRoleOwner})
	handler := NewServerHandler(testServerService{servers: []*models.Server{{
		ID:             "codedock-control-plane",
		Name:           "Codedock Control Plane",
		Status:         models.ServerStatusOnline,
		IsControlPlane: true,
	}}}, nil)
	if err := handler.List(ctx); err != nil {
		t.Fatalf("list servers: %v", err)
	}
	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.Code)
	}
	if body := res.Body.String(); body == "" || !strings.Contains(body, `"isControlPlane":true`) {
		t.Fatalf("expected control plane response, got %s", body)
	}
}

func TestStreamMonitorValidatesFormatAndStreams(t *testing.T) {
	e := echo.New()
	handler := NewServerHandler(testServerService{servers: []*models.Server{{
		ID:     "srv-1",
		Name:   "Test Server",
		Status: models.ServerStatusOnline,
	}}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/system/monitor/stream", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("user", &models.UserClaims{UserID: "u1"})
	_ = handler.StreamMonitor(ctx)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing serverId, got %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/system/monitor/stream?serverId=srv-1&format=unknown", nil)
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	ctx.Set("user", &models.UserClaims{UserID: "u1"})
	_ = handler.StreamMonitor(ctx)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid format, got %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/system/monitor/stream?serverId=missing", nil)
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	ctx.Set("user", &models.UserClaims{UserID: "u1"})
	_ = handler.StreamMonitor(ctx)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for unknown server, got %d", rec.Code)
	}

	reqContext, cancel := context.WithCancel(context.Background())
	cancel()
	req = httptest.NewRequest(http.MethodGet, "/api/system/monitor/stream?serverId=srv-1&format=usage", nil).WithContext(reqContext)
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	ctx.Set("user", &models.UserClaims{UserID: "u1"})
	if err := handler.StreamMonitor(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "event: usage") || !strings.Contains(body, "measuredAt") {
		t.Fatalf("expected usage event in body, got %s", body)
	}
}
