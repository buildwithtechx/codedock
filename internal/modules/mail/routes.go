package mail

import "github.com/labstack/echo/v4"

func RegisterRoutes(authGroup *echo.Group, h *Handler) {
	if authGroup == nil {
		return
	}
	if h == nil {
		h = &Handler{}
	}

	authGroup.GET("/mail/servers", h.ListServers)
	authGroup.POST("/mail/servers", h.CreateServer)
	authGroup.GET("/mail/servers/:id", h.GetServer)
	authGroup.DELETE("/mail/servers/:id", h.DeleteServer)
	authGroup.GET("/mail/servers/:id/domains", h.ListDomains)
	authGroup.POST("/mail/servers/:id/domains", h.CreateDomain)
	authGroup.DELETE("/mail/servers/:id/domains/:domainId", h.DeleteDomain)
	authGroup.GET("/mail/servers/:id/domains/:domainId/mailboxes", h.ListMailboxes)
	authGroup.POST("/mail/servers/:id/domains/:domainId/mailboxes", h.CreateMailbox)
	authGroup.DELETE("/mail/servers/:id/domains/:domainId/mailboxes/:mailboxId", h.DeleteMailbox)
	authGroup.POST("/mail/servers/:id/domains/:domainId/dkim", h.CreateDKIM)
	authGroup.GET("/mail/servers/:id/domains/:domainId/dns", h.GetDNS)
	authGroup.GET("/mail/servers/:id/domains/:domainId/checks", h.GetChecks)
	authGroup.GET("/mail/servers/:id/certs", h.ListCerts)
	authGroup.POST("/mail/servers/:id/certs/probe", h.ProbeCert)
	authGroup.GET("/mail/servers/:id/fail2ban/status", h.GetFail2banStatus)
	authGroup.POST("/mail/servers/:id/fail2ban/ban", h.BanIP)
	authGroup.POST("/mail/servers/:id/fail2ban/unban", h.UnbanIP)
	authGroup.GET("/mail/servers/:id/domains/:domainId/aliases", h.ListAliases)
	authGroup.POST("/mail/servers/:id/domains/:domainId/aliases", h.CreateAlias)
	authGroup.DELETE("/mail/servers/:id/domains/:domainId/aliases/:aliasId", h.DeleteAlias)
}
