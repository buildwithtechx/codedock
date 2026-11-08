package deployments

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (s *DeploymentService) Rollback(ctx context.Context, id, reason string) (*Deployment, error) {
	targetDep, err := s.GetDeployment(ctx, id)
	if err != nil || targetDep == nil {
		return nil, ErrDeploymentNotFound
	}

	rollbackImageRef := strings.TrimSpace(targetDep.ImageRef)
	commitMessage := fmt.Sprintf("Rollback to deployment %s", targetDep.ID)
	if reason != "" {
		commitMessage += ": " + reason
	}

	newDep := &Deployment{
		ID:            uuid.NewString(),
		ServiceID:     targetDep.ServiceID,
		EnvironmentID: targetDep.EnvironmentID,
		ProjectID:     targetDep.ProjectID,
		Status:        DeploymentStatusBuilding,
		CommitHash:    targetDep.CommitHash,
		CommitMessage: commitMessage,
		Branch:        targetDep.Branch,
		ImageRef:      rollbackImageRef,
		Trigger:       "Rollback",
		BuildLogs:     fmt.Sprintf("Rolling back to deployment %s...\n", targetDep.ID),
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	created, err := s.CreateDeployment(ctx, newDep)
	if err != nil {
		return nil, err
	}

	s.ExecuteDeploymentAsync(created)
	return created, nil
}

func (s *DeploymentService) Retry(ctx context.Context, id string) (*Deployment, error) {
	targetDep, err := s.GetDeployment(ctx, id)
	if err != nil || targetDep == nil {
		return nil, ErrDeploymentNotFound
	}

	imageRef := strings.TrimSpace(targetDep.ImageRef)
	newDep := &Deployment{
		ID:            uuid.NewString(),
		ServiceID:     targetDep.ServiceID,
		EnvironmentID: targetDep.EnvironmentID,
		ProjectID:     targetDep.ProjectID,
		Status:        DeploymentStatusPending,
		CommitHash:    targetDep.CommitHash,
		CommitMessage: "Retry deployment " + targetDep.ID,
		Branch:        targetDep.Branch,
		ImageRef:      imageRef,
		Trigger:       "retry",
		BuildLogs:     fmt.Sprintf("⚡ Retrying deployment using artifact %s...\n", imageRef),
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	created, err := s.CreateDeployment(ctx, newDep)
	if err != nil {
		return nil, err
	}

	s.ExecuteDeploymentAsync(created)
	return created, nil
}

func (s *DeploymentService) TogglePin(ctx context.Context, id string) (bool, error) {
	if id == "" {
		return false, ErrDeploymentNotFound
	}
	return s.repo.TogglePinned(ctx, id)
}

func (s *DeploymentService) Cancel(ctx context.Context, id string) error {
	return s.CancelDeployment(id)
}
