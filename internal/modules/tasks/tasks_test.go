package tasks

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

type mockTaskService struct {
	tasks []ScheduledTask
}

func (m *mockTaskService) ListByService(ctx context.Context, serviceID string) ([]ScheduledTask, error) {
	return m.tasks, nil
}
func (m *mockTaskService) ListByProject(ctx context.Context, projectID string) ([]ScheduledTask, error) {
	return m.tasks, nil
}
func (m *mockTaskService) CreateScheduledTask(ctx context.Context, task *ScheduledTask) (*ScheduledTask, error) {
	task.ID = "task-1"
	m.tasks = append(m.tasks, *task)
	return task, nil
}
func (m *mockTaskService) GetScheduledTask(ctx context.Context, id string) (*ScheduledTask, error) {
	for _, t := range m.tasks {
		if t.ID == id {
			return &t, nil
		}
	}
	return nil, nil
}
func (m *mockTaskService) UpdateScheduledTask(ctx context.Context, task *ScheduledTask) error {
	return nil
}
func (m *mockTaskService) DeleteScheduledTask(ctx context.Context, id string) error {
	return nil
}
func (m *mockTaskService) ExecuteScheduledTask(ctx context.Context, id string) (string, error) {
	return "done", nil
}

type mockAccessChecker struct{}

type testTaskClaims struct {
	UserID string
	Role   string
}

func (m *mockAccessChecker) HasServiceAccess(ctx context.Context, userID, role, serviceID string) bool {
	return true
}
func (m *mockAccessChecker) HasProjectAccess(ctx context.Context, userID, role, projectID string) bool {
	return true
}

func TestScheduledTaskHandler_CRUD(t *testing.T) {
	svc := &mockTaskService{
		tasks: []ScheduledTask{
			{ID: "task-existing", ServiceID: "svc-1", Name: "Backup Job", Schedule: "0 0 * * *", Status: "active"},
		},
	}
	h := NewScheduledTaskHandler(svc, &mockAccessChecker{})
	e := echo.New()

	// List
	req := httptest.NewRequest(http.MethodGet, "/scheduled-tasks?serviceId=svc-1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &testTaskClaims{UserID: "usr-1", Role: "admin"})

	if err := h.ListProjectScheduledTasks(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}

	// Create
	body, _ := json.Marshal(ScheduledTask{
		ServiceID: "svc-1",
		Name:      "New Job",
		Schedule:  "*/5 * * * *",
		Command:   "echo 1",
	})
	req = httptest.NewRequest(http.MethodPost, "/scheduled-tasks", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.Set("user", &testTaskClaims{UserID: "usr-1", Role: "admin"})

	if err := h.Create(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", rec.Code)
	}

	// Get
	req = httptest.NewRequest(http.MethodGet, "/scheduled-tasks/task-existing", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("task-existing")
	c.Set("user", &testTaskClaims{UserID: "usr-1", Role: "admin"})

	if err := h.Get(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}

	// Run
	req = httptest.NewRequest(http.MethodPost, "/scheduled-tasks/task-existing/trigger", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("task-existing")
	c.Set("user", &testTaskClaims{UserID: "usr-1", Role: "admin"})

	if err := h.Run(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}
}

func TestScheduledTaskStatusConstants(t *testing.T) {
	if ScheduledTaskStatusActive != "active" {
		t.Fatalf("expected active, got %s", ScheduledTaskStatusActive)
	}
	if ScheduledTaskStatusCompleted != "completed" {
		t.Fatalf("expected completed, got %s", ScheduledTaskStatusCompleted)
	}
}

func TestScheduledTaskModel(t *testing.T) {
	task := ScheduledTask{
		ID:        "t1",
		ServiceID: "s1",
		Name:      "test",
		CreatedAt: time.Now(),
	}
	if task.ID != "t1" {
		t.Fatalf("expected t1, got %s", task.ID)
	}
}
