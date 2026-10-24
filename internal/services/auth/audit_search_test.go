package auth

import (
	"context"
	"errors"
	"testing"

	"codedock/internal/repositories"
)

type fakeSearchRepo struct {
	repositories.ProjectRepository
	ids   []string
	err   error
	terms []string
}

func (f *fakeSearchRepo) SearchIDsByName(_ context.Context, term string, limit int) ([]string, error) {
	f.terms = append(f.terms, term)
	if limit != auditSearchResolverLimit {
		return nil, errors.New("unexpected limit")
	}
	return f.ids, f.err
}

type fakeServerSearchRepo struct {
	repositories.ServerRepository
	ids []string
	err error
}

func (f *fakeServerSearchRepo) SearchIDsByName(_ context.Context, _ string, _ int) ([]string, error) {
	return f.ids, f.err
}

type fakeDatabaseSearchRepo struct {
	repositories.DatabaseRepository
	ids []string
	err error
}

func (f *fakeDatabaseSearchRepo) SearchIDsByName(_ context.Context, _ string, _ int) ([]string, error) {
	return f.ids, f.err
}

type fakeAppSearchRepo struct {
	repositories.AppServiceRepository
	ids []string
	err error
}

func (f *fakeAppSearchRepo) SearchIDsByName(_ context.Context, _ string, _ int) ([]string, error) {
	return f.ids, f.err
}

func TestSearchResourceIDsMergesAndDedupes(t *testing.T) {
	projects := &fakeSearchRepo{ids: []string{"prj-1", "prj-2"}}
	searcher := NewAuditResourceSearcher(
		projects,
		&fakeServerSearchRepo{ids: []string{"prj-2", "srv-1"}},
		&fakeDatabaseSearchRepo{ids: []string{}},
		&fakeAppSearchRepo{ids: []string{"svc-1"}},
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
		&fakeServerSearchRepo{ids: []string{"srv-1"}},
		&fakeDatabaseSearchRepo{err: errors.New("databases down")},
		&fakeAppSearchRepo{ids: []string{}},
	)
	ids := searcher.SearchResourceIDs(context.Background(), "api")
	if len(ids) != 1 || ids[0] != "srv-1" {
		t.Fatalf("expected surviving resolver results, got %v", ids)
	}
}
