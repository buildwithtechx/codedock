package vars

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

type mockService struct {
	items map[string]*Variable
}

func newMockService() *mockService {
	return &mockService{items: make(map[string]*Variable)}
}

func (m *mockService) CreateVariable(ctx context.Context, v *Variable) (*Variable, error) {
	if v.ID == "" {
		v.ID = "var-1"
	}
	m.items[v.ID] = v
	return v, nil
}

func (m *mockService) UpdateVariable(ctx context.Context, v *Variable) error {
	m.items[v.ID] = v
	return nil
}

func (m *mockService) GetVariable(ctx context.Context, id string) (*Variable, error) {
	v, ok := m.items[id]
	if !ok {
		return nil, nil
	}
	return v, nil
}

func (m *mockService) ListVariablesByService(ctx context.Context, serviceID string) ([]*Variable, error) {
	var res []*Variable
	for _, v := range m.items {
		if v.ServiceID == serviceID {
			res = append(res, v)
		}
	}
	return res, nil
}

func (m *mockService) DeleteVariable(ctx context.Context, id string) error {
	delete(m.items, id)
	return nil
}

func TestVariableRoutes(t *testing.T) {
	e := echo.New()
	svc := newMockService()
	h := NewVariableHandler(svc, nil)

	g := e.Group("")
	RegisterRoutes(g, h)

	// Create
	body := `{"key":"PORT","value":"8080","isSecret":false}`
	req := httptest.NewRequest(http.MethodPost, "/services/srv-1/variables", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", rec.Code, rec.Body.String())
	}

	// List
	req = httptest.NewRequest(http.MethodGet, "/services/srv-1/variables", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var listResp struct {
		Data []*Variable `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &listResp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(listResp.Data) != 1 {
		t.Fatalf("expected 1 item, got %d", len(listResp.Data))
	}

	// Suggest
	req = httptest.NewRequest(http.MethodGet, "/services/srv-1/env-suggestions", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	// Delete
	req = httptest.NewRequest(http.MethodDelete, "/services/srv-1/variables/var-1", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}
}
