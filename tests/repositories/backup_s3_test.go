package repositories_test

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"

	_ "modernc.org/sqlite"
)

type mockFailDecryptVault struct{}

func (mockFailDecryptVault) Encrypt(plaintext string) (string, error) {
	return "enc:" + plaintext, nil
}

func (mockFailDecryptVault) Decrypt(ciphertext string) (string, error) {
	if strings.HasPrefix(ciphertext, "corrupt:") {
		return "", errors.New("decryption error: invalid ciphertext")
	}
	return strings.TrimPrefix(ciphertext, "enc:"), nil
}

func TestS3DestinationRepoUpdateAndVaultDecryption(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:?_pragma=foreign_keys(ON)")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := repositories.RunMigrations(db); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	repo := repositories.NewS3DestinationRepo(db, mockFailDecryptVault{})

	dest := &models.S3Destination{
		Name:            "s3-test",
		Bucket:          "my-bucket",
		Endpoint:        "s3.us-east-1.amazonaws.com",
		AccessKeyID:     "ak",
		SecretAccessKey: "secret123",
	}

	if err := repo.CreateS3Destination(context.Background(), dest); err != nil {
		t.Fatalf("failed to create s3 destination: %v", err)
	}

	fetched, err := repo.GetS3Destination(context.Background(), dest.ID)
	if err != nil {
		t.Fatalf("failed to get s3 destination: %v", err)
	}
	if fetched.SecretAccessKey != "secret123" {
		t.Fatalf("expected decrypted secret 'secret123', got '%s'", fetched.SecretAccessKey)
	}

	dest.Name = "s3-test-updated"
	dest.SecretAccessKey = "newsecret"
	if err := repo.UpdateS3Destination(context.Background(), dest); err != nil {
		t.Fatalf("failed to update s3 destination: %v", err)
	}

	updated, err := repo.GetS3Destination(context.Background(), dest.ID)
	if err != nil {
		t.Fatalf("failed to get updated s3 destination: %v", err)
	}
	if updated.Name != "s3-test-updated" {
		t.Fatalf("expected updated name 's3-test-updated', got '%s'", updated.Name)
	}
	if updated.SecretAccessKey != "newsecret" {
		t.Fatalf("expected decrypted secret 'newsecret', got '%s'", updated.SecretAccessKey)
	}

	list, err := repo.ListS3Destinations(context.Background())
	if err != nil {
		t.Fatalf("failed to list s3 destinations: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 destination, got %d", len(list))
	}
	if list[0].Name != "s3-test-updated" {
		t.Fatalf("expected name 's3-test-updated', got '%s'", list[0].Name)
	}
	if list[0].SecretAccessKey != "newsecret" {
		t.Fatalf("expected decrypted secret 'newsecret', got '%s'", list[0].SecretAccessKey)
	}

	_, err = db.Exec("UPDATE s3_destinations SET secret_access_key = ? WHERE id = ?", "corrupt:data", dest.ID)
	if err != nil {
		t.Fatalf("failed to insert corrupt secret: %v", err)
	}

	_, err = repo.GetS3Destination(context.Background(), dest.ID)
	if err == nil {
		t.Fatalf("expected error on corrupt vault decryption in GetS3Destination")
	}

	_, err = repo.ListS3Destinations(context.Background())
	if err == nil {
		t.Fatalf("expected error on corrupt vault decryption in ListS3Destinations")
	}
}

func TestS3DestinationDefaultAndVerification(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:?_pragma=foreign_keys(ON)")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := repositories.RunMigrations(db); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	repo := repositories.NewS3DestinationRepo(db, mockFailDecryptVault{})

	d1 := &models.S3Destination{
		Name:            "dest-1",
		Bucket:          "b1",
		Endpoint:        "s3.us-east-1.amazonaws.com",
		PathPrefix:      "backups/db1",
		IsDefault:       true,
		SecretAccessKey: "key1",
	}
	d2 := &models.S3Destination{
		Name:            "dest-2",
		Bucket:          "b2",
		Endpoint:        "s3.us-east-1.amazonaws.com",
		PathPrefix:      "backups/db2",
		IsDefault:       false,
		SecretAccessKey: "key2",
	}

	if err := repo.CreateS3Destination(context.Background(), d1); err != nil {
		t.Fatalf("create d1: %v", err)
	}
	if err := repo.CreateS3Destination(context.Background(), d2); err != nil {
		t.Fatalf("create d2: %v", err)
	}

	fetched1, _ := repo.GetS3Destination(context.Background(), d1.ID)
	if !fetched1.IsDefault || fetched1.PathPrefix != "backups/db1" {
		t.Fatalf("d1 default or pathPrefix incorrect: %+v", fetched1)
	}

	if err := repo.SetDefaultDestination(context.Background(), d2.ID); err != nil {
		t.Fatalf("set default d2: %v", err)
	}

	fetched1, _ = repo.GetS3Destination(context.Background(), d1.ID)
	fetched2, _ := repo.GetS3Destination(context.Background(), d2.ID)
	if fetched1.IsDefault || !fetched2.IsDefault {
		t.Fatalf("expected d2 to be default and d1 not default")
	}

	if err := repo.RecordVerificationResult(context.Background(), d2.ID, true, ""); err != nil {
		t.Fatalf("record verified: %v", err)
	}
	fetched2, _ = repo.GetS3Destination(context.Background(), d2.ID)
	if fetched2.LastVerifiedAt == nil || fetched2.LastVerifyError != "" {
		t.Fatalf("expected last_verified_at set and error empty: %+v", fetched2)
	}

	if err := repo.RecordVerificationResult(context.Background(), d2.ID, false, "bucket not found"); err != nil {
		t.Fatalf("record verify fail: %v", err)
	}
	fetched2, _ = repo.GetS3Destination(context.Background(), d2.ID)
	if fetched2.LastVerifyError != "bucket not found" {
		t.Fatalf("expected error 'bucket not found', got '%s'", fetched2.LastVerifyError)
	}
}
