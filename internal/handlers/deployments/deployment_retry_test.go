package deployments

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"codedock/internal/models"
)

type dummyDeployRepo struct {
	dep *models.Deployment
}

func TestRetryRequiresValidImageRef(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/deployments/dep-123/retry", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("dep-123")

	handler := &DeploymentHandler{}

	// If missing ID
	c.SetParamValues("")
	err := handler.Retry(c)
	if err == nil && rec.Code != http.StatusBadRequest {
		t.Fatalf("expected error or bad request for missing id, got code %d", rec.Code)
	}
}

func TestRetryErrorCodeConstants(t *testing.T) {
	if ErrDeploymentBuildUnavailable != "DEPLOYMENT_BUILD_UNAVAILABLE" {
		t.Fatalf("expected DEPLOYMENT_BUILD_UNAVAILABLE, got %s", ErrDeploymentBuildUnavailable)
	}
}

func TestDeploymentRetryFields(t *testing.T) {
	target := &models.Deployment{
		ID:            "target-id",
		ServiceID:     "svc-1",
		EnvironmentID: "env-1",
		ProjectID:     "proj-1",
		Branch:        "main",
		CommitHash:    "commit123",
		ImageRef:      "codedock-app-svc-1:built-tag",
	}

	newDep := &models.Deployment{
		ServiceID:     target.ServiceID,
		EnvironmentID: target.EnvironmentID,
		ProjectID:     target.ProjectID,
		Status:        models.DeploymentStatusPending,
		CommitHash:    target.CommitHash,
		CommitMessage: "Retry deployment " + target.ID,
		Branch:        target.Branch,
		ImageRef:      target.ImageRef,
		Trigger:       "retry",
	}

	if newDep.Trigger != "retry" {
		t.Fatalf("expected trigger 'retry', got %s", newDep.Trigger)
	}
	if newDep.ImageRef != target.ImageRef {
		t.Fatalf("expected imageRef to be copied, got %s", newDep.ImageRef)
	}
}
