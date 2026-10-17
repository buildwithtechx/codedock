package system

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"codedock/internal/models"
	systemservices "codedock/internal/services/system"
)

type mockOrgRepo struct {
	orgs []*models.Organization
}

func (m *mockOrgRepo) Create(ctx context.Context, org *models.Organization) error { return nil }
func (m *mockOrgRepo) CreateWithOwner(ctx context.Context, org *models.Organization, owner *models.OrganizationMember) error {
	return nil
}
func (m *mockOrgRepo) GetByID(ctx context.Context, id string) (*models.Organization, error) {
	return nil, nil
}
func (m *mockOrgRepo) ListByUser(ctx context.Context, userID string) ([]*models.Organization, error) {
	return m.orgs, nil
}
func (m *mockOrgRepo) Update(ctx context.Context, org *models.Organization) error { return nil }
func (m *mockOrgRepo) Delete(ctx context.Context, id string) error               { return nil }
func (m *mockOrgRepo) AddMember(ctx context.Context, member *models.OrganizationMember) error {
	return nil
}
func (m *mockOrgRepo) GetMember(ctx context.Context, orgID, userID string) (*models.OrganizationMember, error) {
	return nil, nil
}
func (m *mockOrgRepo) GetMemberByID(ctx context.Context, id string) (*models.OrganizationMember, error) {
	return nil, nil
}
func (m *mockOrgRepo) GetMemberByEmail(ctx context.Context, orgID, email string) (*models.OrganizationMember, error) {
	return nil, nil
}
func (m *mockOrgRepo) ListMembers(ctx context.Context, orgID string) ([]*models.OrganizationMember, error) {
	return nil, nil
}
func (m *mockOrgRepo) UpdateMember(ctx context.Context, member *models.OrganizationMember) error {
	return nil
}
func (m *mockOrgRepo) RemoveMember(ctx context.Context, id string) error { return nil }
func (m *mockOrgRepo) ListInvitesByEmail(ctx context.Context, email string) ([]*models.OrganizationMember, error) {
	return nil, nil
}

type mockNotifRepo struct {
	lastUpserted *models.NotificationSubscription
}

func (m *mockNotifRepo) GetNotificationSettings(ctx context.Context) (*models.NotificationSettings, error) {
	return &models.NotificationSettings{}, nil
}
func (m *mockNotifRepo) UpdateNotificationSettings(ctx context.Context, cfg *models.NotificationSettings) error {
	return nil
}
func (m *mockNotifRepo) ListSubscriptions(ctx context.Context, userID, orgID string) ([]models.NotificationSubscription, error) {
	return []models.NotificationSubscription{}, nil
}
func (m *mockNotifRepo) UpsertSubscription(ctx context.Context, sub *models.NotificationSubscription) error {
	m.lastUpserted = sub
	return nil
}
func (m *mockNotifRepo) ListDefaults(ctx context.Context, orgID string) ([]models.NotificationDefault, error) {
	return []models.NotificationDefault{}, nil
}
func (m *mockNotifRepo) UpsertDefault(ctx context.Context, def *models.NotificationDefault) error {
	return nil
}

func TestNotificationSettingsHandler_UpsertSubscriptionResolvesOrg(t *testing.T) {
	orgRepo := &mockOrgRepo{
		orgs: []*models.Organization{
			{ID: "org-solomon-personal", Name: "Solomon Workspace"},
		},
	}
	notifRepo := &mockNotifRepo{}
	svc := systemservices.NewNotificationSettingsService(notifRepo, nil)
	handler := NewNotificationSettingsHandler(svc, orgRepo)

	e := echo.New()
	body, _ := json.Marshal(models.UpsertSubscriptionRequest{
		Category: "deploy.failed",
		Channels: []string{"email"},
		Enabled:  true,
	})

	req := httptest.NewRequest(http.MethodPut, "/notifications/subscriptions", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &models.UserClaims{UserID: "usr-123", Email: "solomon@example.com"})

	if err := handler.UpsertSubscription(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
	}

	if notifRepo.lastUpserted == nil {
		t.Fatalf("expected subscription to be upserted")
	}

	if notifRepo.lastUpserted.OrganizationID != "org-solomon-personal" {
		t.Fatalf("expected org-solomon-personal, got %q", notifRepo.lastUpserted.OrganizationID)
	}
}

func TestNotificationSettingsHandler_RejectsDefaultOrgWhenUnresolved(t *testing.T) {
	orgRepo := &mockOrgRepo{orgs: []*models.Organization{}}
	notifRepo := &mockNotifRepo{}
	svc := systemservices.NewNotificationSettingsService(notifRepo, nil)
	handler := NewNotificationSettingsHandler(svc, orgRepo)

	e := echo.New()
	body, _ := json.Marshal(models.UpsertSubscriptionRequest{
		Category: "deploy.failed",
		Channels: []string{"email"},
		Enabled:  true,
	})

	req := httptest.NewRequest(http.MethodPut, "/notifications/subscriptions?orgId=default", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &models.UserClaims{UserID: "usr-no-org", Email: "test@example.com"})

	_ = handler.UpsertSubscription(c)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d: %s", rec.Code, rec.Body.String())
	}
}
