package databases

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/http/middleware"
	"codedock/internal/utils"
)

type Handler struct {
	service *Service
}

type DatabaseHandler = Handler

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func NewDatabaseHandler(service *Service) *Handler {
	return NewHandler(service)
}

func (h *Handler) verifyProjectOwnership(c echo.Context, projectID string) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	err := h.service.VerifyProjectOwnership(c.Request().Context(), projectID, user.UserID, string(user.Role))
	if err != nil {
		if errors.Is(err, ErrProjectNotFound) {
			return utils.Error(c, http.StatusNotFound, "project not found")
		}
		return utils.Error(c, http.StatusForbidden, "access denied")
	}
	return nil
}

func (h *Handler) ListDatabases(c echo.Context) error {
	databases, err := h.service.ListDatabases(c.Request().Context(), c.QueryParam("search"))
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if databases == nil {
		databases = []*Database{}
	}
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user != nil && string(user.Role) != "admin" && string(user.Role) != "owner" {
		filtered := make([]*Database, 0)
		for _, db := range databases {
			if h.service.IsProjectMemberOrOwner(c.Request().Context(), db.ProjectID, user.UserID, string(user.Role)) {
				filtered = append(filtered, db)
			}
		}
		return utils.Success(c, "Operation successful", filtered)
	}
	return utils.Success(c, "Operation successful", databases)
}

func (h *Handler) CreateDatabase(c echo.Context) error {
	var req CreateDatabaseRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if err := h.verifyProjectOwnership(c, req.ProjectID); err != nil {
		return err
	}
	db, err := h.service.CreateDatabaseFromRequest(c.Request().Context(), &req)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", db)
}

func (h *Handler) GetDatabase(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing database id parameter")
	}
	db, err := h.service.GetDatabase(c.Request().Context(), id)
	if err != nil || db == nil {
		return utils.Error(c, http.StatusNotFound, "database not found")
	}
	if err := h.verifyProjectOwnership(c, db.ProjectID); err != nil {
		return err
	}
	return utils.Success(c, "Operation successful", db)
}

func (h *Handler) UpdateDatabase(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing database id parameter")
	}
	db, err := h.service.GetDatabase(c.Request().Context(), id)
	if err != nil || db == nil {
		return utils.Error(c, http.StatusNotFound, "database not found")
	}
	if err := h.verifyProjectOwnership(c, db.ProjectID); err != nil {
		return err
	}
	var req UpdateDatabaseRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	db.ExternalDNS = req.ExternalDNS
	db.CustomArgs = req.CustomArgs
	db.LogicalReplication = req.LogicalReplication
	db.CPULimit = req.CPULimit
	db.MemoryLimit = req.MemoryLimit
	if err := h.service.UpdateDatabase(c.Request().Context(), db); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", db)
}

func (h *Handler) DeleteDatabase(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing database id parameter")
	}
	db, err := h.service.GetDatabase(c.Request().Context(), id)
	if err != nil || db == nil {
		return utils.Error(c, http.StatusNotFound, "database not found")
	}
	if err := h.verifyProjectOwnership(c, db.ProjectID); err != nil {
		return err
	}
	if err := h.service.DeleteDatabase(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", map[string]string{"status": "deleted"})
}

func (h *Handler) StartDatabase(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing database id parameter")
	}
	db, err := h.service.GetDatabase(c.Request().Context(), id)
	if err != nil || db == nil {
		return utils.Error(c, http.StatusNotFound, "database not found")
	}
	if err := h.verifyProjectOwnership(c, db.ProjectID); err != nil {
		return err
	}
	dbStarted, err := h.service.StartDatabase(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", dbStarted)
}

func (h *Handler) StopDatabase(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing database id parameter")
	}
	db, err := h.service.GetDatabase(c.Request().Context(), id)
	if err != nil || db == nil {
		return utils.Error(c, http.StatusNotFound, "database not found")
	}
	if err := h.verifyProjectOwnership(c, db.ProjectID); err != nil {
		return err
	}
	if err := h.service.StopDatabase(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", map[string]string{"status": "stopped"})
}

func (h *Handler) RestartDatabase(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing database id parameter")
	}
	db, err := h.service.GetDatabase(c.Request().Context(), id)
	if err != nil || db == nil {
		return utils.Error(c, http.StatusNotFound, "database not found")
	}
	if err := h.verifyProjectOwnership(c, db.ProjectID); err != nil {
		return err
	}
	if err := h.service.StopDatabase(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	dbStarted, err := h.service.StartDatabase(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", dbStarted)
}

func (h *Handler) ImportData(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing database id parameter")
	}
	db, err := h.service.GetDatabase(c.Request().Context(), id)
	if err != nil || db == nil {
		return utils.Error(c, http.StatusNotFound, "database not found")
	}
	if err := h.verifyProjectOwnership(c, db.ProjectID); err != nil {
		return err
	}
	var req ImportDatabaseRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	go func() {
		if err := h.service.ImportData(context.Background(), id, req.SourceURL); err != nil {
			slog.Error("failed to import data", "error", err)
		}
	}()
	return utils.Success(c, "Import started", nil)
}

func (h *Handler) QueryDatabase(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing database id parameter")
	}
	db, err := h.service.GetDatabase(c.Request().Context(), id)
	if err != nil || db == nil {
		return utils.Error(c, http.StatusNotFound, "database not found")
	}
	if err := h.verifyProjectOwnership(c, db.ProjectID); err != nil {
		return err
	}
	var req DatabaseQueryRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	res, err := h.service.QueryDatabase(c.Request().Context(), id, req.Query, db.ProjectID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	userID := ""
	if user := middleware.GetUserClaimsFromContext(c.Request().Context()); user != nil {
		userID = user.UserID
	}
	h.service.LogAction(c.Request().Context(), AuditLogEntry{
		UserID:    userID,
		Action:    "database.query",
		Resource:  id,
		IPAddress: c.RealIP(),
		Details:   map[string]any{"query": req.Query, "database_id": id},
	})
	return utils.Success(c, "Operation successful", res)
}

func (h *Handler) RevealCredentials(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing database id parameter")
	}
	db, err := h.service.GetDatabase(c.Request().Context(), id)
	if err != nil || db == nil {
		return utils.Error(c, http.StatusNotFound, "database not found")
	}
	if err := h.verifyProjectOwnership(c, db.ProjectID); err != nil {
		return err
	}
	userID := ""
	if user := middleware.GetUserClaimsFromContext(c.Request().Context()); user != nil {
		userID = user.UserID
	}
	h.service.LogAction(c.Request().Context(), AuditLogEntry{
		UserID:    userID,
		Action:    "database.credentials_reveal",
		Resource:  id,
		IPAddress: c.RealIP(),
		Details:   map[string]any{"database_id": id},
	})
	return utils.Success(c, "Credentials retrieved successfully", map[string]string{
		"password": db.Password,
	})
}
