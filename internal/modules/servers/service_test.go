package servers

import (
	"context"
	"testing"
)

type mockRepo struct {
	servers map[string]*Server
}

func newMockRepo() *mockRepo {
	return &mockRepo{servers: make(map[string]*Server)}
}

func (m *mockRepo) Create(ctx context.Context, s *Server) error {
	m.servers[s.ID] = s
	return nil
}

func (m *mockRepo) Update(ctx context.Context, s *Server) error {
	m.servers[s.ID] = s
	return nil
}

func (m *mockRepo) GetByID(ctx context.Context, id string) (*Server, error) {
	return m.servers[id], nil
}

func (m *mockRepo) GetByToken(ctx context.Context, token string) (*Server, error) {
	for _, s := range m.servers {
		if s.WorkerToken == token {
			return s, nil
		}
	}
	return nil, nil
}

func (m *mockRepo) ListByUser(ctx context.Context, userID, search string) ([]*Server, error) {
	var list []*Server
	for _, s := range m.servers {
		if s.UserID == userID || (s.UserID == "system" && s.IsLocal) {
			list = append(list, s)
		}
	}
	return list, nil
}

func (m *mockRepo) SearchIDsByName(ctx context.Context, term string, limit int) ([]string, error) {
	var ids []string
	for id := range m.servers {
		ids = append(ids, id)
	}
	return ids, nil
}

func (m *mockRepo) UpdateStatus(ctx context.Context, id string, status ServerStatus) error {
	if s, ok := m.servers[id]; ok {
		s.Status = status
	}
	return nil
}

func (m *mockRepo) UpdateMetrics(ctx context.Context, id string, metrics []byte) error {
	if s, ok := m.servers[id]; ok {
		s.Metrics = metrics
	}
	return nil
}

func (m *mockRepo) Delete(ctx context.Context, id string) error {
	delete(m.servers, id)
	return nil
}

type mockUsers struct {
	users map[string]*UserInfo
}

func (mu *mockUsers) GetUserByID(ctx context.Context, id string) (*UserInfo, error) {
	return mu.users[id], nil
}

func TestServerServiceCreateAndGet(t *testing.T) {
	repo := newMockRepo()
	users := &mockUsers{users: map[string]*UserInfo{
		"user-1": {ID: "user-1", Role: "owner", PlanType: "pro"},
	}}
	svc := NewService(repo, users, nil)

	s, err := svc.CreateServer(context.Background(), "user-1", CreateServerRequest{
		Name:        "Test Worker",
		SSHHost:     "192.168.1.10",
		SSHPassword: "secret-password",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Name != "Test Worker" {
		t.Errorf("expected name Test Worker, got %s", s.Name)
	}

	got, err := svc.GetServer(context.Background(), s.ID, "user-1")
	if err != nil {
		t.Fatalf("unexpected get error: %v", err)
	}
	if got.ID != s.ID {
		t.Errorf("expected ID %s, got %s", s.ID, got.ID)
	}
}

func TestServerServiceListControlPlaneForOwner(t *testing.T) {
	repo := newMockRepo()
	users := &mockUsers{users: map[string]*UserInfo{
		"owner-1": {ID: "owner-1", Role: "owner", PlanType: "pro"},
	}}
	svc := NewService(repo, users, nil)

	list, err := svc.ListServersByUser(context.Background(), "owner-1", "")
	if err != nil {
		t.Fatalf("unexpected list error: %v", err)
	}
	if len(list) != 1 || !list[0].IsControlPlane {
		t.Errorf("expected 1 control plane server, got %d", len(list))
	}
}
