package projects

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/utils"

	"codedock.run/codedock/internal/models"
	authservices "codedock.run/codedock/internal/services/auth"
	projectservices "codedock.run/codedock/internal/services/projects"
)

type ServiceVarHandler struct {
	appService   *projectservices.AppService
	auditService *authservices.AuditService
	envSugg      *projectservices.EnvSuggestionService
}

func NewServiceVarHandler(s *projectservices.AppService, audit *authservices.AuditService, envSugg *projectservices.EnvSuggestionService) *ServiceVarHandler {
	return &ServiceVarHandler{appService: s, auditService: audit, envSugg: envSugg}
}

func (h *ServiceVarHandler) List(c echo.Context) error {
	serviceID := c.Param("serviceId")
	list, err := h.appService.ListVariablesByService(c.Request().Context(), serviceID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", list)
}

func (h *ServiceVarHandler) Create(c echo.Context) error {
	serviceID := c.Param("serviceId")
	var req models.Variable
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	svc, err := h.appService.GetAppService(c.Request().Context(), serviceID)
	if err != nil || svc == nil {
		return utils.Error(c, http.StatusNotFound, "service not found")
	}
	req.ServiceID = serviceID
	req.ProjectID = svc.ProjectID
	req.EnvironmentID = svc.EnvironmentID
	created, err := h.appService.CreateVariable(c.Request().Context(), &req)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	h.auditService.LogAction(c.Request().Context(), authservices.AuditActionOpts{
		UserID:    "system",
		Action:    "service_var.create",
		Resource:  serviceID,
		IPAddress: c.RealIP(),
		Details: map[string]string{
			"variableId": created.ID,
			"key":        created.Key,
		},
	})

	return utils.Created(c, "Created successfully", created)
}

func (h *ServiceVarHandler) Update(c echo.Context) error {
	serviceID := c.Param("serviceId")
	id := c.Param("id")
	var req models.Variable
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if req.ServiceID != "" && req.ServiceID != serviceID {
		return utils.Error(c, http.StatusBadRequest, "service ID mismatch")
	}
	existing, err := h.appService.GetVariable(c.Request().Context(), id)
	if err != nil || existing == nil {
		return utils.Error(c, http.StatusNotFound, "variable not found")
	}
	if existing.ServiceID != serviceID {
		return utils.Error(c, http.StatusBadRequest, "service ID mismatch")
	}
	req.ID = id
	req.ServiceID = serviceID
	if err := h.appService.UpdateVariable(c.Request().Context(), &req); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	h.auditService.LogAction(c.Request().Context(), authservices.AuditActionOpts{
		UserID:    "system",
		Action:    "service_var.update",
		Resource:  serviceID,
		IPAddress: c.RealIP(),
		Details: map[string]string{
			"variableId": id,
			"key":        req.Key,
		},
	})

	return utils.Success(c, "Operation successful", req)
}

func (h *ServiceVarHandler) Delete(c echo.Context) error {
	serviceID := c.Param("serviceId")
	id := c.Param("id")
	existing, err := h.appService.GetVariable(c.Request().Context(), id)
	if err != nil || existing == nil {
		return utils.Error(c, http.StatusNotFound, "variable not found")
	}
	if existing.ServiceID != serviceID {
		return utils.Error(c, http.StatusBadRequest, "service ID mismatch")
	}
	if err := h.appService.DeleteVariable(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	h.auditService.LogAction(c.Request().Context(), authservices.AuditActionOpts{
		UserID:    "system",
		Action:    "service_var.delete",
		Resource:  serviceID,
		IPAddress: c.RealIP(),
		Details: map[string]string{
			"variableId": id,
		},
	})

	return utils.Success(c, "Deleted successfully", nil)
}

func (h *ServiceVarHandler) Suggest(c echo.Context) error {
	serviceID := c.Param("serviceId")
	svc, err := h.appService.GetAppService(c.Request().Context(), serviceID)
	if err != nil || svc == nil {
		return utils.Error(c, http.StatusNotFound, "service not found")
	}

	if svc.RepositoryURL == "" {
		return utils.Success(c, "No repository attached", []any{})
	}

	suggestions, err := h.envSugg.SuggestEnvVars(c.Request().Context(), svc.RepositoryURL, svc.Branch, svc.RootDirectory)
	if err != nil {
		return utils.Success(c, "No suggestions available", []any{})
	}

	return utils.Success(c, "Suggestions loaded", suggestions)
}
