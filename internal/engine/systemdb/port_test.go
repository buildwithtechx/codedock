package systemdb

import (
	"errors"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestResolvePostgresPortPersistsChoice(t *testing.T) {
	dir := t.TempDir()
	first, err := resolvePostgresPort(dir, 5432)
	if err != nil {
		t.Fatalf("resolve port: %v", err)
	}
	second, err := resolvePostgresPort(dir, 5432)
	if err != nil {
		t.Fatalf("resolve port again: %v", err)
	}
	if first != second {
		t.Fatalf("expected stable port, got %d then %d", first, second)
	}
	raw, err := os.ReadFile(filepath.Join(dir, "postgres", "port"))
	if err != nil {
		t.Fatalf("read stored port: %v", err)
	}
	if stored, _ := strconv.Atoi(string(raw[:len(raw)-1])); stored != first {
		t.Fatalf("expected stored port %d, got %q", first, raw)
	}
}

func TestResolvePostgresPortSkipsOccupied(t *testing.T) {
	var ln net.Listener
	var port int
	for range 10 {
		candidate, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("occupy port: %v", err)
		}
		port = candidate.Addr().(*net.TCPAddr).Port
		if port <= 65000 {
			ln = candidate
			break
		}
		candidate.Close()
	}
	if ln == nil {
		t.Skip("no suitable ephemeral port")
	}
	defer ln.Close()
	got, err := resolvePostgresPort(t.TempDir(), port)
	if err != nil {
		t.Fatalf("resolve port: %v", err)
	}
	if got == port {
		t.Fatalf("expected occupied port %d to be skipped", port)
	}
}

func TestRefreshPostgresPortOverwritesStored(t *testing.T) {
	dir := t.TempDir()
	first, err := resolvePostgresPort(dir, 5432)
	if err != nil {
		t.Fatalf("resolve port: %v", err)
	}
	second, err := refreshPostgresPort(dir, first+1)
	if err != nil {
		t.Fatalf("refresh port: %v", err)
	}
	if second == first {
		t.Fatalf("expected refresh to move past %d", first)
	}
	third, err := resolvePostgresPort(dir, 5432)
	if err != nil {
		t.Fatalf("resolve port: %v", err)
	}
	if third != second {
		t.Fatalf("expected stored port %d, got %d", second, third)
	}
}

func TestIsPortConflict(t *testing.T) {
	cases := map[string]bool{
		"failed to start postgres container: Error response from daemon: ports are not available: exposing port TCP 127.0.0.1:5432": true,
		"Bind for 127.0.0.1:5432 failed: port is already allocated":                                                                 true,
		"listen tcp4 127.0.0.1:5432: bind: address already in use":                                                                  true,
		"connection refused": false,
		"no such container":  false,
	}
	for message, want := range cases {
		if got := isPortConflict(errors.New(message)); got != want {
			t.Errorf("isPortConflict(%q) = %v, want %v", message, got, want)
		}
	}
}
