package models

import "strings"

type AuditLog struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	Action    string `json:"action"`
	Resource  string `json:"resource"`
	Details   string `json:"details"`
	IPAddress string `json:"ipAddress"`
	CreatedAt string `json:"createdAt"`
	Category  string `json:"category"`
}

type AuditSettings struct {
	OrganizationID string `json:"organizationId"`
	Enabled        bool   `json:"enabled"`
	RetentionDays  int    `json:"retentionDays"`
	UpdatedAt      string `json:"updatedAt"`
}

type AuditCategory struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

type AuditCategoryCount struct {
	AuditCategory
	Count int `json:"count"`
}

type AuditFacets struct {
	Total      int                  `json:"total"`
	Categories []AuditCategoryCount `json:"categories"`
}

func AuditCategories() []AuditCategory {
	return []AuditCategory{
		{ID: "deployments", Label: "Deployments", Description: "Deploys, rollbacks and releases"},
		{ID: "apps", Label: "Apps & services", Description: "Services, variables, databases and stacks"},
		{ID: "domains", Label: "Domains & SSL", Description: "Domains, DNS records and certificates"},
		{ID: "servers", Label: "Servers", Description: "Servers and SSH access"},
		{ID: "members", Label: "Members & access", Description: "Organizations, members and invites"},
		{ID: "agent", Label: "AI agents", Description: "AI-assisted operations"},
		{ID: "security", Label: "Security", Description: "Sign-ins, tokens and credentials"},
		{ID: "billing", Label: "Billing", Description: "Plans, checkouts and credits"},
		{ID: "system", Label: "System", Description: "Backups, maintenance and everything else"},
	}
}

func categoryPrefixes() map[string][]string {
	return map[string][]string{
		"deployments": {"deployment.", "project.deployment."},
		"apps":        {"service_", "service.", "app.", "database.", "oneclick.", "compose.", "stack.", "environment.", "preview."},
		"domains":     {"domain.", "dns.", "certificate.", "route."},
		"servers":     {"server.", "ssh."},
		"members":     {"member.", "organization.", "invite.", "user.create", "user.delete", "user.role"},
		"agent":       {"ai.", "agent."},
		"security":    {"auth.", "login.", "user.login", "user.logout", "user.password", "user.2fa", "user.totp", "token.", "apikey.", "oauth.", "credentials."},
		"billing":     {"billing.", "subscription.", "checkout.", "credit."},
	}
}

func CategoryForAction(action string) string {
	for category, prefixes := range categoryPrefixes() {
		for _, prefix := range prefixes {
			if strings.HasPrefix(action, prefix) {
				return category
			}
		}
	}
	return "system"
}

func CategoryPrefixes(category string) []string {
	return categoryPrefixes()[category]
}
