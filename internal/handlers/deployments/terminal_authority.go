package deployments

import (
	"context"
	"fmt"
	"sync"
	"time"

	"codedock/internal/models"
	projectservices "codedock/internal/services/projects"
)

// VerifyTerminalAuthority checks user active status and project membership/permissions.
func VerifyTerminalAuthority(
	ctx context.Context,
	projectService *projectservices.ProjectService,
	userRepo userStatusProvider,
	projectID, userID, role string,
) error {
	if userRepo != nil && userID != "" {
		u, err := userRepo.GetUserByID(ctx, userID)
		if err != nil || u == nil || !u.IsActive {
			return fmt.Errorf("user account not found or deactivated")
		}
	}

	if role == "admin" {
		return nil
	}

	if projectService != nil && projectID != "" && userID != "" {
		if !projectService.HasPermission(ctx, projectID, userID, models.UserRole(role), "") {
			return fmt.Errorf("insufficient permissions to access this terminal")
		}
	}

	return nil
}

// StartTerminalAuthorityWatcher runs continuous verification during long-lived
// interactive terminal sessions. If authority is revoked, onDenied is invoked immediately.
func StartTerminalAuthorityWatcher(
	ctx context.Context,
	interval time.Duration,
	verify func() error,
	onDenied func(),
) func() {
	if interval <= 0 {
		interval = 5 * time.Second
	}

	ticker := time.NewTicker(interval)
	done := make(chan struct{})
	var once sync.Once

	stop := func() {
		once.Do(func() {
			close(done)
			ticker.Stop()
		})
	}

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := verify(); err != nil {
					stop()
					if onDenied != nil {
						onDenied()
					}
					return
				}
			}
		}
	}()

	return stop
}
