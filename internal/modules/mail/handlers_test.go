package mail

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

type memoryMailRepo struct {
	servers map[string]*MailServer
}

func newMemoryMailRepo() *memoryMailRepo {
	return &memoryMailRepo{servers: make(map[string]*MailServer)}
}

func (m *memoryMailRepo) Create(ctx context.Context, s *MailServer) error {
	m.servers[s.ID] = s
	return nil
}

func (m *memoryMailRepo) GetByID(ctx context.Context, id string) (*MailServer, error) {
	s, ok := m.servers[id]
	if !ok {
		return nil, ErrMailNotFound
	}
	return s, nil
}

func (m *memoryMailRepo) ListByOrganization(ctx context.Context, orgID string) ([]*MailServer, error) {
	var list []*MailServer
	for _, s := range m.servers {
		if s.OrganizationID == orgID {
			list = append(list, s)
		}
	}
	return list, nil
}

func (m *memoryMailRepo) Delete(ctx context.Context, id string) error {
	delete(m.servers, id)
	return nil
}

func TestMailHandler_CreateServer_InvalidPayload(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/mail/servers", strings.NewReader("bad json"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	repo := newMemoryMailRepo()
	svc := NewService(repo)
	h := NewHandler(svc, nil, nil)

	err := h.CreateServer(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", rec.Code)
	}
}

func TestMailHandler_DeleteServer_NotFound(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/mail/servers/missing", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("missing")

	repo := newMemoryMailRepo()
	svc := NewService(repo)
	h := NewHandler(svc, nil, nil)

	err := h.DeleteServer(c)
	if err == nil {
		t.Fatalf("expected error for missing server, got nil")
	}
	if !errors.Is(err, ErrMailNotFound) {
		t.Fatalf("expected ErrMailNotFound, got: %v", err)
	}
}
