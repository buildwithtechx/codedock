package repositories

import (
	"context"
	"testing"

	"codedock/internal/models"
)

func TestAuditResourceResolverIsUUID(t *testing.T) {
	if !IsUUID("91fdae7e-559c-41af-8426-814fe0d5e643") {
		t.Fatalf("expected valid UUID to return true")
	}
	if IsUUID("my-cool-project") {
		t.Fatalf("expected project slug to return false")
	}
	if IsUUID("") {
		t.Fatalf("expected empty string to return false")
	}
}

func TestAuditResourceResolverExtractDetails(t *testing.T) {
	name, resType := extractNameFromDetails(`{"name":"web-app","type":"service"}`)
	if name != "web-app" || resType != "service" {
		t.Fatalf("expected web-app/service, got %q/%q", name, resType)
	}

	name2, _ := extractNameFromDetails(`{"after":{"projectName":"marketing-site"}}`)
	if name2 != "marketing-site" {
		t.Fatalf("expected marketing-site, got %q", name2)
	}
}

func TestAuditResourceResolverResolveDetailsFallback(t *testing.T) {
	logs := []models.AuditLog{
		{
			ID:       "log-1",
			Resource: "91fdae7e-559c-41af-8426-814fe0d5e643",
			Details:  `{"name":"solomon-project","type":"project"}`,
		},
		{
			ID:       "log-2",
			Resource: "already-named-service",
		},
	}
	resolveResourceNames(context.Background(), nil, logs)

	if logs[0].ResourceName != "solomon-project" {
		t.Fatalf("expected logs[0].ResourceName to be solomon-project, got %q", logs[0].ResourceName)
	}
	if logs[0].ResourceType != "project" {
		t.Fatalf("expected logs[0].ResourceType to be project, got %q", logs[0].ResourceType)
	}
	if logs[1].ResourceName != "already-named-service" {
		t.Fatalf("expected logs[1].ResourceName to be already-named-service, got %q", logs[1].ResourceName)
	}
}
