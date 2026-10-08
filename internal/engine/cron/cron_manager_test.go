package cron

import (
	"testing"
	"time"

	"codedock.run/codedock/internal/models"
)

type reconcileStore struct {
	tasks []models.ScheduledTask
}

func (s *reconcileStore) ListScheduledTasks() ([]models.ScheduledTask, error) {
	return s.tasks, nil
}

func (s *reconcileStore) GetScheduledTask(id string) (*models.ScheduledTask, error) {
	return nil, nil
}

func (s *reconcileStore) GetProject(id string) (*models.ProjectConfig, error) {
	return nil, nil
}

func (s *reconcileStore) GetAppService(id string) (*models.AppService, error) {
	return nil, nil
}

func (s *reconcileStore) UpdateScheduledTaskStatusAndOutput(id string, status models.ScheduledTaskStatus, lastRunAt *time.Time, output string) error {
	return nil
}

func TestReconcileRegistersActiveTasks(t *testing.T) {
	store := &reconcileStore{tasks: []models.ScheduledTask{
		{ID: "task-one", Name: "one", Schedule: "0 * * * *", Command: "echo one", Status: "active"},
		{ID: "task-two", Name: "two", Schedule: "0 * * * *", Command: "echo two", Status: "paused"},
	}}
	manager := NewCronManager(nil, store)
	if err := manager.Reconcile(); err != nil {
		t.Fatalf("reconcile: %v", err)
	}
	if _, ok := manager.entries["task-one"]; !ok {
		t.Fatal("active task was not registered")
	}
	if _, ok := manager.entries["task-two"]; ok {
		t.Fatal("paused task was registered")
	}
}

func TestReconcileRemovesStaleTasks(t *testing.T) {
	store := &reconcileStore{tasks: []models.ScheduledTask{
		{ID: "task-one", Name: "one", Schedule: "0 * * * *", Command: "echo one", Status: "active"},
	}}
	manager := NewCronManager(nil, store)
	if err := manager.Reconcile(); err != nil {
		t.Fatalf("reconcile: %v", err)
	}
	store.tasks = nil
	if err := manager.Reconcile(); err != nil {
		t.Fatalf("reconcile after removal: %v", err)
	}
	if _, ok := manager.entries["task-one"]; ok {
		t.Fatal("removed task entry survived reconcile")
	}
}

func TestReconcilePreservesHousekeepingEntries(t *testing.T) {
	store := &reconcileStore{}
	manager := NewCronManager(nil, store)
	if err := manager.ScheduleDockerCleanup("0 * * * *"); err != nil {
		t.Fatalf("schedule docker cleanup: %v", err)
	}
	if err := manager.ScheduleDiskUsageCheck("0 * * * *", 80); err != nil {
		t.Fatalf("schedule disk usage check: %v", err)
	}
	if err := manager.Reconcile(); err != nil {
		t.Fatalf("reconcile: %v", err)
	}
	if _, ok := manager.entries[dockerCleanupEntryID]; !ok {
		t.Fatal("docker cleanup entry was removed by reconcile")
	}
	if _, ok := manager.entries[diskUsageEntryID]; !ok {
		t.Fatal("disk usage entry was removed by reconcile")
	}
}
