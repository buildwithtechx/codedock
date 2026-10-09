package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"codedock/internal/models"
)

func stubServer(t *testing.T, handler http.HandlerFunc) (*Client, func()) {
	t.Helper()
	server := httptest.NewServer(handler)
	return NewClient(server.URL, "test-token"), server.Close
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func TestListServers(t *testing.T) {
	client, done := stubServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/servers" || r.Method != http.MethodGet {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"data": []models.Server{{ID: "srv-1", Name: "primary"}}})
	})
	defer done()

	servers, err := client.ListServers()
	if err != nil {
		t.Fatal(err)
	}
	if len(servers) != 1 || servers[0].ID != "srv-1" || servers[0].Name != "primary" {
		t.Fatalf("unexpected servers: %+v", servers)
	}
}

func TestCreateServer(t *testing.T) {
	var received models.CreateServerRequest
	client, done := stubServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/servers" || r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{"data": models.Server{ID: "srv-2", Name: received.Name}})
	})
	defer done()

	server, err := client.CreateServer(&models.CreateServerRequest{Name: "edge", IPAddress: "10.0.0.2"})
	if err != nil {
		t.Fatal(err)
	}
	if server.ID != "srv-2" || received.IPAddress != "10.0.0.2" {
		t.Fatalf("unexpected create roundtrip: %+v %+v", server, received)
	}
}

func TestDeleteServer(t *testing.T) {
	var deleted string
	client, done := stubServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/servers/srv-9" || r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		deleted = "srv-9"
		w.WriteHeader(http.StatusNoContent)
	})
	defer done()

	if err := client.DeleteServer("srv-9"); err != nil {
		t.Fatal(err)
	}
	if deleted != "srv-9" {
		t.Fatal("expected delete request")
	}
}

func TestServersSurfaceErrors(t *testing.T) {
	client, done := stubServer(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"message": "bad token"})
	})
	defer done()

	if _, err := client.ListServers(); err == nil {
		t.Error("expected list error")
	}
	if _, err := client.CreateServer(&models.CreateServerRequest{Name: "x"}); err == nil {
		t.Error("expected create error")
	}
	if err := client.DeleteServer("srv-1"); err == nil {
		t.Error("expected delete error")
	}
}
