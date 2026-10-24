package http

import (
	"context"
	"fmt"
)

func (a *engineAdapter) ClaimRecordExpiry(ctx context.Context, id string) (bool, error) {
	store, ok := a.backupRepo.(interface {
		ClaimRecordExpiry(context.Context, string) (bool, error)
	})
	if !ok {
		return false, fmt.Errorf("backup retention protection storage unavailable")
	}
	return store.ClaimRecordExpiry(ctx, id)
}
