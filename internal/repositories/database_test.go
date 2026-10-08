package repositories

import (
	"context"
	"testing"

	"codedock.run/codedock/internal/models"
)

type databaseTestVault struct{}

func (databaseTestVault) Encrypt(plaintext string) (string, error) { return "enc:" + plaintext, nil }
func (databaseTestVault) Decrypt(ciphertext string) (string, error) {
	if len(ciphertext) > 4 && ciphertext[:4] == "enc:" {
		return ciphertext[4:], nil
	}
	return ciphertext, nil
}

func TestDatabaseRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	repo := NewDatabaseRepo(db, databaseTestVault{})
	database := &models.Database{
		ProjectID:          "project-one",
		Name:               "main-db",
		Engine:             "postgres",
		Version:            "16",
		Port:               5432,
		Username:           "app",
		Password:           "secret",
		DatabaseName:       "appdb",
		Status:             "running",
		LogicalReplication: true,
		CPULimit:           1.5,
		MemoryLimit:        1024,
	}
	if err := repo.Create(ctx, database); err != nil {
		t.Fatalf("create database: %v", err)
	}
	stored, err := repo.GetByID(ctx, database.ID)
	if err != nil {
		t.Fatalf("get database: %v", err)
	}
	if stored.Name != "main-db" || stored.Password != "secret" || !stored.LogicalReplication || stored.CPULimit != 1.5 || stored.MemoryLimit != 1024 {
		t.Fatalf("unexpected database row: %+v", stored)
	}
	list, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("list databases: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 database, got %d", len(list))
	}
	byProject, err := repo.ListByProject(ctx, "project-one")
	if err != nil {
		t.Fatalf("list by project: %v", err)
	}
	if len(byProject) != 1 {
		t.Fatalf("expected 1 project database, got %d", len(byProject))
	}
	database.Name = "renamed-db"
	database.Password = "rotated"
	if err := repo.Update(ctx, database); err != nil {
		t.Fatalf("update database: %v", err)
	}
	updated, err := repo.GetByID(ctx, database.ID)
	if err != nil {
		t.Fatalf("get updated database: %v", err)
	}
	if updated.Name != "renamed-db" || updated.Password != "rotated" {
		t.Fatalf("update did not persist: %+v", updated)
	}
	if err := repo.Delete(ctx, database.ID); err != nil {
		t.Fatalf("delete database: %v", err)
	}
	if _, err := repo.GetByID(ctx, database.ID); err == nil {
		t.Fatal("expected deleted database to be gone")
	}
}
