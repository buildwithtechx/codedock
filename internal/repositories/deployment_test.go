package repositories

import (
	"context"
	"testing"

	"codedock/internal/models"
)

func TestDeploymentRepositoryCreateAndRead(t *testing.T) {
	for _, serviceID := range []string{"service-one", ""} {
		t.Run("service="+serviceID, func(t *testing.T) {
			db := openPGTestDB(t)
			for _, query := range []string{
				`INSERT INTO organizations (id, name) VALUES ('org-one', 'Organization')`,
				`INSERT INTO project_apps (id, organization_id, name, slug) VALUES ('app-one', 'org-one', 'App', 'app')`,
				`INSERT INTO projects (id, app_id, organization_id, name, slug) VALUES ('project-one', 'app-one', 'org-one', 'Project', 'project')`,
				`INSERT INTO app_services (id, project_id, environment_id, name) VALUES ('service-one', 'project-one', 'environment-one', 'Service')`,
			} {
				if _, err := db.Exec(query); err != nil {
					t.Fatalf("seed deployment dependencies: %v", err)
				}
			}
			ctx := context.Background()
			repo := NewDeploymentRepo(db)
			deployment := &models.Deployment{
				ProjectID:  "project-one",
				ServiceID:  serviceID,
				Branch:     "main",
				Trigger:    "manual",
				CommitHash: "abc1234",
			}
			if err := repo.Create(ctx, deployment); err != nil {
				t.Fatalf("create deployment without organization: %v", err)
			}
			if deployment.OrganizationID != "org-one" {
				t.Fatalf("unexpected organization: %q", deployment.OrganizationID)
			}
			stored, err := repo.GetByID(ctx, deployment.ID)
			if err != nil {
				t.Fatalf("read deployment: %v", err)
			}
			expectedEnvironment := ""
			if serviceID != "" {
				expectedEnvironment = "environment-one"
			}
			if stored.OrganizationID != "org-one" || stored.EnvironmentID != expectedEnvironment || stored.ServiceID != serviceID || stored.CommitHash != "abc1234" {
				t.Fatalf("unexpected stored deployment: %+v", stored)
			}
			deployments, total, err := repo.ListByService(ctx, serviceID, 25, 0)
			if err != nil {
				t.Fatalf("list service deployments: %v", err)
			}
			if serviceID == "" && (total != 0 || len(deployments) != 0) {
				t.Fatalf("unassigned deployments leaked into service listing: %d %+v", total, deployments)
			}
			if serviceID != "" && (total != 1 || len(deployments) != 1 || deployments[0].EnvironmentID != expectedEnvironment || deployments[0].OrganizationID != "org-one") {
				t.Fatalf("unexpected service deployments: total=%d rows=%+v", total, deployments)
			}
			items, total, err := repo.ListByOrganization(ctx, models.DeploymentListFilter{OrganizationID: "org-one", Limit: 25})
			if err != nil {
				t.Fatalf("list organization deployments: %v", err)
			}
			if total != 1 || len(items) != 1 || items[0].EnvironmentID != expectedEnvironment {
				t.Fatalf("unexpected organization deployments: total=%d rows=%+v", total, items)
			}
			if serviceID != "" {
				if _, err := db.Exec(`DELETE FROM app_services WHERE id = $1`, serviceID); err != nil {
					t.Fatalf("delete service: %v", err)
				}
				stored, err = repo.GetByID(ctx, deployment.ID)
				if err != nil {
					t.Fatalf("read deployment after service deletion: %v", err)
				}
				if stored.ServiceID != "" || stored.EnvironmentID != "" {
					t.Fatalf("unexpected deleted service metadata: %+v", stored)
				}
			}
		})
	}
}

func TestDeploymentUpdateStatusAndRecovery(t *testing.T) {
	db := openPGTestDB(t)
	for _, query := range []string{
		`INSERT INTO organizations (id, name) VALUES ('org-one', 'Organization')`,
		`INSERT INTO project_apps (id, organization_id, name, slug) VALUES ('app-one', 'org-one', 'App', 'app')`,
		`INSERT INTO projects (id, app_id, organization_id, name, slug) VALUES ('project-one', 'app-one', 'org-one', 'Project', 'project')`,
		`INSERT INTO app_services (id, project_id, environment_id, name) VALUES ('service-one', 'project-one', 'environment-one', 'Service')`,
	} {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("seed deployment dependencies: %v", err)
		}
	}
	ctx := context.Background()
	repo := NewDeploymentRepo(db)
	deployment := &models.Deployment{ProjectID: "project-one", ServiceID: "service-one", Branch: "main", Trigger: "manual"}
	if err := repo.Create(ctx, deployment); err != nil {
		t.Fatalf("create deployment: %v", err)
	}
	deployment.Branch = "feature"
	deployment.CommitHash = "def5678"
	if err := repo.Update(ctx, deployment); err != nil {
		t.Fatalf("update deployment: %v", err)
	}
	updated, err := repo.GetByID(ctx, deployment.ID)
	if err != nil {
		t.Fatalf("read updated deployment: %v", err)
	}
	if updated.Branch != "feature" || updated.CommitHash != "def5678" {
		t.Fatalf("update did not persist: %+v", updated)
	}
	if err := repo.UpdateStatus(ctx, deployment.ID, models.DeploymentStatusBuilding, "building", "container-one"); err != nil {
		t.Fatalf("update to building: %v", err)
	}
	building, err := repo.GetByID(ctx, deployment.ID)
	if err != nil {
		t.Fatalf("read building deployment: %v", err)
	}
	if building.FinishedAt != nil {
		t.Fatal("non-terminal status should not set finished_at")
	}
	if err := repo.UpdateStatus(ctx, deployment.ID, models.DeploymentStatusReady, "done", "container-one"); err != nil {
		t.Fatalf("update to ready: %v", err)
	}
	ready, err := repo.GetByID(ctx, deployment.ID)
	if err != nil {
		t.Fatalf("read ready deployment: %v", err)
	}
	if ready.FinishedAt == nil {
		t.Fatal("terminal status should set finished_at")
	}
	interrupted := &models.Deployment{ProjectID: "project-one", ServiceID: "service-one", Branch: "main", Trigger: "manual"}
	if err := repo.Create(ctx, interrupted); err != nil {
		t.Fatalf("create interrupted deployment: %v", err)
	}
	if err := repo.RecoverInterrupted(ctx); err != nil {
		t.Fatalf("recover interrupted: %v", err)
	}
	recovered, err := repo.GetByID(ctx, interrupted.ID)
	if err != nil {
		t.Fatalf("read recovered deployment: %v", err)
	}
	if string(recovered.Status) != "FAILED" {
		t.Fatalf("expected FAILED after recovery, got %q", recovered.Status)
	}
}

func TestDeploymentSearchIsCaseInsensitive(t *testing.T) {
	db := openPGTestDB(t)
	for _, query := range []string{
		`INSERT INTO organizations (id, name) VALUES ('org-one', 'Organization')`,
		`INSERT INTO project_apps (id, organization_id, name, slug) VALUES ('app-one', 'org-one', 'App', 'app')`,
		`INSERT INTO projects (id, app_id, organization_id, name, slug) VALUES ('project-one', 'app-one', 'org-one', 'Project', 'project')`,
		`INSERT INTO app_services (id, project_id, environment_id, name) VALUES ('service-one', 'project-one', 'environment-one', 'BillingService')`,
	} {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("seed deployment dependencies: %v", err)
		}
	}
	ctx := context.Background()
	repo := NewDeploymentRepo(db)
	deployment := &models.Deployment{ProjectID: "project-one", ServiceID: "service-one", Branch: "main", Trigger: "manual"}
	if err := repo.Create(ctx, deployment); err != nil {
		t.Fatalf("create deployment: %v", err)
	}
	items, total, err := repo.ListByOrganization(ctx, models.DeploymentListFilter{OrganizationID: "org-one", Search: "BILLING", Limit: 25})
	if err != nil {
		t.Fatalf("search deployments: %v", err)
	}
	if total != 1 || len(items) != 1 {
		t.Fatalf("expected case-insensitive match, got total=%d rows=%d", total, len(items))
	}
	items, total, err = repo.ListByOrganization(ctx, models.DeploymentListFilter{OrganizationID: "org-one", Search: "nothing-matches-this", Limit: 25})
	if err != nil {
		t.Fatalf("search deployments: %v", err)
	}
	if total != 0 || len(items) != 0 {
		t.Fatalf("expected no matches, got total=%d rows=%d", total, len(items))
	}
}
