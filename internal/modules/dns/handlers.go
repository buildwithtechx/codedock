package dns

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	authservices "codedock/internal/services/auth"
	"codedock/internal/utils"
)

type Handler struct {
	dnsService   *Service
	auditService *authservices.AuditService
}

func NewHandler(dnsService *Service, auditService *authservices.AuditService) *Handler {
	return &Handler{dnsService: dnsService, auditService: auditService}
}

func (h *Handler) Create(c echo.Context) error {
	var req CreateDNSRecordRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	record, err := h.dnsService.CreateRecord(c.Request().Context(), &req)
	if err != nil {
		if errors.Is(err, ErrFieldsRequired) {
			return utils.Error(c, http.StatusBadRequest, err.Error())
		}
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	if h.auditService != nil {
		h.auditService.LogAction(c.Request().Context(), authservices.AuditActionOpts{
			UserID:    utils.AuditActorID(c),
			Action:    "dns.create",
			Resource:  record.ID,
			IPAddress: c.RealIP(),
			Details:   map[string]string{"domain": req.DomainName, "record": req.RecordName, "type": req.RecordType},
		})
	}
	return utils.Created(c, "DNS record created successfully", record)
}

func (h *Handler) List(c echo.Context) error {
	domain := c.QueryParam("domain")

	records, err := h.dnsService.ListByDomain(c.Request().Context(), domain)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "DNS records fetched successfully", records)
}

func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing record id")
	}

	var req UpdateDNSRecordRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	record, err := h.dnsService.UpdateRecord(c.Request().Context(), id, &req)
	if err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			return utils.Error(c, http.StatusNotFound, err.Error())
		}
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "DNS record updated successfully", record)
}

func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing record id")
	}

	if err := h.dnsService.DeleteRecord(c.Request().Context(), id); err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			return utils.Error(c, http.StatusNotFound, err.Error())
		}
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	if h.auditService != nil {
		h.auditService.LogAction(c.Request().Context(), authservices.AuditActionOpts{
			UserID:    utils.AuditActorID(c),
			Action:    "dns.delete",
			Resource:  id,
			IPAddress: c.RealIP(),
		})
	}
	return c.NoContent(http.StatusNoContent)
}
