package backup

import (
	"testing"

	"codedock/internal/models"
)

func TestBackupReconcileRegistersActiveConfigs(t *testing.T) {
	store := newMockStore()
	store.configs["cfg-one"] = &models.BackupConfig{ID: "cfg-one", Name: "one", Status: models.BackupConfigStatusActive, BackupEnabled: true, Schedule: "0 2 * * *"}
	manager := NewBackupManager(nil, store, t.TempDir())
	if err := manager.Reconcile(); err != nil {
		t.Fatalf("reconcile: %v", err)
	}
	if _, ok := manager.entries["cfg-one"]; !ok {
		t.Fatal("active backup config was not registered")
	}
}

func TestBackupReconcileRemovesStaleConfigs(t *testing.T) {
	store := newMockStore()
	store.configs["cfg-one"] = &models.BackupConfig{ID: "cfg-one", Name: "one", Status: models.BackupConfigStatusActive, BackupEnabled: true, Schedule: "0 2 * * *"}
	manager := NewBackupManager(nil, store, t.TempDir())
	if err := manager.Reconcile(); err != nil {
		t.Fatalf("reconcile: %v", err)
	}
	delete(store.configs, "cfg-one")
	if err := manager.Reconcile(); err != nil {
		t.Fatalf("reconcile after removal: %v", err)
	}
	if _, ok := manager.entries["cfg-one"]; ok {
		t.Fatal("removed backup config entry survived reconcile")
	}
}

func TestBackupReconcileSkipsManualConfigs(t *testing.T) {
	store := newMockStore()
	store.configs["cfg-manual"] = &models.BackupConfig{ID: "cfg-manual", Name: "manual", Status: models.BackupConfigStatusActive, BackupEnabled: true, Schedule: "manual"}
	manager := NewBackupManager(nil, store, t.TempDir())
	if err := manager.Reconcile(); err != nil {
		t.Fatalf("reconcile: %v", err)
	}
	if _, ok := manager.entries["cfg-manual"]; ok {
		t.Fatal("manual backup config was registered")
	}
}
