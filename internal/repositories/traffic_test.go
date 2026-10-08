package repositories

import (
	"context"
	"testing"

	"codedock.run/codedock/internal/models"
)

func seedTrafficSamples() []models.TrafficSample {
	return []models.TrafficSample{
		{Time: "2026-10-07T08:00:00Z", ProjectID: "proj-1", Domain: "shop.example.com", Path: "/cart", Status: 200, Bytes: 100, DurationMs: 10, ClientIP: "10.0.0.2"},
		{Time: "2026-10-07T08:00:30Z", ProjectID: "proj-1", Domain: "shop.example.com", Path: "/cart", Status: 500, Bytes: 200, DurationMs: 30, ClientIP: "10.0.0.3"},
		{Time: "2026-10-07T08:01:00Z", ProjectID: "proj-1", Domain: "shop.example.com", Path: "/home", Status: 200, Bytes: 50, DurationMs: 5, ClientIP: "10.0.0.2"},
	}
}

func TestTrafficRecordAndSummary(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	repo := NewTrafficRepository(db)
	if err := repo.RecordBatch(ctx, seedTrafficSamples()); err != nil {
		t.Fatalf("record: %v", err)
	}
	summary, err := repo.Summary(ctx, "proj-1", "", "2026-10-07T08:00", "2026-10-07T08:02")
	if err != nil {
		t.Fatalf("summary: %v", err)
	}
	if summary.Requests != 3 || summary.Bytes != 350 {
		t.Fatalf("expected 3 requests/350 bytes, got %#v", summary)
	}
	if summary.ErrorRate < 0.33 || summary.ErrorRate > 0.34 {
		t.Fatalf("expected ~0.33 error rate, got %#v", summary)
	}
	if summary.AvgDurationMs != 15 {
		t.Fatalf("expected 15ms average, got %#v", summary)
	}
	overview, err := repo.Overview(ctx, "proj-1", "", "2026-10-07T08:00", "2026-10-07T08:02", "1")
	if err != nil {
		t.Fatalf("overview: %v", err)
	}
	if len(overview.Series) != 2 || len(overview.Statuses) != 2 || len(overview.TopPaths) != 2 {
		t.Fatalf("expected series/statuses/paths, got %#v", overview)
	}
	if overview.TopPaths[0].Path != "/cart" || overview.TopPaths[0].Requests != 2 {
		t.Fatalf("expected /cart top path, got %#v", overview.TopPaths)
	}
	geo, err := repo.Geo(ctx, "proj-1", "2026-10-07", "2026-10-07")
	if err != nil {
		t.Fatalf("geo: %v", err)
	}
	if len(geo.Countries) != 1 || geo.Countries[0].Visitors != 2 || geo.Countries[0].Country != "unknown" {
		t.Fatalf("expected unknown country with 2 visitors, got %#v", geo)
	}
	enabled, err := repo.PathsEnabled(ctx, "proj-1")
	if err != nil || enabled {
		t.Fatalf("expected paths disabled by default, got %v %v", enabled, err)
	}
	if err := repo.SetPathsEnabled(ctx, "proj-1", true); err != nil {
		t.Fatalf("enable paths: %v", err)
	}
	if enabled, err := repo.PathsEnabled(ctx, "proj-1"); err != nil || !enabled {
		t.Fatalf("expected paths enabled, got %v %v", enabled, err)
	}
	if err := repo.DeleteBefore(ctx, "2026-10-07T08:01"); err != nil {
		t.Fatalf("retire: %v", err)
	}
	after, err := repo.Summary(ctx, "proj-1", "", "2026-10-07T08:00", "2026-10-07T08:02")
	if err != nil {
		t.Fatalf("summary after retire: %v", err)
	}
	if after.Requests != 1 {
		t.Fatalf("expected 1 request after retire, got %#v", after)
	}
}

func TestAttentionIssueLifecycle(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	orgRepo := NewOrganizationRepository(db)
	if err := orgRepo.Create(ctx, &models.Organization{ID: "org-1", Name: "Acme"}); err != nil {
		t.Fatalf("create org: %v", err)
	}
	repo := NewAttentionRepository(db)
	issue := &models.AttentionIssue{ID: "iss-1", OrganizationID: "org-1", Kind: "deployment-failed", Subject: "svc-1", Severity: models.AttentionSeverityCritical, Title: "Deploy failing", Detail: "log", Remediation: "fix", Action: "redeploy"}
	if err := repo.Upsert(ctx, issue); err != nil {
		t.Fatalf("upsert: %v", err)
	}
	again := &models.AttentionIssue{ID: "iss-2", OrganizationID: "org-1", Kind: "deployment-failed", Subject: "svc-1", Severity: models.AttentionSeverityCritical, Title: "Deploy failing", Detail: "log2", Remediation: "fix"}
	if err := repo.Upsert(ctx, again); err != nil {
		t.Fatalf("reupsert: %v", err)
	}
	listed, err := repo.List(ctx, "org-1", "open", 10)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(listed) != 1 || listed[0].Occurrences != 2 || listed[0].Detail != "log2" {
		t.Fatalf("expected deduped issue with 2 occurrences, got %#v", listed)
	}
	open, err := repo.CountOpen(ctx, "org-1")
	if err != nil || open != 1 {
		t.Fatalf("expected 1 open, got %d %v", open, err)
	}
	if err := repo.SetStatus(ctx, listed[0].ID, "org-1", models.AttentionStatusAcked); err != nil {
		t.Fatalf("ack: %v", err)
	}
	if err := repo.ResolveMissing(ctx, "org-1", []string{"other:live"}); err != nil {
		t.Fatalf("reconcile: %v", err)
	}
	loaded, err := repo.Get(ctx, listed[0].ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if loaded.Status != models.AttentionStatusAcked {
		t.Fatalf("expected acked to survive reconcile, got %s", loaded.Status)
	}
	if err := repo.SetStatus(ctx, listed[0].ID, "org-1", models.AttentionStatusOpen); err != nil {
		t.Fatalf("reopen: %v", err)
	}
	if err := repo.ResolveMissing(ctx, "org-1", []string{"other:live"}); err != nil {
		t.Fatalf("reconcile open: %v", err)
	}
	loaded, err = repo.Get(ctx, listed[0].ID)
	if err != nil {
		t.Fatalf("get after reconcile: %v", err)
	}
	if loaded.Status != models.AttentionStatusResolved {
		t.Fatalf("expected missing open issue resolved, got %s", loaded.Status)
	}
	if err := repo.AppendDetail(ctx, listed[0].ID, "action ran"); err != nil {
		t.Fatalf("append: %v", err)
	}
}
