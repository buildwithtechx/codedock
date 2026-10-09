package models

import "testing"

func TestCategoryForAction(t *testing.T) {
	cases := map[string]string{
		"deployment.trigger":         "deployments",
		"deployment.rollback":        "deployments",
		"project.deployment.trigger": "deployments",
		"app.create":                 "apps",
		"app.delete":                 "apps",
		"service_var.update":         "apps",
		"database.query":             "apps",
		"domain.create":              "domains",
		"dns.delete":                 "domains",
		"server.create":              "servers",
		"member.invite":              "members",
		"member.role":                "members",
		"organization.create":        "members",
		"user.delete":                "members",
		"ai.deploy":                  "agent",
		"auth.login":                 "security",
		"auth.signup":                "security",
		"token.create":               "security",
		"user.password":              "security",
		"billing.checkout":           "billing",
		"backup.trigger":             "system",
		"system.restart":             "system",
		"unknown.action":             "system",
		"":                           "system",
	}
	for action, want := range cases {
		if got := CategoryForAction(action); got != want {
			t.Errorf("CategoryForAction(%q) = %q, want %q", action, got, want)
		}
	}
}

func TestAuditCategoriesCoverPrefixes(t *testing.T) {
	seen := map[string]bool{}
	for _, category := range AuditCategories() {
		if seen[category.ID] {
			t.Fatalf("duplicate category %q", category.ID)
		}
		seen[category.ID] = true
		if category.Label == "" || category.Description == "" {
			t.Errorf("category %q missing label or description", category.ID)
		}
	}
	if len(seen) != 9 {
		t.Errorf("expected 9 categories, got %d", len(seen))
	}
}
