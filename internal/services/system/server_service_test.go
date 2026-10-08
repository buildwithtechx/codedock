package system

import (
	"context"
	"testing"

	"codedock.run/codedock/internal/config"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	"codedock.run/codedock/internal/testdb"
)

func TestListServersByUserIncludesControlPlaneForOwner(t *testing.T) {
	db := testdb.Open(t)
	if err := repositories.RunMigrations(db); err != nil {
		t.Fatalf("run migrations: %v", err)
	}
	userRepo := repositories.NewUserRepo(db)
	user := &models.User{ID: "owner", Email: "owner@example.com", Role: models.UserRoleOwner, IsActive: true}
	if err := userRepo.CreateUser(context.Background(), user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	service := NewServerService(repositories.NewServerRepository(db, nil), userRepo, nil)
	servers, err := service.ListServersByUser(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("list servers: %v", err)
	}
	if len(servers) != 1 || !servers[0].IsControlPlane || servers[0].ID != controlPlaneServerID {
		t.Fatalf("expected control plane server, got %#v", servers)
	}
}

func TestWorkerLimitsOnlyApplyInCloudMode(t *testing.T) {
	cfg := config.Get()
	previous := cfg.Cloud.Enabled
	t.Cleanup(func() { cfg.Cloud.Enabled = previous })
	for _, check := range []struct {
		name    string
		cloud   bool
		plan    string
		allowed bool
	}{
		{"self-hosted", false, "free", true},
		{"cloud hobby", true, "free", false},
		{"cloud pro", true, "pro", true},
	} {
		t.Run(check.name, func(t *testing.T) {
			cfg.Cloud.Enabled = check.cloud
			db := testdb.Open(t)
			if err := repositories.RunMigrations(db); err != nil {
				t.Fatal(err)
			}
			users := repositories.NewUserRepo(db)
			user := &models.User{ID: "owner", Email: "owner@example.com", Role: models.UserRoleOwner, IsActive: true, PlanType: check.plan}
			if err := users.CreateUser(context.Background(), user); err != nil {
				t.Fatal(err)
			}
			service := NewServerService(repositories.NewServerRepository(db, nil), users, nil)
			request := models.CreateServerRequest{Name: "worker", SSHHost: "192.0.2.1", SSHPassword: "test-password"}
			if _, err := service.CreateServer(context.Background(), user.ID, request); err != nil {
				t.Fatal(err)
			}
			_, err := service.CreateServer(context.Background(), user.ID, request)
			if check.allowed && err != nil {
				t.Fatalf("worker creation should be allowed: %v", err)
			}
			if !check.allowed && err == nil {
				t.Fatal("cloud hobby worker limit was bypassed")
			}
		})
	}
}

func TestListServersByAPITokenDoesNotLoadSyntheticUser(t *testing.T) {
	db := testdb.Open(t)
	if err := repositories.RunMigrations(db); err != nil {
		t.Fatal(err)
	}
	service := NewServerService(repositories.NewServerRepository(db, nil), repositories.NewUserRepo(db), nil)
	servers, err := service.ListServersByUser(context.Background(), "api-token-test")
	if err != nil {
		t.Fatal(err)
	}
	for _, server := range servers {
		if server.IsControlPlane {
			t.Fatal("API token received owner control plane")
		}
	}
}
