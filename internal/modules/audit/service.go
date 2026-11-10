package audit

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type auditResourceSearcher struct {
	projects  NameSearchable
	servers   NameSearchable
	databases NameSearchable
	apps      NameSearchable
}

func NewAuditResourceSearcher(
	projects NameSearchable,
	servers NameSearchable,
	databases NameSearchable,
	apps NameSearchable,
) AuditResourceSearcher {
	return &auditResourceSearcher{projects: projects, servers: servers, databases: databases, apps: apps}
}

func (s *auditResourceSearcher) SearchResourceIDs(ctx context.Context, term string) []string {
	resolvers := []func(context.Context) ([]string, error){
		func(ctx context.Context) ([]string, error) {
			if s.projects == nil {
				return nil, nil
			}
			return s.projects.SearchIDsByName(ctx, term, AuditSearchResolverLimit)
		},
		func(ctx context.Context) ([]string, error) {
			if s.servers == nil {
				return nil, nil
			}
			return s.servers.SearchIDsByName(ctx, term, AuditSearchResolverLimit)
		},
		func(ctx context.Context) ([]string, error) {
			if s.databases == nil {
				return nil, nil
			}
			return s.databases.SearchIDsByName(ctx, term, AuditSearchResolverLimit)
		},
		func(ctx context.Context) ([]string, error) {
			if s.apps == nil {
				return nil, nil
			}
			return s.apps.SearchIDsByName(ctx, term, AuditSearchResolverLimit)
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

type AuditService struct {
	repo     AuditLogRepository
	searcher AuditResourceSearcher
}

func NewAuditService(repo AuditLogRepository, searcher AuditResourceSearcher) *AuditService {
	return &AuditService{repo: repo, searcher: searcher}
}

func (s *AuditService) LogAction(ctx context.Context, opts AuditActionOpts) {
	if s == nil || s.repo == nil {
		return
	}
	var detailsStr string
	if opts.Details != nil {
		b, err := json.Marshal(opts.Details)
		if err == nil {
			detailsStr = string(b)
		}
	}

	log := &AuditLog{
		ID:        uuid.New().String(),
		UserID:    opts.UserID,
		Action:    opts.Action,
		Resource:  opts.Resource,
		Details:   detailsStr,
		IPAddress: opts.IPAddress,
	}

	go func() {
		err := s.repo.Create(context.Background(), log)
		if err != nil {
			slog.Error("failed to write audit log", "err", err, "action", opts.Action)
		}
	}()
}

func (s *AuditService) ListLogs(ctx context.Context, limit, offset int) ([]AuditLog, error) {
	return s.ListLogsByCategory(ctx, "", "", limit, offset)
}

func (s *AuditService) ListLogsByCategory(ctx context.Context, category string, q string, limit, offset int) ([]AuditLog, error) {
	logs, err := s.repo.ListFiltered(ctx, s.searchFilter(ctx, CategoryPrefixes(category), q), limit, offset)
	if err != nil {
		return nil, err
	}
	for i := range logs {
		logs[i].Category = CategoryForAction(logs[i].Action)
	}
	return logs, nil
}

func (s *AuditService) searchFilter(ctx context.Context, prefixes []string, q string) AuditLogFilter {
	filter := AuditLogFilter{Prefixes: prefixes}
	if term := strings.TrimSpace(q); term != "" {
		filter.Q = term
		if s.searcher != nil {
			filter.QResourceIDs = s.searcher.SearchResourceIDs(ctx, term)
		}
	}
	return filter
}

func (s *AuditService) Facets(ctx context.Context, category string, q string) (AuditFacets, error) {
	counts, total, err := s.repo.Facets(ctx, s.searchFilter(ctx, CategoryPrefixes(category), q))
	if err != nil {
		return AuditFacets{}, err
	}
	categories := make([]AuditCategoryCount, 0, len(AuditCategories()))
	for _, category := range AuditCategories() {
		categories = append(categories, AuditCategoryCount{
			AuditCategory: category,
			Count:         counts[category.ID],
		})
	}
	actors, _ := s.repo.ListActors(ctx)
	return AuditFacets{Total: total, Categories: categories, Actors: actors}, nil
}

func (s *AuditService) GetSettings(ctx context.Context, orgID string) (*AuditSettings, error) {
	return s.repo.GetSettings(ctx, orgID)
}

func (s *AuditService) UpdateSettings(ctx context.Context, settings *AuditSettings) error {
	return s.repo.UpdateSettings(ctx, settings)
}
