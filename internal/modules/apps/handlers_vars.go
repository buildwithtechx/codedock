package apps

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type VariableHandler struct {
	service *Service
	access  AccessChecker
}

func NewVariableHandler(s *Service, ac AccessChecker) *VariableHandler {
	return &VariableHandler{service: s, access: ac}
}

func (h *VariableHandler) List(c echo.Context) error {
	serviceID := c.Param("id")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing service id")
	}

	vars, err := h.service.ListVariablesByService(c.Request().Context(), serviceID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", vars)
}

func (h *VariableHandler) Create(c echo.Context) error {
	serviceID := c.Param("id")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing service id")
	}

	var req CreateVariableRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if req.Key == "" {
		return utils.Error(c, http.StatusBadRequest, "key is required")
	}

	v := &Variable{
		ServiceID: serviceID,
		Key:       req.Key,
		Value:     req.Value,
		IsSecret:  req.IsSecret,
	}
	created, err := h.service.CreateVariable(c.Request().Context(), v)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", created)
}

func (h *VariableHandler) Update(c echo.Context) error {
	id := c.Param("varId")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing variable id")
	}

	existing, err := h.service.GetVariable(c.Request().Context(), id)
	if err != nil || existing == nil {
		return utils.Error(c, http.StatusNotFound, "variable not found")
	}

	var req UpdateVariableRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	if req.Key != nil {
		existing.Key = *req.Key
	}
	if req.Value != nil {
		existing.Value = *req.Value
	}
	if req.IsSecret != nil {
		existing.IsSecret = *req.IsSecret
	}

	if err := h.service.UpdateVariable(c.Request().Context(), existing); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Updated successfully", existing)
}

func (h *VariableHandler) Delete(c echo.Context) error {
	id := c.Param("varId")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing variable id")
	}

	if err := h.service.DeleteVariable(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
