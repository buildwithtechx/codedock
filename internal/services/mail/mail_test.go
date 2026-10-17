package mail

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"codedock/internal/models"
	"codedock/internal/repositories"
)

type stubMailRepo struct {
	server *models.MailServer
}

func (s *stubMailRepo) Create(ctx context.Context, server *models.MailServer) error {
	s.server = server
	return nil
}

func (s *stubMailRepo) GetByID(ctx context.Context, id string) (*models.MailServer, error) {
	return s.server, nil
}

func (s *stubMailRepo) ListByOrganization(ctx context.Context, orgID string) ([]*models.MailServer, error) {
	return []*models.MailServer{s.server}, nil
}

func (s *stubMailRepo) Delete(ctx context.Context, id string) error {
	return nil
}

var _ repositories.MailRepository = (*stubMailRepo)(nil)

type capturedMailRequest struct {
	method string
	path   string
	auth   string
	body   map[string]any
}

func newMailStub(t *testing.T, status int, response string, captured *capturedMailRequest) *MailService {
	t.Helper()
	var stubServer *httptest.Server
	stubServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured.method = r.Method
		captured.path = r.URL.EscapedPath()
		captured.auth = r.Header.Get("Authorization")
		if r.Body != nil {
			raw, _ := io.ReadAll(r.Body)
			if len(raw) > 0 {
				_ = json.Unmarshal(raw, &captured.body)
			}
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(response))
	}))
	t.Cleanup(stubServer.Close)
	return NewMailService(&stubMailRepo{server: &models.MailServer{
		ID:       "srv-1",
		MaildURL: stubServer.URL,
		APIToken: "token-123",
	}})
}

func TestCreateMailboxSendsBearerToken(t *testing.T) {
	var captured capturedMailRequest
	service := newMailStub(t, 201, `{"id":"m-1","address":"alice@example.com","password":"once"}`, &captured)
	credentials, err := service.CreateMailbox(context.Background(), "srv-1", "d-1", "alice")
	if err != nil {
		t.Fatalf("create mailbox: %v", err)
	}
	if credentials.Password != "once" || credentials.Address != "alice@example.com" {
		t.Fatalf("unexpected credentials %+v", credentials)
	}
	if captured.auth != "Bearer token-123" {
		t.Fatalf("expected bearer token, got %q", captured.auth)
	}
	if captured.body["localPart"] != "alice" {
		t.Fatalf("unexpected body %+v", captured.body)
	}
}

func TestCreateDomainConflict(t *testing.T) {
	var captured capturedMailRequest
	service := newMailStub(t, 409, `{"error":"already exists"}`, &captured)
	if _, err := service.CreateDomain(context.Background(), "srv-1", "example.com"); !errors.Is(err, ErrMailExists) {
		t.Fatalf("expected ErrMailExists, got %v", err)
	}
}

func TestDeleteMapsNotFound(t *testing.T) {
	var captured capturedMailRequest
	service := newMailStub(t, 404, `{"error":"not found"}`, &captured)
	if err := service.DeleteDomain(context.Background(), "srv-1", "missing"); !errors.Is(err, ErrMailNotFound) {
		t.Fatalf("expected ErrMailNotFound, got %v", err)
	}
}

func TestDeleteNoContent(t *testing.T) {
	var captured capturedMailRequest
	service := newMailStub(t, 204, ``, &captured)
	if err := service.DeleteMailbox(context.Background(), "srv-1", "d-1", "m-1"); err != nil {
		t.Fatalf("delete mailbox: %v", err)
	}
	if captured.method != http.MethodDelete || captured.path != "/api/v1/domains/d-1/mailboxes/m-1" {
		t.Fatalf("unexpected request %s %s", captured.method, captured.path)
	}
}

func TestCreateServerValidation(t *testing.T) {
	service := NewMailService(&stubMailRepo{})
	if _, err := service.CreateServer(context.Background(), &models.CreateMailServerRequest{}); err == nil {
		t.Fatalf("expected validation error")
	}
}
