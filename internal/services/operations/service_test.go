package operations

import (
	"codedock/internal/models"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"testing"
	"time"
)

type reviewedStore struct {
	Store
	operation *models.Operation
	claimed   bool
}

func (s *reviewedStore) Get(context.Context, string) (*models.Operation, error) {
	return s.operation, nil
}

type cancellableStore struct {
	Store
	mu        sync.Mutex
	operation *models.Operation
}

func (s *cancellableStore) Get(_ context.Context, _ string) (*models.Operation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.operation, nil
}

func (s *cancellableStore) Claim(_ context.Context, _, _ string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.operation.Status = "RUNNING"
	return nil
}

func (s *cancellableStore) Observe(_ context.Context, _, status, _, _, _ string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.operation.Status = status
	return nil
}

func (s *cancellableStore) status() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.operation.Status
}

func TestCancelReachesOperationOnAnotherNode(t *testing.T) {
	digest := sha256.Sum256([]byte("confirmation"))
	store := &cancellableStore{operation: &models.Operation{ID: "op", UserID: "owner", TokenHash: hex.EncodeToString(digest[:]), Snapshot: "version", Status: "REVIEWED"}}
	executor := NewService(store)
	started := make(chan struct{})
	var once sync.Once
	if err := executor.Apply(context.Background(), "op", "owner", "confirmation", "version", func(ctx context.Context, _ *models.Operation, progress func(string, string) error) error {
		once.Do(func() { close(started) })
		for {
			if err := progress("WORKING", "tick"); err != nil {
				return err
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(10 * time.Millisecond):
			}
		}
	}); err != nil {
		t.Fatalf("apply: %v", err)
	}
	<-started
	other := NewService(store)
	if err := other.Cancel(context.Background(), "op", "owner"); err != nil {
		t.Fatalf("cross-node cancel: %v", err)
	}
	deadline := time.Now().Add(5 * time.Second)
	for store.status() != "INTERRUPTED" {
		if time.Now().After(deadline) {
			t.Fatalf("operation did not interrupt, status %s", store.status())
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func TestCancelFinishedOperationFails(t *testing.T) {
	store := &cancellableStore{operation: &models.Operation{ID: "op", UserID: "owner", Status: "COMPLETED"}}
	service := NewService(store)
	if err := service.Cancel(context.Background(), "op", "owner"); err == nil {
		t.Fatal("cancel of a completed operation succeeded")
	}
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
