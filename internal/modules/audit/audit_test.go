package audit

import (
	"context"
	"errors"
	"testing"
)

type fakeAuditRepo struct {
	AuditLogRepository
	gotFilter AuditLogFilter
}

func (f *fakeAuditRepo) ListFiltered(_ context.Context, filter AuditLogFilter, _, _ int) ([]AuditLog, error) {
	f.gotFilter = filter
	return []AuditLog{{ID: "log-1"}}, nil
}

func (f *fakeAuditRepo) Facets(_ context.Context, filter AuditLogFilter) (map[string]int, int, error) {
	f.gotFilter = filter
	return map[string]int{}, 0, nil
}

func (f *fakeAuditRepo) ListActors(_ context.Context) ([]AuditActor, error) {
	return []AuditActor{{ID: "usr-1", Name: "Alice"}}, nil
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

type fakeSearchRepo struct {
	ids   []string
	err   error
	terms []string
}

func (f *fakeSearchRepo) SearchIDsByName(_ context.Context, term string, limit int) ([]string, error) {
	f.terms = append(f.terms, term)
	if limit != AuditSearchResolverLimit {
		return nil, errors.New("unexpected limit")
	}
	return f.ids, f.err
}

func TestSearchResourceIDsMergesAndDedupes(t *testing.T) {
	projects := &fakeSearchRepo{ids: []string{"prj-1", "prj-2"}}
	searcher := NewAuditResourceSearcher(
		projects,
		&fakeSearchRepo{ids: []string{"prj-2", "srv-1"}},
		&fakeSearchRepo{ids: []string{}},
		&fakeSearchRepo{ids: []string{"svc-1"}},
	)
	ids := searcher.SearchResourceIDs(context.Background(), "api")
	want := []string{"prj-1", "prj-2", "srv-1", "svc-1"}
	if len(ids) != len(want) {
		t.Fatalf("expected %v, got %v", want, ids)
	}
	for i := range want {
		if ids[i] != want[i] {
			t.Fatalf("expected %v, got %v", want, ids)
		}
	}
	if len(projects.terms) != 1 || projects.terms[0] != "api" {
		t.Fatalf("expected term to reach resolvers, got %v", projects.terms)
	}
}

func TestSearchResourceIDsIsolatesResolverFailures(t *testing.T) {
	searcher := NewAuditResourceSearcher(
		&fakeSearchRepo{err: errors.New("projects down")},
		&fakeSearchRepo{ids: []string{"srv-1"}},
		&fakeSearchRepo{err: errors.New("databases down")},
		&fakeSearchRepo{ids: []string{}},
	)
	ids := searcher.SearchResourceIDs(context.Background(), "api")
	if len(ids) != 1 || ids[0] != "srv-1" {
		t.Fatalf("expected surviving resolver results, got %v", ids)
	}
}
