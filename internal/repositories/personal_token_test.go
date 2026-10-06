package repositories

import (
	"codedock.run/codedock/internal/models"
	"context"
	"testing"
)

func TestPersonalTokenLookupAndRevocation(t *testing.T) {
	db := openTestDB(t)
	db.SetMaxOpenConns(1)
	if err := RunMigrations(db); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO users(id,email,name,password_hash) VALUES('user','user@example.com','User','hash')`); err != nil {
		t.Fatal(err)
	}
	repo := NewUserRepo(db)
	ctx := context.Background()
	pat := &models.PersonalAccessToken{ID: "pat", UserID: "user", Name: "automation", TokenHash: "digest", Prefix: "vpt_1234", AccessLevel: "read", ProjectScope: "all"}
	if err := repo.CreatePAT(ctx, pat); err != nil {
		t.Fatal(err)
	}
	found, err := repo.GetPATByHash(ctx, "digest")
	if err != nil || found.UserID != "user" || found.ExpiresAt != nil {
		t.Fatalf("lookup failed: %v", err)
	}
	if err := repo.DeletePAT(ctx, "pat", "user"); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.GetPATByHash(ctx, "digest"); err == nil {
		t.Fatal("revoked token remained usable")
	}
}
