package dns

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestDNSHandler_Create_InvalidPayload(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/dns", strings.NewReader("bad json"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	repo := newMemoryDNSRepo()
	svc := NewService(repo, nil)
	h := NewHandler(svc, nil)

	err := h.Create(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", rec.Code)
	}
}

func TestDNSHandler_Create_Success(t *testing.T) {
	e := echo.New()
	body := `{"domainName":"codedock.io","recordType":"A","recordName":"api","recordValue":"10.0.0.1","ttl":300}`
	req := httptest.NewRequest(http.MethodPost, "/dns", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	repo := newMemoryDNSRepo()
	svc := NewService(repo, nil)
	h := NewHandler(svc, nil)

	err := h.Create(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", rec.Code)
	}
}

func TestDNSHandler_Delete_MissingID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/dns/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("")

	repo := newMemoryDNSRepo()
	svc := NewService(repo, nil)
	h := NewHandler(svc, nil)

	err := h.Delete(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", rec.Code)
	}
}

func TestDNSHandler_List_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/dns?domain=codedock.io", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	repo := newMemoryDNSRepo()
	_ = repo.Create(context.Background(), &DNSRecord{
		ID:          "rec-1",
		DomainName:  "codedock.io",
		RecordType:  "A",
		RecordName:  "api",
		RecordValue: "10.0.0.1",
	})
	svc := NewService(repo, nil)
	h := NewHandler(svc, nil)

	err := h.List(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}
}
