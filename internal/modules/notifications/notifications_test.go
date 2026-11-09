package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

type mockNotifService struct {
	lastOrgID string
}

type testSupportClaims struct {
	UserID string
	Email  string
	Role   string
}

func (m *mockNotifService) GetNotificationSettings(ctx context.Context) (*NotificationSettings, error) {
	return &NotificationSettings{}, nil
}
func (m *mockNotifService) UpdateNotificationSettings(ctx context.Context, settings *NotificationSettings) error {
	return nil
}
func (m *mockNotifService) ListSubscriptions(ctx context.Context, userID, orgID string) (any, error) {
	return []any{}, nil
}
func (m *mockNotifService) UpsertSubscription(ctx context.Context, userID, orgID, category string, channels []string, enabled bool) (any, error) {
	m.lastOrgID = orgID
	return map[string]any{"status": "ok"}, nil
}
func (m *mockNotifService) ListDefaults(ctx context.Context, orgID string) (any, error) {
	return []any{}, nil
}
func (m *mockNotifService) UpsertDefault(ctx context.Context, orgID, category string, channels []string, enabled bool) (any, error) {
	return map[string]any{"status": "ok"}, nil
}
func (m *mockNotifService) ListDeliveries(ctx context.Context, limit int) (any, error) {
	return []any{}, nil
}
func (m *mockNotifService) RegisterDeviceToken(ctx context.Context, userID, token, platform string) (any, error) {
	return map[string]any{"status": "ok"}, nil
}
func (m *mockNotifService) DeleteDeviceToken(ctx context.Context, userID, token string) error {
	return nil
}
func (m *mockNotifService) ListDeviceTokens(ctx context.Context, userID string) (any, error) {
	return []any{}, nil
}
func (m *mockNotifService) ResolveOrgID(ctx context.Context, userID string) string {
	if userID == "usr-no-org" {
		return ""
	}
	return "org-solomon-personal"
}

func TestNotificationSettingsHandler_UpsertSubscriptionResolvesOrg(t *testing.T) {
	svc := &mockNotifService{}
	handler := NewNotificationSettingsHandler(svc)

	e := echo.New()
	body, _ := json.Marshal(UpsertSubscriptionRequest{
		Category: "deploy.failed",
		Channels: []string{"email"},
		Enabled:  true,
	})

	req := httptest.NewRequest(http.MethodPut, "/notifications/subscriptions", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &testSupportClaims{UserID: "usr-123", Email: "solomon@example.com"})

	if err := handler.UpsertSubscription(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
	}

	if svc.lastOrgID != "org-solomon-personal" {
		t.Fatalf("expected org-solomon-personal, got %q", svc.lastOrgID)
	}
}

func TestNotificationSettingsHandler_RejectsDefaultOrgWhenUnresolved(t *testing.T) {
	svc := &mockNotifService{}
	handler := NewNotificationSettingsHandler(svc)

	e := echo.New()
	body, _ := json.Marshal(UpsertSubscriptionRequest{
		Category: "deploy.failed",
		Channels: []string{"email"},
		Enabled:  true,
	})

	req := httptest.NewRequest(http.MethodPut, "/notifications/subscriptions?orgId=default", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &testSupportClaims{UserID: "usr-no-org", Email: "test@example.com"})

	_ = handler.UpsertSubscription(c)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d: %s", rec.Code, rec.Body.String())
	}
}
