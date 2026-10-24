package repositories

import (
	"context"
	"errors"
	"testing"
)

func TestControlPlaneSnapshotRequiresDumper(t *testing.T) {
	db := openPGTestDB(t)
	_, err := NewBackupRepo(db, nil).ControlPlaneSnapshot(context.Background())
	if err == nil {
		t.Fatal("expected error when snapshotter is not configured")
	}
}

func TestControlPlaneSnapshotUsesDumper(t *testing.T) {
	db := openPGTestDB(t)
	repo := NewBackupRepo(db, nil)
	repo.SetSnapshotDumper(func(ctx context.Context) ([]byte, error) {
		return []byte("dump-bytes"), nil
	})
	data, err := repo.ControlPlaneSnapshot(context.Background())
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	if string(data) != "dump-bytes" {
		t.Fatalf("unexpected snapshot payload %q", data)
	}
	repo.SetSnapshotDumper(func(ctx context.Context) ([]byte, error) {
		return nil, errors.New("dump failed")
	})
	if _, err := repo.ControlPlaneSnapshot(context.Background()); err == nil {
		t.Fatal("expected dumper error to propagate")
	}
}
