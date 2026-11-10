package deployments

import (
	"context"
	"sync"
	"time"
)

func VerifyTerminalAuthority(
	ctx context.Context,
	projectLookup ProjectLookup,
	userLookup UserStatusLookup,
	projectID, userID, role string,
) error {
	if userLookup != nil && userID != "" {
		active, err := userLookup.IsUserActive(ctx, userID)
		if err != nil || !active {
			return ErrUserDeactivated
		}
	}

	if role == "admin" {
		return nil
	}

	if projectLookup != nil && projectID != "" && userID != "" {
		if !projectLookup.HasPermission(ctx, projectID, userID, role, "") {
			return ErrTerminalUnauthorized
		}
	}

	return nil
}

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
