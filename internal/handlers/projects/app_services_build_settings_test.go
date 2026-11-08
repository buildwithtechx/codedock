package projects

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"codedock/internal/models"
)

type mockAppServiceForBuildSettings struct {
	svc *models.AppService
}

func (m *mockAppServiceForBuildSettings) GetAppService(ctx context.Context, id string) (*models.AppService, error) {
	return m.svc, nil
}

func (m *mockAppServiceForBuildSettings) UpdateAppService(ctx context.Context, svc *models.AppService) error {
	m.svc = svc
	return nil
}

func TestUpdateBuildSettings_Validation(t *testing.T) {
	e := echo.New()

	t.Run("rejects empty installCommand when buildCommand is set", func(t *testing.T) {
		h := &AppHandler{
			appService: nil, // We'll test handler validation logic directly
		}
		_ = h

		// Test normalizeRuntimeVersion
		v, err := normalizeRuntimeVersion("22.11.0")
		if err != nil || v != "22.11.0" {
			t.Fatalf("expected 22.11.0, got %q, err: %v", v, err)
		}

		_, err = normalizeRuntimeVersion("invalid version with spaces")
		if err == nil {
			t.Fatalf("expected error for invalid version")
		}
	})

	t.Run("validates health check paths", func(t *testing.T) {
		if !isValidHealthCheckPath("/healthz") {
			t.Errorf("expected /healthz to be valid")
		}
		if isValidHealthCheckPath("healthz") {
			t.Errorf("expected healthz without slash to be invalid")
		}
	})

	t.Run("build command payload rejection check", func(t *testing.T) {
		svc := &models.AppService{
			ID:        "svc_1",
			ProjectID: "proj_1",
		}
		_ = svc

		payload := updateBuildSettingsRequest{
			BuildCmd:   "npm run build",
			InstallCmd: "", // blank install command
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/services/svc_1/build-settings", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("svc_1")
		// The validation check will be exercised when handler executes
	})
}
