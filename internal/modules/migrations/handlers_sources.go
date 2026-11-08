package migrations

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type LifecycleHandler struct {
	service LifecycleService
}

func NewLifecycleHandler(service LifecycleService) *LifecycleHandler {
	return &LifecycleHandler{service: service}
}

func (h *LifecycleHandler) ListSources(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	sources, err := h.service.ListSources(c.Request().Context(), userID, c.Param("id"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	if sources == nil {
		sources = []*MigrationSource{}
	}
	return utils.Success(c, "Migration sources retrieved", sources)
}

func (h *LifecycleHandler) TestSource(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req CreateMigrationSourceRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	fingerprint, err := h.service.TestSource(c.Request().Context(), userID, c.Param("id"), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Source verified", map[string]string{"fingerprint": fingerprint})
}

func (h *LifecycleHandler) CreateSource(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req CreateMigrationSourceRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	source, err := h.service.CreateSource(c.Request().Context(), userID, c.Param("id"), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Migration source connected", source)
}

func (h *LifecycleHandler) DeleteSource(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	if err := h.service.DeleteSource(c.Request().Context(), userID, c.Param("id"), c.Param("sourceId")); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Migration source removed", nil)
}

func (h *LifecycleHandler) Scan(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req ScanMigrationRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	stack, err := h.service.Scan(c.Request().Context(), userID, c.Param("id"), req.SourceID)
	if err != nil {
		return utils.Error(c, http.StatusBadGateway, err.Error())
	}
	return utils.Success(c, "Scan complete", stack)
}

func (h *LifecycleHandler) ScanStream(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	sourceID := c.QueryParam("sourceId")
	if sourceID == "" {
		return utils.Error(c, http.StatusBadRequest, "sourceId is required")
	}
	writer, flush, ok := startEventStream(c)
	if !ok {
		return utils.Error(c, http.StatusInternalServerError, "streaming unsupported")
	}
	err := h.service.ScanStream(c.Request().Context(), userID, c.Param("id"), sourceID, func(event string, data any) error {
		var encoded string
		if s, ok := data.(string); ok {
			encoded = s
		} else {
			b, _ := json.Marshal(data)
			encoded = string(b)
		}
		writeStreamEvent(writer, flush, event, encoded)
		return nil
	})
	if err != nil {
		writeStreamEvent(writer, flush, "ERROR", err.Error())
		return nil
	}
	return nil
}

func (h *LifecycleHandler) RevealEnv(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req RevealMigrationEnvRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	env, err := h.service.RevealEnv(c.Request().Context(), userID, c.Param("id"), req.SourceID, req.ContainerID)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Environment revealed", env)
}

func (h *LifecycleHandler) Adopt(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req AdoptMigrationRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	result, err := h.service.Adopt(c.Request().Context(), userID, c.Param("id"), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Services adopted", result)
}

func (h *LifecycleHandler) Reimport(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req ReimportMigrationRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	result, err := h.service.Reimport(c.Request().Context(), userID, c.Param("id"), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Project reimported", result)
}

func (h *LifecycleHandler) RepoCompose(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req RepoComposeMigrationRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	stack, err := h.service.RepoCompose(c.Request().Context(), userID, c.Param("id"), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Repository compose read", stack)
}

func (h *LifecycleHandler) Preview(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req PreviewMigrationRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	preview, err := h.service.Preview(c.Request().Context(), userID, c.Param("id"), req)
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
