package e2e_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"testing"

	"github.com/docker/docker/client"
	_ "modernc.org/sqlite"

	codedockhttp "codedock.run/codedock/internal/http"
	"codedock.run/codedock/internal/repositories"
	"codedock.run/codedock/internal/utils"
)

type e2eHarness struct {
	server    *httptest.Server
	db        *sql.DB
	vault     *utils.Vault
	dataDir   string
	client    *http.Client
	authToken string
}

func newE2EHarness(t *testing.T) *e2eHarness {
	dataDir := t.TempDir()
	vlt, err := utils.NewVault(dataDir)
	if err != nil {
		t.Fatalf("failed to create vault: %v", err)
	}

	dbPath := filepath.Join(dataDir, "codedock_e2e.db")
	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)")
	if err != nil {
		t.Fatalf("failed to open sqlite database: %v", err)
	}

	if err := repositories.RunMigrations(db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	jwtSecret := "super-secure-e2e-test-jwt-secret-32-chars-minimum!"
	t.Setenv("CODEDOCK_JWT_SECRET", jwtSecret)

	dockerClient, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	s, err := codedockhttp.NewServer(db, vlt, nil, nil, dockerClient, dataDir)
	if err != nil {
		t.Fatalf("failed to create http server: %v", err)
	}

	ts := httptest.NewTLSServer(s)
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("failed to create cookie jar: %v", err)
	}

	hc := ts.Client()
	hc.Jar = jar

	return &e2eHarness{
		server:  ts,
		db:      db,
		vault:   vlt,
		dataDir: dataDir,
		client:  hc,
	}
}

func (h *e2eHarness) Close() {
	if h.server != nil {
		h.server.Close()
	}
	if h.db != nil {
		_ = h.db.Close()
	}
}

func (h *e2eHarness) post(urlPath string, payload any, headers map[string]string) (*http.Response, map[string]any, error) {
	return h.doRequest("POST", urlPath, payload, headers)
}

func (h *e2eHarness) get(urlPath string, headers map[string]string) (*http.Response, map[string]any, error) {
	return h.doRequest("GET", urlPath, nil, headers)
}

func (h *e2eHarness) put(urlPath string, payload any, headers map[string]string) (*http.Response, map[string]any, error) {
	return h.doRequest("PUT", urlPath, payload, headers)
}

func (h *e2eHarness) delete(urlPath string, headers map[string]string) (*http.Response, map[string]any, error) {
	return h.doRequest("DELETE", urlPath, nil, headers)
}

func (h *e2eHarness) getCSRFToken(targetURL *url.URL) string {
	if targetURL == nil {
		targetURL, _ = url.Parse(h.server.URL + "/api/auth/me")
	}
	if h.client.Jar != nil {
		for _, cookie := range h.client.Jar.Cookies(targetURL) {
			if cookie.Name == "csrf_token" || cookie.Name == "_csrf" {
				return cookie.Value
			}
		}
	}
	req, err := http.NewRequest("GET", h.server.URL+"/api/auth/me", nil)
	if err == nil {
		if h.authToken != "" {
			req.Header.Set("Authorization", "Bearer "+h.authToken)
		}
		res, err := h.client.Do(req)
		if err == nil {
			_ = res.Body.Close()
		}
	}
	if h.client.Jar != nil {
		for _, cookie := range h.client.Jar.Cookies(targetURL) {
			if cookie.Name == "csrf_token" || cookie.Name == "_csrf" {
				return cookie.Value
			}
		}
	}
	return ""
}

func (h *e2eHarness) doRequest(method, urlPath string, payload any, headers map[string]string) (*http.Response, map[string]any, error) {
	var bodyReader io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal json payload: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, h.server.URL+urlPath, bodyReader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if h.authToken != "" && req.Header.Get("Authorization") == "" {
		req.Header.Set("Authorization", "Bearer "+h.authToken)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	if method == "POST" || method == "PUT" || method == "DELETE" {
		if req.Header.Get("X-CSRF-Token") == "" {
			if csrfVal := h.getCSRFToken(req.URL); csrfVal != "" {
				req.Header.Set("X-CSRF-Token", csrfVal)
			}
		}
	}

	res, err := h.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer res.Body.Close()

	respBody, _ := io.ReadAll(res.Body)
	var respMap map[string]any
	if len(respBody) > 0 {
		_ = json.Unmarshal(respBody, &respMap)
	}

	if respMap != nil && (urlPath == "/api/auth/signup" || urlPath == "/api/auth/signin") {
		if data, ok := respMap["data"].(map[string]any); ok {
			if token, ok := data["token"].(string); ok && token != "" {
				h.authToken = token
			}
		}
	}

	return res, respMap, nil
}
