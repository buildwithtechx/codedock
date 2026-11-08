package environments

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

type RouteRuleService struct {
	repo RouteRuleRepository
}

func NewRouteRuleService(repo RouteRuleRepository) *RouteRuleService {
	return &RouteRuleService{repo: repo}
}

func (s *RouteRuleService) List(ctx context.Context, serviceID string) ([]*RouteRule, error) {
	return s.repo.ListByService(ctx, serviceID)
}

func (s *RouteRuleService) Create(ctx context.Context, serviceID string, req CreateRouteRuleRequest) (*RouteRule, error) {
	if req.Name == "" {
		return nil, errors.New("rule name is required")
	}
	if req.RuleType == "" {
		return nil, errors.New("rule type is required")
	}
	if req.Spec == nil {
		return nil, errors.New("rule spec is required")
	}

	specBytes, err := json.Marshal(req.Spec)
	if err != nil {
		return nil, fmt.Errorf("marshal rule spec: %w", err)
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	rule := &RouteRule{
		ServiceID: serviceID,
		Name:      req.Name,
		Enabled:   enabled,
		RuleType:  req.RuleType,
		SpecJSON:  string(specBytes),
	}

	if err := s.repo.Create(ctx, rule); err != nil {
		return nil, err
	}
	return rule, nil
}

func (s *RouteRuleService) Update(ctx context.Context, id string, req UpdateRouteRuleRequest) (*RouteRule, error) {
	var specJSON *string
	if req.Spec != nil {
		bytes, err := json.Marshal(req.Spec)
		if err != nil {
			return nil, fmt.Errorf("marshal rule spec: %w", err)
		}
		str := string(bytes)
		specJSON = &str
	}

	if err := s.repo.Update(ctx, id, req.Name, req.Enabled, specJSON); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *RouteRuleService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *RouteRuleService) GetByID(ctx context.Context, id string) (*RouteRule, error) {
	return s.repo.GetByID(ctx, id)
}
