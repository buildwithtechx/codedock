package takeover

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

type mockVault struct{}

func (m *mockVault) Encrypt(s string) (string, error) { return "enc:" + s, nil }
func (m *mockVault) Decrypt(s string) (string, error) {
	if strings.HasPrefix(s, "enc:") {
		return strings.TrimPrefix(s, "enc:"), nil
	}
	return s, nil
}

type mockRepo struct {
	runs map[string]*TakeoverRun
}

func newMockRepo() *mockRepo {
	return &mockRepo{runs: make(map[string]*TakeoverRun)}
}

func (m *mockRepo) Create(ctx context.Context, run *TakeoverRun) error {
	run.ID = "run-123"
	m.runs[run.ID] = run
	return nil
}

func (m *mockRepo) GetByID(ctx context.Context, id string) (*TakeoverRun, error) {
	run, ok := m.runs[id]
	if !ok {
		return nil, nil
	}
	return run, nil
}

func (m *mockRepo) UpdateStatus(ctx context.Context, id string, status TakeoverStatus, errMsg string) error {
	if run, ok := m.runs[id]; ok {
		run.Status = status
		run.Error = errMsg
	}
	return nil
}

func (m *mockRepo) UpdateDiscovered(ctx context.Context, id string, discoveredJSON string) error {
	if run, ok := m.runs[id]; ok {
		run.DiscoveredJSON = discoveredJSON
	}
	return nil
}

func (m *mockRepo) UpdateAdopted(ctx context.Context, id string, projectIDs []string) error {
	if run, ok := m.runs[id]; ok {
		run.AdoptedProjectIDs = strings.Join(projectIDs, ",")
	}
	return nil
}

func (m *mockRepo) ListByUser(ctx context.Context, userID string) ([]*TakeoverRun, error) {
	var list []*TakeoverRun
	for _, r := range m.runs {
		if r.UserID == userID {
			list = append(list, r)
		}
	}
	return list, nil
}

type mockProjectCreator struct {
	createdProjects []string
	createdApps     []string
}

func (m *mockProjectCreator) CreateProject(ctx context.Context, id, name, description string) error {
	m.createdProjects = append(m.createdProjects, id)
	return nil
}

func (m *mockProjectCreator) CreateApp(ctx context.Context, id, projectID, name, imageRef, runtimeMode string, port int, status string) error {
	m.createdApps = append(m.createdApps, id)
	return nil
}

type testClaims struct {
	UserID string
}

func TestSerializeDeserializeStack(t *testing.T) {
	stack := &DiscoveredStack{
		Host:     "127.0.0.1",
		Platform: TakeoverPlatformDocker,
		Containers: []DiscoveredContainer{
			{Name: "web", Image: "nginx:latest", Ports: []string{"80/tcp"}},
		},
	}

	raw, err := SerializeStack(stack)
	if err != nil {
		t.Fatalf("SerializeStack failed: %v", err)
	}

	deserialized, err := DeserializeStack(raw)
	if err != nil {
		t.Fatalf("DeserializeStack failed: %v", err)
	}

	if deserialized.Host != stack.Host {
		t.Errorf("expected host %s, got %s", stack.Host, deserialized.Host)
	}
	if len(deserialized.Containers) != 1 {
		t.Fatalf("expected 1 container, got %d", len(deserialized.Containers))
	}
	if deserialized.Containers[0].Name != "web" {
		t.Errorf("expected container name 'web', got %s", deserialized.Containers[0].Name)
	}
}

func TestTakeoverHandler_Validation(t *testing.T) {
	e := echo.New()
	repo := newMockRepo()
	creator := &mockProjectCreator{}
	svc := NewTakeoverService(repo, creator)
	h := NewTakeoverHandler(svc, repo)

	t.Run("Scan unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/system/takeover/scan", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err := h.Scan(c)
		if err != nil {
			t.Fatal(err)
		}
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", rec.Code)
		}
	})

	t.Run("Scan invalid payload", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/system/takeover/scan", strings.NewReader(`{invalid}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", &testClaims{UserID: "user-1"})
		err := h.Scan(c)
		if err != nil {
			t.Fatal(err)
		}
		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("Scan missing required fields", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/system/takeover/scan", strings.NewReader(`{"host":""}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", &testClaims{UserID: "user-1"})
		err := h.Scan(c)
		if err != nil {
			t.Fatal(err)
		}
		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})
}
