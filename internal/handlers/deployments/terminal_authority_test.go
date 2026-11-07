package deployments

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"codedock/internal/models"
)

type mockUserProvider struct {
	user *models.User
	err  error
}

func (m *mockUserProvider) GetUserByID(_ context.Context, _ string) (*models.User, error) {
	return m.user, m.err
}

func TestVerifyTerminalAuthorityAdmin(t *testing.T) {
	err := VerifyTerminalAuthority(context.Background(), nil, nil, "proj-1", "user-1", "admin")
	if err != nil {
		t.Fatalf("expected admin to have authority, got error: %v", err)
	}
}

func TestVerifyTerminalAuthorityDeactivatedUser(t *testing.T) {
	up := &mockUserProvider{
		user: &models.User{ID: "user-1", IsActive: false},
	}
	err := VerifyTerminalAuthority(context.Background(), nil, up, "proj-1", "user-1", "member")
	if err == nil {
		t.Fatal("expected error for deactivated user")
	}
}

func TestStartTerminalAuthorityWatcherDenied(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var calls int32
	var deniedTriggered int32

	stop := StartTerminalAuthorityWatcher(
		ctx,
		10*time.Millisecond,
		func() error {
			current := atomic.AddInt32(&calls, 1)
			if current >= 2 {
				return fmt.Errorf("permission revoked")
			}
			return nil
		},
		func() {
			atomic.StoreInt32(&deniedTriggered, 1)
		},
	)
	defer stop()

	time.Sleep(50 * time.Millisecond)

	if atomic.LoadInt32(&deniedTriggered) != 1 {
		t.Fatal("expected onDenied callback to be triggered when permission revoked")
	}
}
