package auth

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"

	"codedock/internal/models"
	"codedock/internal/repositories"
	"github.com/google/uuid"
)

type AuditService struct {
	repo     repositories.AuditLogRepository
	searcher AuditResourceSearcher
}

func NewAuditService(repo repositories.AuditLogRepository, searcher AuditResourceSearcher) *AuditService {
	return &AuditService{repo: repo, searcher: searcher}
}

type AuditActionOpts struct {
	UserID    string
	Action    string
	Resource  string
	IPAddress string
	Details   any
}

func (s *AuditService) LogAction(ctx context.Context, opts AuditActionOpts) {
	var detailsStr string
	if opts.Details != nil {
		b, err := json.Marshal(opts.Details)
		if err == nil {
			detailsStr = string(b)
		}
	}

	log := &models.AuditLog{
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

func (s *AuditService) ListLogs(ctx context.Context, limit, offset int) ([]models.AuditLog, error) {
	return s.ListLogsByCategory(ctx, "", "", limit, offset)
}

func (s *AuditService) ListLogsByCategory(ctx context.Context, category string, q string, limit, offset int) ([]models.AuditLog, error) {
	logs, err := s.repo.ListFiltered(ctx, s.searchFilter(ctx, models.CategoryPrefixes(category), q), limit, offset)
	if err != nil {
		return nil, err
	}
	for i := range logs {
		logs[i].Category = models.CategoryForAction(logs[i].Action)
	}
	return logs, nil
}

func (s *AuditService) searchFilter(ctx context.Context, prefixes []string, q string) models.AuditLogFilter {
	filter := models.AuditLogFilter{Prefixes: prefixes}
	if term := strings.TrimSpace(q); term != "" {
		filter.Q = term
		if s.searcher != nil {
			filter.QResourceIDs = s.searcher.SearchResourceIDs(ctx, term)
		}
	}
	return filter
}

func (s *AuditService) Facets(ctx context.Context, category string, q string) (models.AuditFacets, error) {
	counts, total, err := s.repo.Facets(ctx, s.searchFilter(ctx, models.CategoryPrefixes(category), q))
	if err != nil {
		return models.AuditFacets{}, err
	}
	categories := make([]models.AuditCategoryCount, 0, len(models.AuditCategories()))
	for _, category := range models.AuditCategories() {
		categories = append(categories, models.AuditCategoryCount{
			AuditCategory: category,
			Count:         counts[category.ID],
		})
	}
	actors, _ := s.repo.ListActors(ctx)
	return models.AuditFacets{Total: total, Categories: categories, Actors: actors}, nil
}

func (s *AuditService) GetSettings(ctx context.Context, orgID string) (*models.AuditSettings, error) {
	return s.repo.GetSettings(ctx, orgID)
}

func (s *AuditService) UpdateSettings(ctx context.Context, settings *models.AuditSettings) error {
	return s.repo.UpdateSettings(ctx, settings)
}
