package testdb

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Open(t *testing.T) *sql.DB {
	t.Helper()
	dsn := os.Getenv("CODEDOCK_TEST_PG_URL")
	if dsn == "" {
		t.Skip("CODEDOCK_TEST_PG_URL is not set")
	}
	parsed, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("parse test database url: %v", err)
	}
	raw := make([]byte, 6)
	if _, err := rand.Read(raw); err != nil {
		t.Fatalf("random database suffix: %v", err)
	}
	name := fmt.Sprintf("rt%s", hex.EncodeToString(raw))
	parsed.Path = "/" + name
	target := parsed.String()
	parsed.Path = "/postgres"
	maintenanceDSN := parsed.String()
	maintenance, err := sql.Open("pgx", maintenanceDSN)
	if err != nil {
		t.Fatalf("open maintenance database: %v", err)
	}
	defer maintenance.Close()
	if _, err := maintenance.Exec(`CREATE DATABASE "` + name + `"`); err != nil {
		t.Fatalf("create isolated test database: %v", err)
	}
	t.Cleanup(func() {
		cleanup, err := sql.Open("pgx", maintenanceDSN)
		if err != nil {
			return
		}
		defer cleanup.Close()
		_, _ = cleanup.Exec(`DROP DATABASE IF EXISTS "` + name + `" WITH (FORCE)`)
	})
	db, err := sql.Open("pgx", target)
	if err != nil {
		t.Fatalf("open isolated test database: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}
