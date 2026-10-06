package repositories

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
)

func TestControlPlaneSnapshotIncludesCommittedWALWrites(t *testing.T) {
	path := filepath.Join(t.TempDir(), "source.db")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	for _, query := range []string{"PRAGMA journal_mode=WAL", "PRAGMA wal_autocheckpoint=0", "CREATE TABLE sample(value TEXT)", "INSERT INTO sample VALUES('committed')"} {
		if _, err := db.Exec(query); err != nil {
			t.Fatal(err)
		}
	}
	data, err := NewBackupRepo(db, nil).ControlPlaneSnapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	snapshot := filepath.Join(t.TempDir(), "snapshot.db")
	if err := os.WriteFile(snapshot, data, 0600); err != nil {
		t.Fatal(err)
	}
	restored, err := sql.Open("sqlite", snapshot)
	if err != nil {
		t.Fatal(err)
	}
	defer restored.Close()
	var value string
	if err := restored.QueryRow("SELECT value FROM sample").Scan(&value); err != nil || value != "committed" {
		t.Fatal("snapshot lost WAL write", err)
	}
	var integrity string
	if err := restored.QueryRow("PRAGMA quick_check").Scan(&integrity); err != nil || integrity != "ok" {
		t.Fatal("invalid snapshot", err)
	}
}
