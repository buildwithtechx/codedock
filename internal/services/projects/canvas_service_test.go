package projects

import (
	"context"
	"testing"

	"codedock.run/codedock/internal/models"
)

type canvasServiceRepoStub struct {
	organizationID string
}

func (s *canvasServiceRepoStub) ListCanvasSummaries(_ context.Context, organizationID string) ([]models.CanvasSummary, error) {
	s.organizationID = organizationID
	return []models.CanvasSummary{}, nil
}

func (s *canvasServiceRepoStub) GetCanvasSummary(_ context.Context, _ string) (*models.CanvasSummary, error) {
	return nil, nil
}

func (s *canvasServiceRepoStub) GetEnvironmentCanvas(_ context.Context, _ string) (*models.EnvironmentCanvas, error) {
	return nil, nil
}

func TestCanvasServiceListSummariesRequiresOrganization(t *testing.T) {
	repo := &canvasServiceRepoStub{}
	service := NewCanvasService(repo)

	if _, err := service.ListSummaries(context.Background(), ""); err == nil {
		t.Fatal("expected organization validation error")
	}
	if repo.organizationID != "" {
		t.Fatalf("unexpected repository call for organization %q", repo.organizationID)
	}
}

func TestCanvasServiceListSummariesScopesRepositoryToOrganization(t *testing.T) {
	repo := &canvasServiceRepoStub{}
	service := NewCanvasService(repo)

	if _, err := service.ListSummaries(context.Background(), "organization-1"); err != nil {
		t.Fatalf("list summaries: %v", err)
	}
	if repo.organizationID != "organization-1" {
		t.Fatalf("repository organization = %q, want organization-1", repo.organizationID)
	}
}
