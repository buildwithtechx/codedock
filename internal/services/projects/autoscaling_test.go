package projects

import (
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
	"context"
	"errors"
	"testing"
)

type autoscalingTestPolicies struct{ failure error }

func (s autoscalingTestPolicies) Get(context.Context, string) (*models.AutoscalingPolicy, error) {
	return models.DefaultAutoscalingPolicy("app"), s.failure
}
func (s autoscalingTestPolicies) Save(context.Context, *models.AutoscalingPolicy) error {
	return s.failure
}

type autoscalingTestApps struct{ failure error }

func (s autoscalingTestApps) GetByID(context.Context, string) (*models.AppService, error) {
	if s.failure != nil {
		return nil, s.failure
	}
	return &models.AppService{ID: "app", ProjectID: "project"}, nil
}

type autoscalingTestProjects struct{ remote bool }

func (s autoscalingTestProjects) Get(context.Context, string) (*models.ProjectConfig, error) {
	p := &models.ProjectConfig{ID: "project"}
	if s.remote {
		p.ServerID = "remote"
	}
	return p, nil
}
func TestAutoscalingSupportAndValidationErrors(t *testing.T) {
	for _, remote := range []bool{false, true} {
		service := NewAutoscalingService(autoscalingTestPolicies{}, autoscalingTestApps{}, autoscalingTestProjects{remote})
		p, err := service.Get(context.Background(), "app")
		if err != nil || p.Supported == remote {
			t.Fatalf("wrong target support: %v", err)
		}
		p.Enabled = true
		err = service.Save(context.Background(), p)
		if remote && !utils.IsValidation(err) {
			t.Fatal("unsupported target is not a validation error")
		}
		if !remote && err != nil {
			t.Fatal(err)
		}
	}
	service := NewAutoscalingService(autoscalingTestPolicies{errors.New("storage failed")}, autoscalingTestApps{}, autoscalingTestProjects{})
	if err := service.Save(context.Background(), models.DefaultAutoscalingPolicy("app")); err == nil || utils.IsValidation(err) {
		t.Fatal("storage error classified as invalid policy")
	}
	p := models.DefaultAutoscalingPolicy("app")
	p.MinReplicas = 0
	if err := service.Save(context.Background(), p); !utils.IsValidation(err) {
		t.Fatal("invalid limits classified as storage failure")
	}
}

func TestAutoscalingPreservesMissingService(t *testing.T) {
	service := NewAutoscalingService(autoscalingTestPolicies{}, autoscalingTestApps{utils.NewNotFoundError("Service", "app")}, autoscalingTestProjects{})
	if _, err := service.Get(context.Background(), "app"); !utils.IsNotFound(err) {
		t.Fatal("missing service lost not-found classification")
	}
	if err := service.Save(context.Background(), models.DefaultAutoscalingPolicy("app")); !utils.IsNotFound(err) {
		t.Fatal("missing service save lost not-found classification")
	}
}
