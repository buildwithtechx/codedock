package operations

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

type mockOpService struct {
	ops []Operation
}

type testOpClaims struct {
	UserID string
	Email  string
}

func (m *mockOpService) Get(ctx context.Context, id string) (*Operation, error) {
	for _, o := range m.ops {
		if o.ID == id {
			return &o, nil
		}
	}
	return nil, nil
}

func (m *mockOpService) List(ctx context.Context, project, user string) ([]Operation, error) {
	return m.ops, nil
}

func (m *mockOpService) Cancel(ctx context.Context, id, user string) error {
	return nil
}

func TestOperationHandler_List(t *testing.T) {
	svc := &mockOpService{
		ops: []Operation{
			{ID: "op-1", UserID: "usr-1", Status: "RUNNING"},
		},
	}
	h := NewOperationHandler(svc)
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/operations", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &testOpClaims{UserID: "usr-1", Email: "test@example.com"})

	if err := h.List(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}
}

func TestOperationHandler_Get(t *testing.T) {
	svc := &mockOpService{
		ops: []Operation{
			{ID: "op-1", UserID: "usr-1", Status: "COMPLETED"},
		},
	}
	h := NewOperationHandler(svc)
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/operations/op-1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("operationId")
	c.SetParamValues("op-1")
	c.Set("user", &testOpClaims{UserID: "usr-1", Email: "test@example.com"})

	if err := h.Get(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}
}
