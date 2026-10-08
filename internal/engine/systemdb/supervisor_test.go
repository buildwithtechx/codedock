package systemdb

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestPasswordRoundtrip(t *testing.T) {
	dir := t.TempDir()
	first, err := loadOrGeneratePassword(dir)
	if err != nil {
		t.Fatalf("generate password: %v", err)
	}
	if len(first) != 64 {
		t.Fatalf("expected 64 hex chars, got %d", len(first))
	}
	second, err := loadOrGeneratePassword(dir)
	if err != nil {
		t.Fatalf("reload password: %v", err)
	}
	if first != second {
		t.Fatal("expected stable password across loads")
	}
	if runtime.GOOS == "windows" {
		return
	}
	info, err := os.Stat(filepath.Join(dir, "postgres", "password"))
	if err != nil {
		t.Fatalf("stat password file: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("expected 0600 password file, got %o", info.Mode().Perm())
	}
}

func TestPasswordRegeneratesWhenShort(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "postgres"), 0o700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "postgres", "password"), []byte("short\n"), 0o600); err != nil {
		t.Fatalf("write short password: %v", err)
	}
	password, err := loadOrGeneratePassword(dir)
	if err != nil {
		t.Fatalf("load password: %v", err)
	}
	if len(password) != 64 {
		t.Fatalf("expected regenerated 64-char password, got %d chars", len(password))
	}
}

func TestDatabaseURLFormat(t *testing.T) {
	supervisor := NewSupervisor(nil, ContainerName, t.TempDir(), "postgres:16", 5433)
	url := supervisor.databaseURL("secret")
	want := "postgres://codedock:secret@127.0.0.1:5433/codedock?sslmode=disable"
	if url != want {
		t.Fatalf("expected %q, got %q", want, url)
	}
}

func TestDatabaseURLHonorsConnectHost(t *testing.T) {
	supervisor := NewSupervisor(nil, ContainerName, t.TempDir(), "postgres:16", 5432)
	supervisor.SetConnectHost("codedock-postgres")
	url := supervisor.databaseURL("secret")
	want := "postgres://codedock:secret@codedock-postgres:5432/codedock?sslmode=disable"
	if url != want {
		t.Fatalf("expected %q, got %q", want, url)
	}
	supervisor.SetConnectHost("")
	if supervisor.databaseURL("secret") != want {
		t.Fatal("empty host override cleared the configured host")
	}
}
