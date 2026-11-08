package migrations

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

type mockBundleService struct {
	exported []byte
	imported any
}

func (m *mockBundleService) Export(ctx context.Context, passphrase string) ([]byte, error) {
	return m.exported, nil
}

func (m *mockBundleService) Import(ctx context.Context, data []byte, passphrase string) (any, error) {
	return m.imported, nil
}

type mockUserCounter struct {
	count int
}

func (m *mockUserCounter) CountUsers(ctx context.Context) (int, error) {
	return m.count, nil
}

func TestBundleHandler_ExportRequiresPassphrase(t *testing.T) {
	e := echo.New()
	handler := NewBundleHandler(&mockBundleService{}, &mockUserCounter{}, nil)

	req := httptest.NewRequest(http.MethodPost, "/system/export", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.Export(c); err != nil {
		t.Fatalf("unexpected handler error: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestLifecycleHandler_ListSourcesRequiresAuth(t *testing.T) {
	e := echo.New()
	handler := NewLifecycleHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/organizations/org-1/migration/sources", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.ListSources(c); err != nil {
		t.Fatalf("unexpected handler error: %v", err)
	}
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
}
