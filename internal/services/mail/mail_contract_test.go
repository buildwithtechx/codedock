package mail

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"codedock/internal/models"
)

type recordedCall struct {
	method string
	path   string
	body   []byte
}

func newMaildStub(t *testing.T, calls *[]recordedCall, bodies map[string][]byte) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	respond := func(w http.ResponseWriter, r *http.Request, status int, payload any) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("expected bearer auth, got %q", got)
		}
		var body []byte
		if r.Body != nil {
			body, _ = io.ReadAll(r.Body)
		}
		*calls = append(*calls, recordedCall{method: r.Method, path: r.URL.RequestURI(), body: body})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}
	register := func(base string) {
		mux.HandleFunc(base+"/domains", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				respond(w, r, http.StatusOK, []models.MailDomain{{ID: "dom-1", Name: "example.com"}})
				return
			}
			var req models.CreateMailDomainRequest
			raw, _ := io.ReadAll(r.Body)
			_ = json.Unmarshal(raw, &req)
			bodies["create-domain"] = raw
			respond(w, r, http.StatusCreated, models.MailDomain{ID: "dom-2", Name: req.Name})
		})
		mux.HandleFunc(base+"/domains/", func(w http.ResponseWriter, r *http.Request) {
			rest := strings.TrimPrefix(r.URL.Path, base+"/domains/")
			switch {
			case rest == "missing":
				respond(w, r, http.StatusNotFound, map[string]string{"error": "nope"})
			case strings.HasSuffix(rest, "/mailboxes") && r.Method == http.MethodGet:
				respond(w, r, http.StatusOK, []models.Mailbox{{ID: "mb-1", DomainID: "dom-1", LocalPart: "ops", Address: "ops@example.com"}})
			case strings.HasSuffix(rest, "/mailboxes") && r.Method == http.MethodPost:
				var req models.CreateMailboxRequest
				raw, _ := io.ReadAll(r.Body)
				_ = json.Unmarshal(raw, &req)
				bodies["create-mailbox"] = raw
				respond(w, r, http.StatusCreated, models.MailboxCredentials{Mailbox: models.Mailbox{ID: "mb-2", DomainID: "dom-1", LocalPart: req.LocalPart}, Password: "secret"})
			case strings.HasSuffix(rest, "/dkim"):
				respond(w, r, http.StatusCreated, models.DKIMKey{ID: "key-1", DomainID: "dom-1", Selector: "mail", Algorithm: "rsa"})
			case strings.HasSuffix(rest, "/dns"):
				respond(w, r, http.StatusOK, models.DNSResponse{Domain: "example.com", Records: []models.MailDNSRecord{}})
			case strings.HasSuffix(rest, "/checks"):
				respond(w, r, http.StatusOK, models.MailCheckReport{Domain: "example.com", Checks: []models.MailCheck{}})
			case r.Method == http.MethodDelete:
				respond(w, r, http.StatusNoContent, nil)
			default:
				respond(w, r, http.StatusNotFound, map[string]string{"error": "nope"})
			}
		})
		mux.HandleFunc(base+"/certs", func(w http.ResponseWriter, r *http.Request) {
			respond(w, r, http.StatusOK, []models.MailCert{{ID: "cert-1", Hostname: "mail.example.com"}})
		})
		mux.HandleFunc(base+"/certs/probe", func(w http.ResponseWriter, r *http.Request) {
			respond(w, r, http.StatusOK, models.MailCert{ID: "cert-1", Hostname: "mail.example.com"})
		})
	}
	register("/api/v1")
	return httptest.NewServer(mux)
}

func TestMaildContractPaths(t *testing.T) {
	var calls []recordedCall
	bodies := map[string][]byte{}
	maild := newMaildStub(t, &calls, bodies)
	defer maild.Close()
	service := NewMailService(&stubMailRepo{server: &models.MailServer{ID: "srv-1", MaildURL: maild.URL, APIToken: "test-token"}})
	ctx := context.Background()

	if _, err := service.ListDomains(ctx, "srv-1"); err != nil {
		t.Fatalf("list domains: %v", err)
	}
	if _, err := service.CreateDomain(ctx, "srv-1", "example.com"); err != nil {
		t.Fatalf("create domain: %v", err)
	}
	if err := service.DeleteDomain(ctx, "srv-1", "dom-2"); err != nil {
		t.Fatalf("delete domain: %v", err)
	}
	if _, err := service.ListMailboxes(ctx, "srv-1", "dom-1"); err != nil {
		t.Fatalf("list mailboxes: %v", err)
	}
	creds, err := service.CreateMailbox(ctx, "srv-1", "dom-1", "ops")
	if err != nil {
		t.Fatalf("create mailbox: %v", err)
	}
	if creds.Password != "secret" {
		t.Fatalf("expected credentials password, got %+v", creds)
	}
	if err := service.DeleteMailbox(ctx, "srv-1", "dom-1", "mb-2"); err != nil {
		t.Fatalf("delete mailbox: %v", err)
	}
	if _, err := service.CreateDKIM(ctx, "srv-1", "dom-1", "mail", "rsa"); err != nil {
		t.Fatalf("create dkim: %v", err)
	}
	if _, err := service.GetDNS(ctx, "srv-1", "dom-1"); err != nil {
		t.Fatalf("get dns: %v", err)
	}
	if _, err := service.GetChecks(ctx, "srv-1", "dom-1", "1.2.3.4"); err != nil {
		t.Fatalf("get checks: %v", err)
	}
	if _, err := service.ListCerts(ctx, "srv-1"); err != nil {
		t.Fatalf("list certs: %v", err)
	}
	if _, err := service.ProbeCert(ctx, "srv-1", "mail.example.com", 443); err != nil {
		t.Fatalf("probe cert: %v", err)
	}

	want := []recordedCall{
		{method: "GET", path: "/api/v1/domains"},
		{method: "POST", path: "/api/v1/domains"},
		{method: "DELETE", path: "/api/v1/domains/dom-2"},
		{method: "GET", path: "/api/v1/domains/dom-1/mailboxes"},
		{method: "POST", path: "/api/v1/domains/dom-1/mailboxes"},
		{method: "DELETE", path: "/api/v1/domains/dom-1/mailboxes/mb-2"},
		{method: "POST", path: "/api/v1/domains/dom-1/dkim"},
		{method: "GET", path: "/api/v1/domains/dom-1/dns"},
		{method: "GET", path: "/api/v1/domains/dom-1/checks?ip=1.2.3.4"},
		{method: "GET", path: "/api/v1/certs"},
		{method: "POST", path: "/api/v1/certs/probe"},
	}
	if len(calls) != len(want) {
		t.Fatalf("expected %d maild calls, got %d (%v)", len(want), len(calls), calls)
	}
	for i := range want {
		if calls[i].method != want[i].method || calls[i].path != want[i].path {
			t.Fatalf("call %d: expected %s %s, got %s %s", i, want[i].method, want[i].path, calls[i].method, calls[i].path)
		}
	}
	if !strings.Contains(string(bodies["create-domain"]), `"name":"example.com"`) {
		t.Fatalf("expected domain name in create body, got %s", bodies["create-domain"])
	}
	if !strings.Contains(string(bodies["create-mailbox"]), `"localPart":"ops"`) {
		t.Fatalf("expected localPart in mailbox body, got %s", bodies["create-mailbox"])
	}
}

func TestMaildContractNotFound(t *testing.T) {
	var calls []recordedCall
	maild := newMaildStub(t, &calls, map[string][]byte{})
	defer maild.Close()
	service := NewMailService(&stubMailRepo{server: &models.MailServer{ID: "srv-1", MaildURL: maild.URL, APIToken: "test-token"}})
	err := service.DeleteDomain(context.Background(), "srv-1", "missing")
	if err != ErrMailNotFound {
		t.Fatalf("expected not-found mapping, got %v", err)
	}
}

func TestMaildUnreachable(t *testing.T) {
	service := NewMailService(&stubMailRepo{server: &models.MailServer{ID: "srv-1", MaildURL: "http://127.0.0.1:1", APIToken: "test-token"}})
	_, err := service.ListDomains(context.Background(), "srv-1")
	if err == nil || !strings.Contains(err.Error(), "unreachable") {
		t.Fatalf("expected unreachable error, got %v", err)
	}
}

func TestMaildContractFallbackToLegacy(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/domains", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]models.MailDomain{{ID: "dom-legacy", Name: "legacy.com"}})
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	svc := NewMailService(&stubMailRepo{server: &models.MailServer{ID: "srv-legacy", MaildURL: server.URL, APIToken: "token"}})
	domains, err := svc.ListDomains(context.Background(), "srv-legacy")
	if err != nil {
		t.Fatalf("expected fallback to succeed, got %v", err)
	}
	if len(domains) != 1 || domains[0].Name != "legacy.com" {
		t.Fatalf("unexpected domains: %+v", domains)
	}
}
