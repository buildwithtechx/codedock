package system

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/services/migrations"
	"codedock.run/codedock/internal/utils"
)

type MigrationLifecycleHandler struct {
	migrations *migrations.Service
}

func NewMigrationLifecycleHandler(service *migrations.Service) *MigrationLifecycleHandler {
	return &MigrationLifecycleHandler{migrations: service}
}

func lifecycleClaims(c echo.Context) (*models.UserClaims, error) {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return nil, utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	return userClaims, nil
}

func (h *MigrationLifecycleHandler) ListSources(c echo.Context) error {
	claims, err := lifecycleClaims(c)
	if err != nil {
		return err
	}
	sources, err := h.migrations.ListSources(c.Request().Context(), claims.UserID, c.Param("id"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	if sources == nil {
		sources = []*models.MigrationSource{}
	}
	return utils.Success(c, "Migration sources retrieved", sources)
}

func (h *MigrationLifecycleHandler) TestSource(c echo.Context) error {
	claims, err := lifecycleClaims(c)
	if err != nil {
		return err
	}
	var req models.CreateMigrationSourceRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	fingerprint, err := h.migrations.TestSource(c.Request().Context(), claims.UserID, c.Param("id"), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Source verified", map[string]string{"fingerprint": fingerprint})
}

func (h *MigrationLifecycleHandler) CreateSource(c echo.Context) error {
	claims, err := lifecycleClaims(c)
	if err != nil {
		return err
	}
	var req models.CreateMigrationSourceRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	source, err := h.migrations.CreateSource(c.Request().Context(), claims.UserID, c.Param("id"), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Migration source connected", source)
}

func (h *MigrationLifecycleHandler) DeleteSource(c echo.Context) error {
	claims, err := lifecycleClaims(c)
	if err != nil {
		return err
	}
	if err := h.migrations.DeleteSource(c.Request().Context(), claims.UserID, c.Param("id"), c.Param("sourceId")); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Migration source removed", nil)
}

func (h *MigrationLifecycleHandler) Scan(c echo.Context) error {
	claims, err := lifecycleClaims(c)
	if err != nil {
		return err
	}
	var req models.ScanMigrationRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	stack, err := h.migrations.ScanSource(c.Request().Context(), claims.UserID, c.Param("id"), req.SourceID, nil)
	if err != nil {
		return utils.Error(c, http.StatusBadGateway, err.Error())
	}
	return utils.Success(c, "Scan complete", stack)
}

func (h *MigrationLifecycleHandler) ScanStream(c echo.Context) error {
	claims, err := lifecycleClaims(c)
	if err != nil {
		return err
	}
	sourceID := c.QueryParam("sourceId")
	if sourceID == "" {
		return utils.Error(c, http.StatusBadRequest, "sourceId is required")
	}
	writer, flush, ok := startEventStream(c)
	if !ok {
		return utils.Error(c, http.StatusInternalServerError, "streaming unsupported")
	}
	stack, err := h.migrations.ScanSource(c.Request().Context(), claims.UserID, c.Param("id"), sourceID, func(step, detail string) {
		writeStreamEvent(writer, flush, step, detail)
	})
	if err != nil {
		writeStreamEvent(writer, flush, "ERROR", err.Error())
		return nil
	}
	encoded, _ := json.Marshal(stack)
	writeStreamEvent(writer, flush, "RESULT", string(encoded))
	return nil
}

func (h *MigrationLifecycleHandler) RevealEnv(c echo.Context) error {
	claims, err := lifecycleClaims(c)
	if err != nil {
		return err
	}
	var req models.RevealMigrationEnvRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	env, err := h.migrations.RevealSourceEnv(c.Request().Context(), claims.UserID, c.Param("id"), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Environment revealed", env)
}

func (h *MigrationLifecycleHandler) Adopt(c echo.Context) error {
	claims, err := lifecycleClaims(c)
	if err != nil {
		return err
	}
	var req models.AdoptMigrationRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	result, err := h.migrations.AdoptSource(c.Request().Context(), claims.UserID, c.Param("id"), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Services adopted", result)
}

func (h *MigrationLifecycleHandler) Reimport(c echo.Context) error {
	claims, err := lifecycleClaims(c)
	if err != nil {
		return err
	}
	var req models.ReimportMigrationRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	result, err := h.migrations.ReimportSource(c.Request().Context(), claims.UserID, c.Param("id"), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Project reimported", result)
}

func (h *MigrationLifecycleHandler) RepoCompose(c echo.Context) error {
	claims, err := lifecycleClaims(c)
	if err != nil {
		return err
	}
	var req models.RepoComposeMigrationRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	stack, err := h.migrations.RepoCompose(c.Request().Context(), claims.UserID, c.Param("id"), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Repository compose read", stack)
}

func (h *MigrationLifecycleHandler) Preview(c echo.Context) error {
	claims, err := lifecycleClaims(c)
	if err != nil {
		return err
	}
	var req models.PreviewMigrationRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	preview, err := h.migrations.PreviewMigration(c.Request().Context(), claims.UserID, c.Param("id"), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Migration preview ready", preview)
}

func startEventStream(c echo.Context) (http.ResponseWriter, http.Flusher, bool) {
	writer := c.Response().Writer
	flush, ok := writer.(http.Flusher)
	if !ok {
		return nil, nil, false
	}
	header := writer.Header()
	header.Set("Content-Type", "text/event-stream")
	header.Set("Cache-Control", "no-cache")
	header.Set("Connection", "keep-alive")
	return writer, flush, true
}

func writeStreamEvent(writer http.ResponseWriter, flush http.Flusher, event, data string) {
	fmt.Fprintf(writer, "event: %s\ndata: %s\n\n", event, data)
	flush.Flush()
}
