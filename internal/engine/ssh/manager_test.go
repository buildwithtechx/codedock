package ssh

import "testing"

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
