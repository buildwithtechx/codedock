package system

import (
	"context"
	"database/sql"
	"testing"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"

	_ "modernc.org/sqlite"
)

func TestListServersByUserIncludesControlPlaneForOwner(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:?_pragma=foreign_keys(ON)")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer db.Close()
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
