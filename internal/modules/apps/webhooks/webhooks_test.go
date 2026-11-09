package webhooks

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

type mockWebhookService struct {
	webhooks map[string]*Webhook
	drains   map[string]*LogDrain
}

func newMockWebhookService() *mockWebhookService {
	return &mockWebhookService{
		webhooks: make(map[string]*Webhook),
		drains:   make(map[string]*LogDrain),
	}
}

func (m *mockWebhookService) CreateWebhook(ctx context.Context, w *Webhook) (*Webhook, error) {
	if w.ID == "" {
		w.ID = "wh-1"
	}
	m.webhooks[w.ID] = w
	return w, nil
}

func (m *mockWebhookService) ListWebhooks(ctx context.Context, serviceID string) ([]*Webhook, error) {
	var list []*Webhook
	for _, w := range m.webhooks {
		if w.ServiceID == serviceID {
			list = append(list, w)
		}
	}
	return list, nil
}

func (m *mockWebhookService) DeleteWebhook(ctx context.Context, id, serviceID string) error {
	delete(m.webhooks, id)
	return nil
}

func (m *mockWebhookService) CreateLogDrain(ctx context.Context, d *LogDrain) (*LogDrain, error) {
	if d.ID == "" {
		d.ID = "drain-1"
	}
	m.drains[d.ID] = d
	return d, nil
}

func (m *mockWebhookService) ListLogDrains(ctx context.Context, serviceID string) ([]*LogDrain, error) {
	var list []*LogDrain
	for _, d := range m.drains {
		if d.ServiceID == serviceID {
			list = append(list, d)
		}
	}
	return list, nil
}

func (m *mockWebhookService) DeleteLogDrain(ctx context.Context, id, serviceID string) error {
	delete(m.drains, id)
	return nil
}

func TestWebhookAndLogDrainRoutes(t *testing.T) {
	e := echo.New()
	svc := newMockWebhookService()
	h := NewWebhookHandler(svc)

	g := e.Group("")
	RegisterRoutes(g, h)

	// Create Webhook
	whBody := `{"url":"https://example.com/hook","secret":"s3cret","events":"push"}`
	req := httptest.NewRequest(http.MethodPost, "/apps/srv-1/webhooks", strings.NewReader(whBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}

	// List Webhooks
	req = httptest.NewRequest(http.MethodGet, "/apps/srv-1/webhooks", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var whListResp struct {
		Data []*Webhook `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &whListResp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(whListResp.Data) != 1 {
		t.Fatalf("expected 1 webhook, got %d", len(whListResp.Data))
	}

	// Delete Webhook
	req = httptest.NewRequest(http.MethodDelete, "/apps/srv-1/webhooks/wh-1", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}

	// Create Log Drain
	drainBody := `{"projectId":"proj-1","drainType":"http","endpointUrl":"https://logs.example.com"}`
	req = httptest.NewRequest(http.MethodPost, "/apps/srv-1/log-drains", strings.NewReader(drainBody))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}

	// List Log Drains
	req = httptest.NewRequest(http.MethodGet, "/apps/srv-1/log-drains", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	// Delete Log Drain
	req = httptest.NewRequest(http.MethodDelete, "/apps/srv-1/log-drains/drain-1", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
}
