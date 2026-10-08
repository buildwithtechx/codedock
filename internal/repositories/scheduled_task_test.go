package repositories

import (
	"context"
	"testing"
	"time"

	"codedock.run/codedock/internal/models"
)

func TestScheduledTaskRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	for _, query := range []string{
		`INSERT INTO organizations (id, name) VALUES ('org-one', 'Organization')`,
		`INSERT INTO project_apps (id, organization_id, name, slug) VALUES ('app-one', 'org-one', 'App', 'app')`,
		`INSERT INTO projects (id, app_id, organization_id, name, slug) VALUES ('project-one', 'app-one', 'org-one', 'Project', 'project')`,
		`INSERT INTO app_services (id, project_id, name) VALUES ('service-one', 'project-one', 'Service')`,
	} {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("seed scheduled task dependencies: %v", err)
		}
	}
	ctx := context.Background()
	repo := NewScheduledTaskRepo(db)
	task := &models.ScheduledTask{ServiceID: "service-one", Name: "nightly", Schedule: "0 0 * * *", Command: "backup"}
	if err := repo.Create(ctx, task); err != nil {
		t.Fatalf("create task: %v", err)
	}
	if task.ID == "" {
		t.Fatal("expected generated task id")
	}
	loaded, err := repo.GetByID(ctx, task.ID)
	if err != nil {
		t.Fatalf("get task: %v", err)
	}
	if loaded.Name != "nightly" || loaded.Schedule != "0 0 * * *" || loaded.Command != "backup" || string(loaded.Status) != "active" {
		t.Fatalf("unexpected task row: %+v", loaded)
	}
	if loaded.CreatedAt.IsZero() || loaded.UpdatedAt.IsZero() {
		t.Fatal("expected task timestamps to roundtrip")
	}
	all, err := repo.ListAll(ctx)
	if err != nil {
		t.Fatalf("list all: %v", err)
	}
	if len(all) != 1 || all[0].ID != task.ID {
		t.Fatalf("expected 1 task, got %d", len(all))
	}
	byProject, err := repo.ListByProject(ctx, "project-one")
	if err != nil {
		t.Fatalf("list by project: %v", err)
	}
	if len(byProject) != 1 || byProject[0].ID != task.ID {
		t.Fatalf("expected 1 project task, got %d", len(byProject))
	}
	byService, err := repo.ListByService(ctx, "service-one")
	if err != nil {
		t.Fatalf("list by service: %v", err)
	}
	if len(byService) != 1 || byService[0].ID != task.ID {
		t.Fatalf("expected 1 service task, got %d", len(byService))
	}
	task.Name = "nightly-renamed"
	task.Schedule = "0 1 * * *"
	if err := repo.Update(ctx, task); err != nil {
		t.Fatalf("update task: %v", err)
	}
	updated, err := repo.GetByID(ctx, task.ID)
	if err != nil {
		t.Fatalf("get updated task: %v", err)
	}
	if updated.Name != "nightly-renamed" || updated.Schedule != "0 1 * * *" {
		t.Fatalf("update did not persist: %+v", updated)
	}
	ranAt := time.Now().UTC().Truncate(time.Second)
	if err := repo.UpdateStatus(ctx, task.ID, models.ScheduledTaskStatusCompleted, &ranAt, "ok"); err != nil {
		t.Fatalf("update status: %v", err)
	}
	completed, err := repo.GetByID(ctx, task.ID)
	if err != nil {
		t.Fatalf("get completed task: %v", err)
	}
	if string(completed.Status) != "completed" || completed.LastOutput != "ok" || !completed.LastRunAt.Truncate(time.Second).Equal(ranAt) {
		t.Fatalf("status update did not persist: %+v", completed)
	}
	if err := repo.Delete(ctx, task.ID); err != nil {
		t.Fatalf("delete task: %v", err)
	}
	if _, err := repo.GetByID(ctx, task.ID); err == nil {
		t.Fatal("expected deleted task to be gone")
	}
}
