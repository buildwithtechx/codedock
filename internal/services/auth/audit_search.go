package auth

import (
	"context"
	"log/slog"
	"sync"

	"codedock/internal/repositories"
)

const auditSearchResolverLimit = 200

type AuditResourceSearcher interface {
	SearchResourceIDs(ctx context.Context, term string) []string
}

type auditResourceSearcher struct {
	projects  repositories.ProjectRepository
	servers   repositories.ServerRepository
	databases repositories.DatabaseRepository
	apps      repositories.AppServiceRepository
}

func NewAuditResourceSearcher(
	projects repositories.ProjectRepository,
	servers repositories.ServerRepository,
	databases repositories.DatabaseRepository,
	apps repositories.AppServiceRepository,
) AuditResourceSearcher {
	return &auditResourceSearcher{projects: projects, servers: servers, databases: databases, apps: apps}
}

func (s *auditResourceSearcher) SearchResourceIDs(ctx context.Context, term string) []string {
	resolvers := []func(context.Context) ([]string, error){
		func(ctx context.Context) ([]string, error) {
			return s.projects.SearchIDsByName(ctx, term, auditSearchResolverLimit)
		},
		func(ctx context.Context) ([]string, error) {
			return s.servers.SearchIDsByName(ctx, term, auditSearchResolverLimit)
		},
		func(ctx context.Context) ([]string, error) {
			return s.databases.SearchIDsByName(ctx, term, auditSearchResolverLimit)
		},
		func(ctx context.Context) ([]string, error) {
			return s.apps.SearchIDsByName(ctx, term, auditSearchResolverLimit)
		},
	}
	results := make([][]string, len(resolvers))
	var wg sync.WaitGroup
	for i, resolve := range resolvers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ids, err := resolve(ctx)
			if err != nil {
				slog.Error("audit resource search resolver failed", "index", i, "err", err)
				return
			}
			results[i] = ids
		}()
	}
	wg.Wait()
	seen := make(map[string]struct{})
	merged := make([]string, 0)
	for _, ids := range results {
		for _, id := range ids {
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			merged = append(merged, id)
		}
	}
	return merged
}
