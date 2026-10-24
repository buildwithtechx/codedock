package repositories

import (
	"context"
	"testing"

	"codedock/internal/models"
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
	list, err := repo.List(ctx, "")
	if err != nil {
		t.Fatalf("list databases: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 database, got %d", len(list))
	}
	byProject, err := repo.ListByProject(ctx, "project-one", "")
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

func TestDatabaseSearchFiltersByNameEngineAndDatabaseName(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	repo := NewDatabaseRepo(db, databaseTestVault{})
	for _, seed := range []models.Database{
		{ProjectID: "project-one", Name: "main-db", Engine: "postgres", Version: "16", Port: 5432, Username: "app", Password: "secret", DatabaseName: "appdb"},
		{ProjectID: "project-one", Name: "cache-db", Engine: "redis", Version: "7", Port: 6379, Username: "default", Password: "secret", DatabaseName: "cache"},
	} {
		seed := seed
		if err := repo.Create(ctx, &seed); err != nil {
			t.Fatalf("create database: %v", err)
		}
	}
	cases := []struct {
		name   string
		search string
		want   int
	}{
		{name: "name", search: "main-db", want: 1},
		{name: "engine case-insensitive", search: "POSTGRES", want: 1},
		{name: "database name", search: "appdb", want: 1},
		{name: "no match", search: "nothing-matches-this", want: 0},
		{name: "underscore literal", search: "main_db", want: 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			list, err := repo.List(ctx, tc.search)
			if err != nil {
				t.Fatalf("search databases: %v", err)
			}
			if len(list) != tc.want {
				t.Fatalf("expected %d matches for %q, got %d", tc.want, tc.search, len(list))
			}
			byProject, err := repo.ListByProject(ctx, "project-one", tc.search)
			if err != nil {
				t.Fatalf("search project databases: %v", err)
			}
			if len(byProject) != tc.want {
				t.Fatalf("expected %d matches for %q, got %d", tc.want, tc.search, len(byProject))
			}
		})
	}
}

func TestDatabaseSearchIDsByName(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	repo := NewDatabaseRepo(db, databaseTestVault{})
	seed := models.Database{ProjectID: "project-one", Name: "main-db", Engine: "postgres", Version: "16", Port: 5432, Username: "app", Password: "secret", DatabaseName: "appdb"}
	if err := repo.Create(ctx, &seed); err != nil {
		t.Fatalf("create database: %v", err)
	}
	ids, err := repo.SearchIDsByName(ctx, "main", 0)
	if err != nil {
		t.Fatalf("search ids: %v", err)
	}
	if len(ids) != 1 || ids[0] != seed.ID {
		t.Fatalf("expected only the seeded database id, got %v", ids)
	}
	ids, err = repo.SearchIDsByName(ctx, "nothing-matches-this", 0)
	if err != nil {
		t.Fatalf("search ids: %v", err)
	}
	if len(ids) != 0 {
		t.Fatalf("expected no ids, got %v", ids)
	}
}
