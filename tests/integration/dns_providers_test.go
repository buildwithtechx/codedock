package integration_test

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"

	"codedock.run/codedock/internal/models"
	dnsproviders "codedock.run/codedock/internal/services/system/dns_providers"
)

type providerSettingsRepo struct {
	settings *models.ServerSettings
}

func (r providerSettingsRepo) GetServerSettings(context.Context) (*models.ServerSettings, error) {
	return r.settings, nil
}

func (providerSettingsRepo) UpdateServerSettings(context.Context, *models.ServerSettings) error {
	return nil
}

func (providerSettingsRepo) ListProjects(context.Context) ([]map[string]any, error) {
	return nil, nil
}

type rewriteTransport struct {
	base string
}

func (t rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	base, err := url.Parse(t.base)
	if err != nil {
		return nil, err
	}
	clone := req.Clone(req.Context())
	clone.URL.Scheme = base.Scheme
	clone.URL.Host = base.Host
	return http.DefaultTransport.RoundTrip(clone)
}

func providerClient(serverURL string) func() *http.Client {
	return func() *http.Client {
		return &http.Client{Transport: rewriteTransport{base: serverURL}}
	}
}

func newIPv4TestServer(handler http.Handler) *httptest.Server {
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	server := &httptest.Server{Listener: listener, Config: &http.Server{Handler: handler}}
	server.Start()
	return server
}

func TestCloudflareProvisionAndDeprovision(t *testing.T) {
	deleted := false
	created := false
	var stateMu sync.Mutex
	server := newIPv4TestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/client/v4/zones":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"result":[{"id":"zone-id"}],"success":true}`))
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/dns_records"):
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"result":[{"id":"record-id","content":"192.0.2.10"}],"result_info":{"page":1,"total_pages":1}}`))
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/dns_records"):
			stateMu.Lock()
			created = true
			stateMu.Unlock()
			w.WriteHeader(http.StatusCreated)
		case r.Method == http.MethodDelete && strings.Contains(r.URL.Path, "/dns_records/"):
			stateMu.Lock()
			deleted = true
			stateMu.Unlock()
			w.WriteHeader(http.StatusOK)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	service := dnsproviders.NewWithHTTPClient(providerSettingsRepo{settings: &models.ServerSettings{CloudflareAPIToken: "token"}}, providerClient(server.URL))
	if _, err := service.ProvisionRecord(context.Background(), "app.example.com", "A", "192.0.2.10"); err == nil {
		t.Fatal("expected existing Cloudflare record conflict")
	}

	deprovisionSettings := providerSettingsRepo{settings: &models.ServerSettings{CloudflareAPIToken: "token"}}
	deprovisionService := dnsproviders.NewWithHTTPClient(deprovisionSettings, providerClient(server.URL))
	if err := deprovisionService.DeprovisionRecord(context.Background(), "app.example.com", "A", "192.0.2.10", "cloudflare"); err != nil {
		t.Fatalf("deprovision Cloudflare record: %v", err)
	}
	stateMu.Lock()
	createdResult, deletedResult := created, deleted
	stateMu.Unlock()
	if createdResult {
		t.Fatal("Cloudflare created a record despite an existing record")
	}
	if !deletedResult {
		t.Fatal("Cloudflare did not delete the matching record")
	}
}

func TestNamecheapProvisionAndDeprovision(t *testing.T) {
	setHostsCalls := 0
	retained := false
	var stateMu sync.Mutex
	server := newIPv4TestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		command := r.URL.Query().Get("Command")
		w.Header().Set("Content-Type", "application/xml")
		if command == "namecheap.domains.dns.getHosts" {
			_, _ = w.Write([]byte(`<ApiResponse Status="OK"><CommandResponse><HostRecords><host Name="app" Type="A" Address="192.0.2.10"/><host Name="www" Type="CNAME" Address="app.example.com"/></HostRecords></CommandResponse></ApiResponse>`))
			return
		}
		if command == "namecheap.domains.dns.setHosts" {
			stateMu.Lock()
			setHostsCalls++
			stateMu.Unlock()
			if r.URL.Query().Get("HostName1") == "app" {
				stateMu.Lock()
				retained = true
				stateMu.Unlock()
			}
			_, _ = w.Write([]byte(`<ApiResponse Status="OK"><CommandResponse/></ApiResponse>`))
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	settings := &models.ServerSettings{NamecheapAPIUser: "user", NamecheapAPIKey: "key", NamecheapClientIP: "127.0.0.1"}
	service := dnsproviders.NewWithHTTPClient(providerSettingsRepo{settings: settings}, providerClient(server.URL))
	if _, err := service.ProvisionRecord(context.Background(), "app.example.com", "A", "192.0.2.20"); err == nil {
		t.Fatal("expected Namecheap conflict for an existing host/type")
	}
	if err := service.DeprovisionRecord(context.Background(), "app.example.com", "A", "192.0.2.10", "namecheap"); err != nil {
		t.Fatalf("deprovision Namecheap record: %v", err)
	}
	stateMu.Lock()
	setHostsResult, retainedResult := setHostsCalls, retained
	stateMu.Unlock()
	if retainedResult {
		t.Fatal("Namecheap deprovision request retained the managed host")
	}
	if setHostsResult != 1 {
		t.Fatalf("expected one Namecheap setHosts call, got %d", setHostsResult)
	}
}

func TestSpaceshipProvisionAndDeprovision(t *testing.T) {
	putCalled := false
	deleteCalled := false
	recordPresent := false
	invalidPayload := false
	var stateMu sync.Mutex
	server := newIPv4TestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/dns/records/example.com" {
			http.NotFound(w, r)
			return
		}
		if r.Header.Get("X-API-Key") != "key" || r.Header.Get("X-API-Secret") != "secret" {
			http.Error(w, "missing credentials", http.StatusUnauthorized)
			return
		}
		switch r.Method {
		case http.MethodGet:
			stateMu.Lock()
			present := recordPresent
			stateMu.Unlock()
			if present {
				_, _ = w.Write([]byte(`{"items":[{"name":"app","type":"A","address":"192.0.2.30","ttl":3600}]}`))
			} else {
				_, _ = w.Write([]byte(`{"items":[]}`))
			}
		case http.MethodPut:
			var payload struct {
				Items []map[string]any `json:"items"`
			}
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || len(payload.Items) != 1 {
				stateMu.Lock()
				invalidPayload = true
				stateMu.Unlock()
				http.Error(w, "invalid payload", http.StatusBadRequest)
				return
			}
			stateMu.Lock()
			putCalled = true
			recordPresent = true
			stateMu.Unlock()
			w.WriteHeader(http.StatusNoContent)
		case http.MethodDelete:
			stateMu.Lock()
			deleteCalled = true
			recordPresent = false
			stateMu.Unlock()
			w.WriteHeader(http.StatusNoContent)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	settings := &models.ServerSettings{SpaceshipAPIKey: "key", SpaceshipAPISecret: "secret"}
	service := dnsproviders.NewWithHTTPClient(providerSettingsRepo{settings: settings}, providerClient(server.URL))
	if _, err := service.ProvisionRecord(context.Background(), "app.example.com", "A", "192.0.2.30"); err != nil {
		t.Fatalf("provision Spaceship record: %v", err)
	}
	if err := service.DeprovisionRecord(context.Background(), "app.example.com", "A", "192.0.2.30", "spaceship"); err != nil {
		t.Fatalf("deprovision Spaceship record: %v", err)
	}
	stateMu.Lock()
	putResult, deleteResult, invalidResult := putCalled, deleteCalled, invalidPayload
	stateMu.Unlock()
	if invalidResult {
		t.Fatal("Spaceship received an invalid record update payload")
	}
	if !putResult {
		t.Fatal("Spaceship did not receive the PUT record update")
	}
	if !deleteResult {
		t.Fatal("Spaceship did not receive the record deletion request")
	}
}
