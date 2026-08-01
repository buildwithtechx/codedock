package repositories_test

import (
	"context"
	"database/sql"
	"testing"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"

	_ "modernc.org/sqlite"
)

func TestOrganizationRepositoryOnFreshDatabase(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:?_pragma=foreign_keys(ON)")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := repositories.RunMigrations(db); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	repo := repositories.NewOrganizationRepository(db)
	org := &models.Organization{ID: "org-1", Name: "Test Organization"}
	if err := repo.Create(context.Background(), org); err != nil {
		t.Fatalf("create organization: %v", err)
	}
	loaded, err := repo.GetByID(context.Background(), org.ID)
	if err != nil || loaded == nil {
		t.Fatalf("get organization: org=%v err=%v", loaded, err)
	}
	if _, err := db.Exec(`INSERT INTO users (id, email, password_hash) VALUES (?, ?, ?)`, "user-1", "member@example.com", "test-hash"); err != nil {
		t.Fatalf("create member user: %v", err)
	}

	member := &models.OrganizationMember{
		ID:             "member-1",
		OrganizationID: org.ID,
		UserID:         "user-1",
		Email:          "member@example.com",
		Permission:     models.MemberPermissionMember,
		Status:         models.MemberStatusAccepted,
	}
	if err := repo.AddMember(context.Background(), member); err != nil {
		t.Fatalf("add organization member: %v", err)
	}
	members, err := repo.ListMembers(context.Background(), org.ID)
	if err != nil || len(members) != 1 {
		t.Fatalf("list organization members: count=%d err=%v", len(members), err)
	}
}
