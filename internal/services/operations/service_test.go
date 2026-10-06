package operations

import (
	"codedock.run/codedock/internal/models"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

type reviewedStore struct {
	Store
	operation *models.Operation
	claimed   bool
}

func (s *reviewedStore) Get(context.Context, string) (*models.Operation, error) {
	return s.operation, nil
}
func (s *reviewedStore) Claim(context.Context, string, string) error { s.claimed = true; return nil }
func TestInvalidReviewsNeverClaimOrExecute(t *testing.T) {
	digest := sha256.Sum256([]byte("valid-confirmation"))
	for _, test := range []struct{ user, token, snapshot string }{{"other", "valid-confirmation", "version"}, {"owner", "invalid", "version"}, {"owner", "valid-confirmation", "changed"}} {
		store := &reviewedStore{operation: &models.Operation{UserID: "owner", TokenHash: hex.EncodeToString(digest[:]), Snapshot: "version", Status: "REVIEWED"}}
		service := NewService(store)
		if err := service.Apply(context.Background(), "operation", test.user, test.token, test.snapshot, func(context.Context, *models.Operation, func(string, string) error) error {
			t.Error("rejected review executed")
			return nil
		}); err == nil || store.claimed {
			t.Fatal("invalid review claimed target")
		}
	}
}
