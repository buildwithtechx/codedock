package repositories_test

import (
	"context"
	"database/sql"
	"testing"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"

	_ "modernc.org/sqlite"
)

func TestDeploymentRepositoryCreateAndRead(t *testing.T) {
	for _, serviceID := range []string{"service-one", ""} {
		t.Run("service="+serviceID, func(t *testing.T) {
			db, err := sql.Open("sqlite", ":memory:?_pragma=foreign_keys(ON)")
			if err != nil {
				t.Fatalf("open database: %v", err)
			}
			db.SetMaxOpenConns(1)
			t.Cleanup(func() {
				if err := db.Close(); err != nil {
					t.Errorf("close database: %v", err)
				}
			})
			if err := repositories.RunMigrations(db); err != nil {
				t.Fatalf("run migrations: %v", err)
			}
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
			repo := repositories.NewDeploymentRepo(db)
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
				if _, err := db.Exec(`DELETE FROM app_services WHERE id = ?`, serviceID); err != nil {
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
