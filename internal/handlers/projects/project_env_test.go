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
	projectservices "codedock/internal/services/projects"
)

type stubEnvRepo struct {
	vars map[string]string
}

func (s *stubEnvRepo) GetVars(_ context.Context, _ string) (map[string]string, error) {
	return s.vars, nil
}

func (s *stubEnvRepo) SetVar(_ context.Context, _, key, value string) error {
	if s.vars == nil {
		s.vars = make(map[string]string)
	}
	s.vars[key] = value
	return nil
}

func (s *stubEnvRepo) DeleteVar(_ context.Context, _, key string) error {
	delete(s.vars, key)
	return nil
}

func (s *stubEnvRepo) ListByProject(_ context.Context, _ string) ([]models.EnvironmentConfig, error) {
	return nil, nil
}

func (s *stubEnvRepo) Create(_ context.Context, _ *models.EnvironmentConfig) error {
	return nil
}

func (s *stubEnvRepo) Get(_ context.Context, _ string) (*models.EnvironmentConfig, error) {
	return nil, nil
}

func (s *stubEnvRepo) Delete(_ context.Context, _ string) error {
	return nil
}

func TestMigrateCommands(t *testing.T) {
	e := echo.New()
	stub := &stubEnvRepo{vars: make(map[string]string)}
	envService := projectservices.NewEnvironmentService(stub, nil, stub, nil)
	handler := NewProjectEnvHandler(envService, nil)

	t.Run("rejects non-migratable shell commands", func(t *testing.T) {
		body, _ := json.Marshal(MigrateCommandEnvironmentRequest{
			Commands: map[string]string{
				"startCommand": "PORT=$DANGEROUS start",
			},
		})
		req := httptest.NewRequest(http.MethodPost, "/projects/p-1/env/migrate-commands", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("p-1")

		_ = handler.MigrateCommands(c)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for unmigratable command, got %d", rec.Code)
		}
	})

	t.Run("rejects conflicting variable values across commands", func(t *testing.T) {
		body, _ := json.Marshal(MigrateCommandEnvironmentRequest{
			Commands: map[string]string{
				"buildCommand": "NODE_ENV=development npm run build",
				"startCommand": "NODE_ENV=production npm start",
			},
		})
		req := httptest.NewRequest(http.MethodPost, "/projects/p-1/env/migrate-commands", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("p-1")

		_ = handler.MigrateCommands(c)
		if rec.Code != http.StatusConflict {
			t.Fatalf("expected 409 for conflicting values, got %d", rec.Code)
		}
	})

	t.Run("successfully extracts variables and cleans commands", func(t *testing.T) {
		body, _ := json.Marshal(MigrateCommandEnvironmentRequest{
			Commands: map[string]string{
				"buildCommand": "API_URL=https://api.test npm run build",
				"startCommand": "API_URL=https://api.test PORT=3000 node server.js",
			},
		})
		req := httptest.NewRequest(http.MethodPost, "/projects/p-1/env/migrate-commands", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("p-1")

		err := handler.MigrateCommands(c)
		if err != nil {
			t.Fatalf("unexpected handler error: %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
		}
		if stub.vars["API_URL"] != "https://api.test" || stub.vars["PORT"] != "3000" {
			t.Fatalf("variables were not saved: %+v", stub.vars)
		}
	})
}
