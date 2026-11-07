package deployments

import (
	"github.com/labstack/echo/v4"
	"net/http/httptest"
	"testing"
)

func TestDeploymentPermissionHelpersReturnAnErrorWithoutClaims(t *testing.T) {
	context := echo.New().NewContext(httptest.NewRequest("POST", "/deployments/deployment/cancel", nil), httptest.NewRecorder())
	handler := &DeploymentHandler{}
	if err := handler.verifyProjectAdmin(context, "project"); err == nil {
		t.Fatal("admin denial returned nil and allowed the controller to continue")
	}
	if err := handler.verifyProjectOwnership(context, "project"); err == nil {
		t.Fatal("ownership denial returned nil and allowed the controller to continue")
	}
}
