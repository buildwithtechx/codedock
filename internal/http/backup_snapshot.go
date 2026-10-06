package http

import (
	"context"
	"fmt"
)

func (a *engineAdapter) ControlPlaneSnapshot(ctx context.Context) ([]byte, error) {
	store, ok := a.backupRepo.(interface {
		ControlPlaneSnapshot(context.Context) ([]byte, error)
	})
	if !ok {
		return nil, fmt.Errorf("consistent control-plane snapshots unavailable")
	}
	return store.ControlPlaneSnapshot(ctx)
}
