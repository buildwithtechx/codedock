package repositories

import (
	"context"
	"testing"
)

func TestCanvasSummariesResolveDeployTargets(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	repo := NewCanvasRepo(db, NewEnvironmentRepo(db))
	if _, err := db.Exec(`INSERT INTO organizations (id, name) VALUES ('org-1', 'Acme')`); err != nil {
		t.Fatalf("seed organization: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO servers (id, organization_id, user_id, name, is_local) VALUES
		('srv-local', 'org-1', 'user-1', 'Local', TRUE),
		('srv-remote', 'org-1', 'user-1', 'Hetzner 1', FALSE)`); err != nil {
		t.Fatalf("seed servers: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO project_apps (id, organization_id, name, slug) VALUES
		('app-1', 'org-1', 'App 1', 'app-1'),
		('app-2', 'org-1', 'App 2', 'app-2'),
		('app-3', 'org-1', 'App 3', 'app-3')`); err != nil {
		t.Fatalf("seed project apps: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO projects (id, app_id, organization_id, server_id, name, slug) VALUES
		('proj-local', 'app-1', 'org-1', NULL, 'Local project', 'local-project'),
		('proj-on-local-server', 'app-2', 'org-1', 'srv-local', 'Local server project', 'local-server-project'),
		('proj-remote', 'app-3', 'org-1', 'srv-remote', 'Remote project', 'remote-project')`); err != nil {
		t.Fatalf("seed projects: %v", err)
	}

	summaries, err := repo.ListCanvasSummaries(ctx, "org-1")
	if err != nil {
		t.Fatalf("list canvas summaries: %v", err)
	}
	byID := make(map[string]string, len(summaries))
	names := make(map[string]string, len(summaries))
	for _, summary := range summaries {
		byID[summary.ID] = summary.DeployTarget
		names[summary.ID] = summary.ServerName
	}
	if byID["proj-local"] != "local" || names["proj-local"] != "" {
		t.Fatalf("unassigned project target = %q name = %q, want local/empty", byID["proj-local"], names["proj-local"])
	}
	if byID["proj-on-local-server"] != "local" || names["proj-on-local-server"] != "" {
		t.Fatalf("local-server project target = %q name = %q, want local/empty", byID["proj-on-local-server"], names["proj-on-local-server"])
	}
	if byID["proj-remote"] != "server" || names["proj-remote"] != "Hetzner 1" {
		t.Fatalf("remote project target = %q name = %q, want server/Hetzner 1", byID["proj-remote"], names["proj-remote"])
	}

	single, err := repo.GetCanvasSummary(ctx, "proj-remote")
	if err != nil {
		t.Fatalf("get canvas summary: %v", err)
	}
	if single.DeployTarget != "server" || single.ServerName != "Hetzner 1" {
		t.Fatalf("single summary target = %q name = %q, want server/Hetzner 1", single.DeployTarget, single.ServerName)
	}
}
