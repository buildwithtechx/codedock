package repositories_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"

	_ "modernc.org/sqlite"
)

func TestDeploymentRepositoryListByOrganization(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:?_pragma=foreign_keys(ON)")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer db.Close()
	if err := repositories.RunMigrations(db); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	ctx := context.Background()
	now := time.Now().UTC()
	seedDeployment(t, db, "org-one", "project-one", "service-one", "alpha", "api", "deployment-one", "READY", "main", "abc1234", now)
	seedDeployment(t, db, "org-one", "project-two", "service-two", "beta", "worker", "deployment-two", "FAILED", "release", "def5678", now.Add(time.Minute))
	seedDeployment(t, db, "org-two", "project-three", "service-three", "gamma", "api", "deployment-three", "READY", "main", "ghi9012", now.Add(2*time.Minute))

	repo := repositories.NewDeploymentRepo(db)
	deployments, total, err := repo.ListByOrganization(ctx, models.DeploymentListFilter{
		OrganizationID: "org-one",
		Limit:          25,
	})
	if err != nil {
		t.Fatalf("list organization deployments: %v", err)
	}
	if total != 2 || len(deployments) != 2 {
		t.Fatalf("expected two organization deployments, got total=%d rows=%d", total, len(deployments))
	}
	if deployments[0].ProjectName != "beta" || deployments[0].ServiceName != "worker" {
		t.Fatalf("expected newest deployment labels, got %+v", deployments[0])
	}

	filtered, filteredTotal, err := repo.ListByOrganization(ctx, models.DeploymentListFilter{
		OrganizationID: "org-one",
		Status:         "READY",
		Search:         "api",
		Limit:          25,
	})
	if err != nil {
		t.Fatalf("filter organization deployments: %v", err)
	}
	if filteredTotal != 1 || len(filtered) != 1 || filtered[0].ID != "deployment-one" {
		t.Fatalf("unexpected filtered deployments: total=%d rows=%+v", filteredTotal, filtered)
	}
}

func seedDeployment(
	t *testing.T,
	db *sql.DB,
	organizationID, projectID, serviceID, projectName, serviceName, deploymentID, status, branch, commitHash string,
	createdAt time.Time,
) {
	t.Helper()
	if _, err := db.Exec(`INSERT OR IGNORE INTO organizations (id, name) VALUES (?, ?)`, organizationID, organizationID); err != nil {
		t.Fatalf("create organization: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO projects (id, name, organization_id) VALUES (?, ?, ?)`, projectID, projectName, organizationID); err != nil {
		t.Fatalf("create project: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO app_services (id, project_id, environment_id, name, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`, serviceID, projectID, "environment-"+projectID, serviceName, createdAt, createdAt); err != nil {
		t.Fatalf("create service: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO deployments (id, service_id, environment_id, project_id, status, commit_hash, branch, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, deploymentID, serviceID, "environment-"+projectID, projectID, status, commitHash, branch, createdAt, createdAt); err != nil {
		t.Fatalf("create deployment: %v", err)
	}
}
