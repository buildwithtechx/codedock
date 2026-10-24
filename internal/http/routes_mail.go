package http

import "github.com/labstack/echo/v4"

func (s *Server) registerMailRoutes(authGroup *echo.Group) {
	authGroup.GET("/mail/servers", s.mailHandler.ListServers)
	authGroup.POST("/mail/servers", s.mailHandler.CreateServer)
	authGroup.GET("/mail/servers/:id", s.mailHandler.GetServer)
	authGroup.DELETE("/mail/servers/:id", s.mailHandler.DeleteServer)
	authGroup.GET("/mail/servers/:id/domains", s.mailHandler.ListDomains)
	authGroup.POST("/mail/servers/:id/domains", s.mailHandler.CreateDomain)
	authGroup.DELETE("/mail/servers/:id/domains/:domainId", s.mailHandler.DeleteDomain)
	authGroup.GET("/mail/servers/:id/domains/:domainId/mailboxes", s.mailHandler.ListMailboxes)
	authGroup.POST("/mail/servers/:id/domains/:domainId/mailboxes", s.mailHandler.CreateMailbox)
	authGroup.DELETE("/mail/servers/:id/domains/:domainId/mailboxes/:mailboxId", s.mailHandler.DeleteMailbox)
	authGroup.POST("/mail/servers/:id/domains/:domainId/dkim", s.mailHandler.CreateDKIM)
	authGroup.GET("/mail/servers/:id/domains/:domainId/dns", s.mailHandler.GetDNS)
	authGroup.GET("/mail/servers/:id/domains/:domainId/checks", s.mailHandler.GetChecks)
	authGroup.GET("/mail/servers/:id/certs", s.mailHandler.ListCerts)
	authGroup.POST("/mail/servers/:id/certs/probe", s.mailHandler.ProbeCert)
}
