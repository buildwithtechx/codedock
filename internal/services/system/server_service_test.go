package system

import (
	"context"
	"testing"

	cdkdb "codedock/db"
	"codedock/internal/config"
	"codedock/internal/models"
	"codedock/internal/repositories"
	"codedock/internal/testdb"
)

func TestListServersByUserIncludesControlPlaneForOwner(t *testing.T) {
	db := testdb.Open(t)
	if err := cdkdb.RunMigrations(db); err != nil {
		t.Fatalf("run migrations: %v", err)
	}
	userRepo := repositories.NewUserRepo(db)
	user := &models.User{ID: "owner", Email: "owner@example.com", Role: models.UserRoleOwner, IsActive: true}
	if err := userRepo.CreateUser(context.Background(), user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	service := NewServerService(repositories.NewServerRepository(db, nil), userRepo, nil)
	servers, err := service.ListServersByUser(context.Background(), user.ID, "")
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
			if err := cdkdb.RunMigrations(db); err != nil {
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
	if err := cdkdb.RunMigrations(db); err != nil {
		t.Fatal(err)
	}
	service := NewServerService(repositories.NewServerRepository(db, nil), repositories.NewUserRepo(db), nil)
	servers, err := service.ListServersByUser(context.Background(), "api-token-test", "")
	if err != nil {
		t.Fatal(err)
	}
	for _, server := range servers {
		if server.IsControlPlane {
			t.Fatal("API token received owner control plane")
		}
	}
}

func TestUpdateServerBuildCacheBudget(t *testing.T) {
	db := testdb.Open(t)
	if err := cdkdb.RunMigrations(db); err != nil {
		t.Fatal(err)
	}
	users := repositories.NewUserRepo(db)
	user := &models.User{ID: "owner-budget", Email: "budget@example.com", Role: models.UserRoleOwner, IsActive: true}
	if err := users.CreateUser(context.Background(), user); err != nil {
		t.Fatal(err)
	}
	serverRepo := repositories.NewServerRepository(db, nil)
	service := NewServerService(serverRepo, users, nil)

	remoteServer := &models.Server{
		ID:          "srv-remote-budget",
		UserID:      user.ID,
		Name:        "remote-worker",
		IPAddress:   "192.0.2.10",
		SSHHost:     "192.0.2.10",
		SSHPort:     22,
		SSHUser:     "root",
		SSHPassword: "secret-password",
		Status:      models.ServerStatusOnline,
	}
	if err := serverRepo.Create(context.Background(), remoteServer); err != nil {
		t.Fatal(err)
	}

	t.Run("updates build cache budget successfully", func(t *testing.T) {
		budget := int64(10 * 1024 * 1024 * 1024) // 10 GiB
		updated, err := service.UpdateServer(context.Background(), remoteServer.ID, user.ID, models.UpdateServerRequest{
			BuildCacheBudgetBytes: &budget,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if updated.BuildCacheBudgetBytes == nil || *updated.BuildCacheBudgetBytes != budget {
			t.Fatalf("expected budget %d, got %v", budget, updated.BuildCacheBudgetBytes)
		}
		if updated.Capabilities == nil || !updated.Capabilities.BuildCacheBudget {
			t.Fatal("expected BuildCacheBudget capability to be true")
		}
	})

	t.Run("clears budget when budget is zero", func(t *testing.T) {
		zero := int64(0)
		updated, err := service.UpdateServer(context.Background(), remoteServer.ID, user.ID, models.UpdateServerRequest{
			BuildCacheBudgetBytes: &zero,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if updated.BuildCacheBudgetBytes != nil {
			t.Fatalf("expected nil budget, got %v", updated.BuildCacheBudgetBytes)
		}
	})

	t.Run("rejects budget below 1 GiB", func(t *testing.T) {
		tooSmall := int64(512 * 1024 * 1024)
		_, err := service.UpdateServer(context.Background(), remoteServer.ID, user.ID, models.UpdateServerRequest{
			BuildCacheBudgetBytes: &tooSmall,
		})
		if err == nil {
			t.Fatal("expected error for budget < 1 GiB")
		}
	})

	t.Run("rejects budget above 64 TiB", func(t *testing.T) {
		tooLarge := int64(70368744177664 + 1)
		_, err := service.UpdateServer(context.Background(), remoteServer.ID, user.ID, models.UpdateServerRequest{
			BuildCacheBudgetBytes: &tooLarge,
		})
		if err == nil {
			t.Fatal("expected error for budget > 64 TiB")
		}
	})
}
