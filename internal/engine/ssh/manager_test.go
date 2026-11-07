package ssh

import (
	"testing"

	"codedock/internal/models"
)

func TestEvictedClientReleaseDoesNotAffectReplacement(t *testing.T) {
	manager := NewSSHManager(nil)
	old := &clientEntry{client: &Client{}, refs: 1}
	manager.clients["server"] = old
	manager.RemoveClient("server")
	if _, exists := manager.clients["server"]; exists {
		t.Fatal("evicted client remains cached")
	}
	replacement := &clientEntry{client: &Client{}, refs: 2}
	manager.clients["server"] = replacement
	manager.releaseClient(old)
	if old.refs != 0 || replacement.refs != 2 || manager.clients["server"] != replacement {
		t.Fatal("old release changed the replacement client")
	}
}

func TestAssertOperatorSSHAccess(t *testing.T) {
	server := &models.Server{
		ID:             "srv-1",
		OrganizationID: "org-1",
	}

	if err := AssertOperatorSSHAccess(server, "org-1"); err != nil {
		t.Fatalf("expected access permitted for same org, got %v", err)
	}

	if err := AssertOperatorSSHAccess(server, "org-2"); err == nil {
		t.Fatalf("expected tenant boundary error for different org")
	}

	if err := AssertOperatorSSHAccess(nil, "org-1"); err == nil {
		t.Fatalf("expected error for nil server")
	}
}
