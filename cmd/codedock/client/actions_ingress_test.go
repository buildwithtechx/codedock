package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"codedock/internal/models"
)

func TestIngressClientOperations(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == "GET" && r.URL.Path == "/api/v1/dns/ingress":
			json.NewEncoder(w).Encode(map[string]any{
				"status": "success",
				"data": []*models.IngressConnection{
					{
						ID:       "ing-1",
						ServerID: "srv-1",
						Name:     "Test Tunnel",
						Status:   "ready",
					},
				},
			})
		case r.Method == "POST" && r.URL.Path == "/api/v1/dns/ingress":
			json.NewEncoder(w).Encode(map[string]any{
				"status": "success",
				"data": &models.IngressConnection{
					ID:       "ing-2",
					ServerID: "srv-2",
					Name:     "New Tunnel",
					Status:   "provisioning",
				},
			})
		case r.Method == "POST" && r.URL.Path == "/api/v1/dns/ingress/ing-1/retry":
			json.NewEncoder(w).Encode(map[string]any{
				"status": "success",
				"data": &models.IngressConnection{
					ID:     "ing-1",
					Name:   "Test Tunnel",
					Status: "reconciling",
				},
			})
		case r.Method == "DELETE" && r.URL.Path == "/api/v1/dns/ingress/ing-1":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]any{"status": "success"})
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	cli := NewClient(ts.URL, "token")

	// Test List
	list, err := cli.ListIngressConnections()
	if err != nil {
		t.Fatalf("ListIngressConnections failed: %v", err)
	}
	if len(list) != 1 || list[0].Name != "Test Tunnel" {
		t.Fatalf("unexpected list response: %+v", list)
	}

	// Test Connect
	conn, err := cli.ConnectIngress("srv-2", "New Tunnel", "managed")
	if err != nil {
		t.Fatalf("ConnectIngress failed: %v", err)
	}
	if conn.ID != "ing-2" {
		t.Fatalf("unexpected connect response: %+v", conn)
	}

	// Test Retry
	retried, err := cli.RetryIngress("ing-1")
	if err != nil {
		t.Fatalf("RetryIngress failed: %v", err)
	}
	if retried.Status != "reconciling" {
		t.Fatalf("unexpected retry status: %s", retried.Status)
	}

	// Test Disconnect
	if err := cli.DisconnectIngress("ing-1"); err != nil {
		t.Fatalf("DisconnectIngress failed: %v", err)
	}
}

func TestActionsClientOperations(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == "GET" && r.URL.Path == "/api/v1/actions/workflows":
			json.NewEncoder(w).Encode(map[string]any{
				"status": "success",
				"data": []*models.ActionWorkflow{
					{ID: "wf-1", Name: "CI Test", Trigger: "push", Enabled: true},
				},
			})
		case r.Method == "GET" && r.URL.Path == "/api/v1/actions/runs":
			json.NewEncoder(w).Encode(map[string]any{
				"status": "success",
				"data": []*models.ActionRun{
					{ID: "run-1", WorkflowID: "wf-1", Name: "CI Test", Number: 1, Status: "ready"},
				},
			})
		case r.Method == "POST" && r.URL.Path == "/api/v1/actions/runs/run-1/retry":
			json.NewEncoder(w).Encode(map[string]any{
				"status": "success",
				"data": &models.ActionRun{
					ID: "run-2", WorkflowID: "wf-1", Name: "CI Test", Number: 2, Status: "running",
				},
			})
		case r.Method == "POST" && r.URL.Path == "/api/v1/actions/runs/run-1/cancel":
			json.NewEncoder(w).Encode(map[string]any{
				"status": "success",
				"data": &models.ActionRun{
					ID: "run-1", WorkflowID: "wf-1", Name: "CI Test", Number: 1, Status: "cancelled",
				},
			})
		case r.Method == "GET" && r.URL.Path == "/api/v1/actions/runners":
			json.NewEncoder(w).Encode(map[string]any{
				"status": "success",
				"data": []*models.ActionRunner{
					{ID: "runner-1", Name: "Default Runner", Mode: "container", Status: "online", Enabled: true},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	cli := NewClient(ts.URL, "token")

	workflows, err := cli.ListWorkflows()
	if err != nil || len(workflows) != 1 {
		t.Fatalf("ListWorkflows failed: %v", err)
	}

	runs, err := cli.ListRuns()
	if err != nil || len(runs) != 1 {
		t.Fatalf("ListRuns failed: %v", err)
	}

	retry, err := cli.RetryRun("run-1")
	if err != nil || retry.Status != "running" {
		t.Fatalf("RetryRun failed: %v", err)
	}

	cancelled, err := cli.CancelRun("run-1")
	if err != nil || cancelled.Status != "cancelled" {
		t.Fatalf("CancelRun failed: %v", err)
	}

	runners, err := cli.ListRunners()
	if err != nil || len(runners) != 1 {
		t.Fatalf("ListRunners failed: %v", err)
	}
}
