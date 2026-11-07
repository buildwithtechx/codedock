package repositories

import (
	"context"
	"testing"

	"codedock/internal/models"
)

func TestIngressRepositoryInMemory(t *testing.T) {
	repo := NewIngressRepository(nil)
	ctx := context.Background()

	conn := &models.IngressConnection{
		ID:        "conn-1",
		ServerID:  "server-1",
		Name:      "Production Tunnel",
		Status:    "ready",
		CreatedAt: "2026-10-17T00:00:00Z",
		UpdatedAt: "2026-10-17T00:00:00Z",
	}

	if err := repo.SaveConnection(ctx, conn); err != nil {
		t.Fatalf("failed to save connection: %v", err)
	}

	list, err := repo.ListConnections(ctx, "server-1")
	if err != nil {
		t.Fatalf("failed to list connections: %v", err)
	}
	if len(list) != 1 || list[0].ID != "conn-1" {
		t.Fatalf("unexpected list result: %+v", list)
	}

	fetched, err := repo.GetConnection(ctx, "conn-1")
	if err != nil || fetched.Name != "Production Tunnel" {
		t.Fatalf("unexpected connection fetched: %+v, err: %v", fetched, err)
	}

	domainID := "domain-1"
	host := &models.IngressHostname{
		ID:           "host-1",
		DomainID:     &domainID,
		ConnectionID: "conn-1",
		Hostname:     "app.example.com",
		Status:       "ready",
		Phase:        "dns",
	}

	if err := repo.SaveHostname(ctx, host); err != nil {
		t.Fatalf("failed to save hostname: %v", err)
	}

	fetchedHost, err := repo.GetHostname(ctx, domainID)
	if err != nil || fetchedHost == nil || fetchedHost.Hostname != "app.example.com" {
		t.Fatalf("unexpected hostname fetched: %+v, err: %v", fetchedHost, err)
	}

	if err := repo.DeleteHostname(ctx, domainID); err != nil {
		t.Fatalf("failed to delete hostname: %v", err)
	}
	delHost, _ := repo.GetHostname(ctx, domainID)
	if delHost != nil {
		t.Fatalf("expected nil after delete, got %+v", delHost)
	}

	if err := repo.DeleteConnection(ctx, "conn-1"); err != nil {
		t.Fatalf("failed to delete connection: %v", err)
	}
	delList, _ := repo.ListConnections(ctx, "server-1")
	if len(delList) != 0 {
		t.Fatalf("expected 0 connections, got %d", len(delList))
	}
}
