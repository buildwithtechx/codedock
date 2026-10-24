package projects

import (
	"codedock/internal/models"
	"codedock/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func (s *AppService) CreateApplicationSetup(ctx context.Context, app *models.AppService) (*models.AppService, error) {
	if app.ID == "" {
		return s.CreateAppService(ctx, app)
	}
	existing, err := s.repo.GetByID(ctx, app.ID)
	if err == nil {
		return matchingApplicationSetup(existing, app)
	}
	if !utils.IsNotFound(err) {
		return nil, err
	}
	created, createErr := s.CreateAppService(ctx, app)
	if createErr == nil {
		return created, nil
	}
	existing, err = s.repo.GetByID(ctx, app.ID)
	if err == nil {
		return matchingApplicationSetup(existing, app)
	}
	return nil, createErr
}

func matchingApplicationSetup(existing, requested *models.AppService) (*models.AppService, error) {
	normalize := func(app *models.AppService) ([]byte, error) {
		value := *app
		value.ID, value.GitUserID, value.ContainerID, value.Icon, value.DeployToken, value.AppID = "", "", "", "", "", ""
		value.Status = ""
		value.CreatedAt, value.UpdatedAt = time.Time{}, time.Time{}
		value.Volumes = nil
		if value.Replicas <= 0 {
			value.Replicas = 1
		}
		return json.Marshal(value)
	}
	left, err := normalize(existing)
	if err != nil {
		return nil, err
	}
	right, err := normalize(requested)
	if err != nil {
		return nil, err
	}
	if string(left) != string(right) {
		return nil, fmt.Errorf("application already exists with different configuration; reload it before continuing")
	}
	return existing, nil
}
