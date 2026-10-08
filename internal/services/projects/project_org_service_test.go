package projects

import (
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	"codedock.run/codedock/internal/testdb"
	"context"
	"sync"
	"testing"
)

func TestConcurrentDefaultOrganizationProvisioning(t *testing.T) {
	db := testdb.Open(t)
	if err := repositories.RunMigrations(db); err != nil {
		t.Fatal(err)
	}
	if err := repositories.NewUserRepo(db).CreateUser(context.Background(), &models.User{ID: "user", Email: "user@example.com", Role: models.UserRoleMember, IsActive: true}); err != nil {
		t.Fatal(err)
	}
	service := NewProjectService(nil, nil, nil, nil, nil, repositories.NewOrganizationRepository(db))
	results := make(chan string, 12)
	failures := make(chan error, 12)
	var group sync.WaitGroup
	for range 12 {
		group.Add(1)
		go func() {
			defer group.Done()
			org, err := service.GetOrCreateDefaultOrganization(context.Background(), "user")
			if err != nil {
				failures <- err
				return
			}
			results <- org.ID
		}()
	}
	group.Wait()
	close(results)
	close(failures)
	for err := range failures {
		t.Fatal(err)
	}
	ids := make(map[string]bool)
	for id := range results {
		ids[id] = true
	}
	if len(ids) != 1 {
		t.Fatalf("created multiple default organizations: %v", ids)
	}
	organizations, err := service.ListOrganizationsByUser(context.Background(), "user")
	if err != nil || len(organizations) != 1 {
		t.Fatalf("unexpected organizations: %v %v", organizations, err)
	}
}
