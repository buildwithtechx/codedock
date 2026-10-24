package system

import (
	"net/http"

	authservices "codedock/internal/services/auth"
	systemservices "codedock/internal/services/system"

	"github.com/labstack/echo/v4"

	"codedock/internal/models"

	"codedock/internal/utils"
)

type DNSHandler struct {
	dnsService   *systemservices.DNSService
	auditService *authservices.AuditService
}

func NewDNSHandler(dnsService *systemservices.DNSService, auditService *authservices.AuditService) *DNSHandler {
	return &DNSHandler{dnsService: dnsService, auditService: auditService}
}

func (h *DNSHandler) Create(c echo.Context) error {
	var req models.CreateDNSRecordRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	record, err := h.dnsService.CreateRecord(c.Request().Context(), &req)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	h.auditService.LogAction(c.Request().Context(), authservices.AuditActionOpts{
		UserID:    utils.AuditActorID(c),
		Action:    "dns.create",
		Resource:  record.ID,
		IPAddress: c.RealIP(),
		Details:   map[string]string{"domain": req.DomainName, "record": req.RecordName, "type": req.RecordType},
	})
	return utils.Created(c, "DNS record created successfully", record)
}

func (h *DNSHandler) List(c echo.Context) error {
	domain := c.QueryParam("domain")

	records, err := h.dnsService.ListByDomain(c.Request().Context(), domain)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "DNS records fetched successfully", records)
}

func (h *DNSHandler) Update(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing record id")
	}

	var req models.UpdateDNSRecordRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	record, err := h.dnsService.UpdateRecord(c.Request().Context(), id, &req)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "DNS record updated successfully", record)
}

func (h *DNSHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing record id")
	}

	if err := h.dnsService.DeleteRecord(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	h.auditService.LogAction(c.Request().Context(), authservices.AuditActionOpts{
		UserID:    utils.AuditActorID(c),
		Action:    "dns.delete",
		Resource:  id,
		IPAddress: c.RealIP(),
	})
	return c.NoContent(http.StatusNoContent)
}
