package projects

import (
	"context"
	"fmt"

	"codedock/internal/models"
	"codedock/internal/modules/ai"
	"codedock/internal/repositories"
)

type AIAnalysisService struct {
	deploymentRepo repositories.DeploymentRepository
	appRepo        repositories.AppServiceRepository
	aiService      *ai.AnalysisService
}

type AIExplanationResponse = ai.AIExplanationResponse

func NewAIAnalysisService(d repositories.DeploymentRepository, a repositories.AppServiceRepository, aiRepo repositories.AISettingsRepository) *AIAnalysisService {
	return &AIAnalysisService{
		deploymentRepo: d,
		appRepo:        a,
		aiService:      ai.NewAnalysisService(&aiSettingsRepoAdapter{repo: aiRepo}),
	}
}

func (s *AIAnalysisService) ExplainDeploymentFailure(ctx context.Context, deploymentID string) (*ai.AIExplanationResponse, error) {
	deployment, err := s.deploymentRepo.GetByID(ctx, deploymentID)
	if err != nil {
		return nil, fmt.Errorf("deployment not found: %w", err)
	}
	if deployment.Status != models.DeploymentStatusFailed {
		return nil, ai.ErrDeploymentNotFailed
	}

	app, err := s.appRepo.GetByID(ctx, deployment.ServiceID)
	if err != nil {
		return nil, fmt.Errorf("service not found: %w", err)
	}

	dCtx := &ai.DeploymentContext{
		ServiceName:    app.Name,
		RuntimeMode:    string(app.RuntimeMode),
		RepositoryURL:  app.RepositoryURL,
		Branch:         app.Branch,
		RootDirectory:  app.RootDirectory,
		InstallCommand: app.InstallCommand,
		BuildCommand:   app.BuildCommand,
		StartCommand:   app.StartCommand,
		Status:         string(deployment.Status),
		BuildLogs:      deployment.BuildLogs,
	}

	return s.aiService.ExplainDeploymentFailure(ctx, dCtx)
}
