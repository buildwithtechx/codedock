package system

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/models"
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

func (s testServerService) ListServersByUser(context.Context, string) ([]*models.Server, error) {
	return s.servers, nil
}

func (s testServerService) GetServer(context.Context, string) (*models.Server, error) {
	return nil, nil
}

func (s testServerService) DeleteServer(context.Context, string, string) error {
	return nil
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
	}}})
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
