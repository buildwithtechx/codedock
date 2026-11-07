package repositories

import (
	"context"
	"testing"
	"time"

	"codedock/internal/models"
	"codedock/internal/utils"
)

func TestProjectRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	if _, err := db.Exec(`INSERT INTO organizations (id, name) VALUES ('org-one', 'Organization')`); err != nil {
		t.Fatalf("seed organization: %v", err)
	}
	ctx := context.Background()
	repo := NewProjectRepo(db, NewEnvironmentRepo(db))
	project := &models.ProjectConfig{OrganizationID: "org-one", Name: "My Project"}
	if err := repo.Create(ctx, project); err != nil {
		t.Fatalf("create project: %v", err)
	}
	if project.ID == "" || project.AppID == "" || project.Slug != "my-project" {
		t.Fatalf("expected generated ids and slug, got %+v", project)
	}
	if project.CreatedAt.IsZero() || project.UpdatedAt.IsZero() {
		t.Fatal("expected timestamps to roundtrip")
	}
	got, err := repo.Get(ctx, project.ID)
	if err != nil {
		t.Fatalf("get project: %v", err)
	}
	if got == nil || got.Name != "My Project" || got.OrganizationID != "org-one" {
		t.Fatalf("unexpected project row: %+v", got)
	}
	if !got.CreatedAt.Truncate(time.Second).Equal(project.CreatedAt.Truncate(time.Second)) {
		t.Fatal("created_at did not roundtrip")
	}
	scoped, err := repo.GetByOrganization(ctx, project.ID, "org-one")
	if err != nil || scoped == nil {
		t.Fatalf("get by organization: %v %+v", err, scoped)
	}
	missing, err := repo.GetByOrganization(ctx, project.ID, "org-other")
	if err != nil || missing != nil {
		t.Fatalf("expected nil for wrong organization, got %+v %v", missing, err)
	}
	projects, total, err := repo.ListByOrganization(ctx, "org-one", "", 10, 0)
	if err != nil {
		t.Fatalf("list by organization: %v", err)
	}
	if total != 1 || len(projects) != 1 || projects[0].ID != project.ID {
		t.Fatalf("unexpected organization listing: total=%d len=%d", total, len(projects))
	}
	all, total, err := repo.ListAll(ctx, "", 10, 0)
	if err != nil {
		t.Fatalf("list all: %v", err)
	}
	if total != 1 || len(all) != 1 {
		t.Fatalf("unexpected full listing: total=%d len=%d", total, len(all))
	}
	for _, query := range []string{
		`INSERT INTO users (id, email, password_hash) VALUES ('user-one', 'ada@example.com', 'hash')`,
		`INSERT INTO organization_members (id, organization_id, user_id, email) VALUES ('member-one', 'org-one', 'user-one', 'ada@example.com')`,
	} {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("seed membership: %v", err)
		}
	}
	count, err := repo.CountByUser(ctx, "user-one")
	if err != nil {
		t.Fatalf("count by user: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 project for user, got %d", count)
	}
	if _, err := db.Exec(`INSERT INTO servers (id, user_id, name) VALUES ('server-one', 'user-one', 'Server')`); err != nil {
		t.Fatalf("seed server: %v", err)
	}
	if err := repo.SetServer(ctx, project.ID, "server-one"); err != nil {
		t.Fatalf("set server: %v", err)
	}
	pinned, err := repo.Get(ctx, project.ID)
	if err != nil || pinned == nil || pinned.ServerID != "server-one" {
		t.Fatalf("server was not pinned: %+v %v", pinned, err)
	}
	if err := repo.SetServer(ctx, project.ID, ""); err != nil {
		t.Fatalf("clear server: %v", err)
	}
	unpinned, err := repo.Get(ctx, project.ID)
	if err != nil || unpinned == nil || unpinned.ServerID != "" {
		t.Fatalf("server was not cleared: %+v %v", unpinned, err)
	}
	envs, err := repo.environments.ListByProject(ctx, project.ID)
	if err != nil {
		t.Fatalf("list default environments: %v", err)
	}
	if len(envs) != 1 || !envs[0].IsDefault {
		t.Fatalf("expected default environment, got %+v", envs)
	}
	vault, err := utils.NewVault(t.TempDir())
	if err != nil {
		t.Fatalf("new vault: %v", err)
	}
	vars := NewEnvRepo(db, vault)
	if err := vars.SetVar(ctx, project.ID, "API_KEY", "secret"); err != nil {
		t.Fatalf("set var: %v", err)
	}
	values, err := vars.GetVars(ctx, project.ID)
	if err != nil {
		t.Fatalf("get vars: %v", err)
	}
	if values["API_KEY"] != "secret" {
		t.Fatalf("unexpected vars: %v", values)
	}
	if err := vars.SetVar(ctx, project.ID, "API_KEY", "rotated"); err != nil {
		t.Fatalf("overwrite var: %v", err)
	}
	values, err = vars.GetVars(ctx, project.ID)
	if err != nil || values["API_KEY"] != "rotated" {
		t.Fatalf("var overwrite did not persist: %v %v", values, err)
	}
	if err := repo.Delete(ctx, project.ID); err != nil {
		t.Fatalf("delete project: %v", err)
	}
	deleted, err := repo.Get(ctx, project.ID)
	if err != nil || deleted != nil {
		t.Fatalf("expected deleted project to be gone, got %+v %v", deleted, err)
	}
}

func TestProjectSearchIDsByName(t *testing.T) {
	db := openPGTestDB(t)
	if _, err := db.Exec(`INSERT INTO organizations (id, name) VALUES ('org-one', 'Organization')`); err != nil {
		t.Fatalf("seed organization: %v", err)
	}
	ctx := context.Background()
	repo := NewProjectRepo(db, NewEnvironmentRepo(db))
	billing := &models.ProjectConfig{OrganizationID: "org-one", Name: "Billing API"}
	if err := repo.Create(ctx, billing); err != nil {
		t.Fatalf("create project: %v", err)
	}
	marketing := &models.ProjectConfig{OrganizationID: "org-one", Name: "Marketing Site"}
	if err := repo.Create(ctx, marketing); err != nil {
		t.Fatalf("create project: %v", err)
	}
	ids, err := repo.SearchIDsByName(ctx, "billing", 0)
	if err != nil {
		t.Fatalf("search ids: %v", err)
	}
	if len(ids) != 1 || ids[0] != billing.ID {
		t.Fatalf("expected only the billing project id, got %v", ids)
	}
	ids, err = repo.SearchIDsByName(ctx, "nothing-matches-this", 0)
	if err != nil {
		t.Fatalf("search ids: %v", err)
	}
	if len(ids) != 0 {
		t.Fatalf("expected no ids, got %v", ids)
	}
}

func TestProjectSearchFiltersByNameSlugAndDescription(t *testing.T) {
	db := openPGTestDB(t)
	if _, err := db.Exec(`INSERT INTO organizations (id, name) VALUES ('org-one', 'Organization')`); err != nil {
		t.Fatalf("seed organization: %v", err)
	}
	ctx := context.Background()
	repo := NewProjectRepo(db, NewEnvironmentRepo(db))
	for _, seed := range []models.ProjectConfig{
		{OrganizationID: "org-one", Name: "Billing API", Description: "handles invoices"},
		{OrganizationID: "org-one", Name: "Marketing Site"},
	} {
		seed := seed
		if err := repo.Create(ctx, &seed); err != nil {
			t.Fatalf("create project: %v", err)
		}
	}
	cases := []struct {
		name   string
		search string
		want   int
	}{
		{name: "name case-insensitive", search: "BILLING", want: 1},
		{name: "description", search: "invoices", want: 1},
		{name: "slug", search: "marketing-site", want: 1},
		{name: "no match", search: "nothing-matches-this", want: 0},
		{name: "underscore literal", search: "billin_", want: 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			projects, total, err := repo.ListByOrganization(ctx, "org-one", tc.search, 10, 0)
			if err != nil {
				t.Fatalf("search projects: %v", err)
			}
			if total != tc.want || len(projects) != tc.want {
				t.Fatalf("expected %d matches for %q, got total=%d rows=%d", tc.want, tc.search, total, len(projects))
			}
			all, allTotal, err := repo.ListAll(ctx, tc.search, 10, 0)
			if err != nil {
				t.Fatalf("search all projects: %v", err)
			}
			if allTotal != tc.want || len(all) != tc.want {
				t.Fatalf("expected %d matches for %q, got total=%d rows=%d", tc.want, tc.search, allTotal, len(all))
			}
		})
	}
}
