package auth

import (
	"context"
	"testing"

	"codedock/internal/models"
	"codedock/internal/repositories"
)

type fakeAuditRepo struct {
	repositories.AuditLogRepository
	gotFilter models.AuditLogFilter
}

func (f *fakeAuditRepo) ListFiltered(_ context.Context, filter models.AuditLogFilter, _, _ int) ([]models.AuditLog, error) {
	f.gotFilter = filter
	return []models.AuditLog{{ID: "log-1"}}, nil
}

func (f *fakeAuditRepo) Facets(_ context.Context, filter models.AuditLogFilter) (map[string]int, int, error) {
	f.gotFilter = filter
	return map[string]int{}, 0, nil
}

func (f *fakeAuditRepo) ListActors(_ context.Context) ([]models.AuditActor, error) {
	return []models.AuditActor{{ID: "usr-1", Name: "Alice"}}, nil
}

type fakeResourceSearcher struct {
	ids   []string
	terms []string
}

func (f *fakeResourceSearcher) SearchResourceIDs(_ context.Context, term string) []string {
	f.terms = append(f.terms, term)
	return f.ids
}

func TestListLogsByCategoryResolvesResourceIDs(t *testing.T) {
	repo := &fakeAuditRepo{}
	searcher := &fakeResourceSearcher{ids: []string{"prj-1"}}
	service := NewAuditService(repo, searcher)
	logs, err := service.ListLogsByCategory(context.Background(), "", "billing", 10, 0)
	if err != nil {
		t.Fatalf("list logs: %v", err)
	}
	if len(logs) != 1 {
		t.Fatalf("expected repo rows, got %d", len(logs))
	}
	if repo.gotFilter.Q != "billing" {
		t.Fatalf("expected q to reach repo, got %+v", repo.gotFilter)
	}
	if len(repo.gotFilter.QResourceIDs) != 1 || repo.gotFilter.QResourceIDs[0] != "prj-1" {
		t.Fatalf("expected resolved ids to reach repo, got %+v", repo.gotFilter)
	}
	if len(searcher.terms) != 1 || searcher.terms[0] != "billing" {
		t.Fatalf("expected trimmed term at searcher, got %v", searcher.terms)
	}
}

func TestListLogsWithoutQuerySkipsResolution(t *testing.T) {
	repo := &fakeAuditRepo{}
	searcher := &fakeResourceSearcher{ids: []string{"prj-1"}}
	service := NewAuditService(repo, searcher)
	if _, err := service.ListLogsByCategory(context.Background(), "servers", "   ", 10, 0); err != nil {
		t.Fatalf("list logs: %v", err)
	}
	if repo.gotFilter.Q != "" || len(repo.gotFilter.QResourceIDs) != 0 {
		t.Fatalf("expected no resolution for blank q, got %+v", repo.gotFilter)
	}
	if len(searcher.terms) != 0 {
		t.Fatalf("expected searcher untouched, got %v", searcher.terms)
	}
	if len(repo.gotFilter.Prefixes) == 0 {
		t.Fatalf("expected category prefixes to reach repo, got %+v", repo.gotFilter)
	}
}

func TestFacetsCarrySearchFilter(t *testing.T) {
	repo := &fakeAuditRepo{}
	service := NewAuditService(repo, &fakeResourceSearcher{})
	if _, err := service.Facets(context.Background(), "", "billing"); err != nil {
		t.Fatalf("facets: %v", err)
	}
	if repo.gotFilter.Q != "billing" {
		t.Fatalf("expected q to reach facets, got %+v", repo.gotFilter)
	}
}
