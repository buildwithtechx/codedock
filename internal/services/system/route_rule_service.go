package system

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
)

type RouteRuleService struct {
	repo repositories.RouteRuleRepository
}

func NewRouteRuleService(repo repositories.RouteRuleRepository) *RouteRuleService {
	return &RouteRuleService{repo: repo}
}

func (s *RouteRuleService) List(ctx context.Context, serviceID string) ([]*models.RouteRule, error) {
	return s.repo.ListByService(ctx, serviceID)
}

func (s *RouteRuleService) Create(ctx context.Context, serviceID string, req models.CreateRouteRuleRequest) (*models.RouteRule, error) {
	if err := validateRuleType(req.RuleType); err != nil {
		return nil, err
	}
	specJSON, err := json.Marshal(req.Spec)
	if err != nil {
		return nil, fmt.Errorf("marshal spec: %w", err)
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	rule := &models.RouteRule{
		ID:        uuid.NewString(),
		ServiceID: serviceID,
		Name:      req.Name,
		Enabled:   enabled,
		RuleType:  req.RuleType,
		SpecJSON:  string(specJSON),
	}
	if err := s.repo.Create(ctx, rule); err != nil {
		return nil, fmt.Errorf("create route rule: %w", err)
	}
	return rule, nil
}

func (s *RouteRuleService) Update(ctx context.Context, id string, req models.UpdateRouteRuleRequest) (*models.RouteRule, error) {
	var specJSON *string
	if req.Spec != nil {
		b, err := json.Marshal(req.Spec)
		if err != nil {
			return nil, fmt.Errorf("marshal spec: %w", err)
		}
		str := string(b)
		specJSON = &str
	}
	if err := s.repo.Update(ctx, id, req.Name, req.Enabled, specJSON); err != nil {
		return nil, fmt.Errorf("update route rule: %w", err)
	}
	return s.repo.GetByID(ctx, id)
}

func (s *RouteRuleService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *RouteRuleService) GetByID(ctx context.Context, id string) (*models.RouteRule, error) {
	return s.repo.GetByID(ctx, id)
}

func validateRuleType(rt models.RouteRuleType) error {
	switch rt {
	case models.RouteRuleTypeRateLimit, models.RouteRuleTypeIPAllowlist,
		models.RouteRuleTypeIPBlocklist, models.RouteRuleTypeHeaders:
		return nil
	}
	return fmt.Errorf("unknown rule type: %s", rt)
}
