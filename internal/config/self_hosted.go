package config

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"codedock/internal/models"
)

type SelfHostedStore interface {
	Load(context.Context) (*models.SelfHostedConfig, error)
	Save(context.Context, *models.SelfHostedConfig) error
}

func PrepareSelfHosted(ctx context.Context, cfg *models.Config, store SelfHostedStore) error {
	if cfg.Cloud.Enabled {
		return nil
	}
	stored, err := store.Load(ctx)
	if err != nil {
		return err
	}
	changed := false
	for _, pair := range [][2]*string{
		{&cfg.Security.JWTSecret, &stored.JWTSecret},
		{&cfg.Security.RefreshSecret, &stored.RefreshSecret},
		{&cfg.Telemetry.Salt, &stored.TelemetrySalt},
	} {
		if *pair[1] == "" {
			*pair[1] = *pair[0]
			if *pair[1] == "" {
				random := make([]byte, 32)
				if _, err := rand.Read(random); err != nil {
					return fmt.Errorf("generate authentication secret: %w", err)
				}
				*pair[1] = hex.EncodeToString(random)
			}
			changed = true
		}
		if *pair[0] == "" {
			*pair[0] = *pair[1]
		}
	}
	if cfg.Security.TLSEmail == "" {
		cfg.Security.TLSEmail = stored.TLSEmail
	}
	if cfg.Domains.WildcardDomain == "" {
		cfg.Domains.WildcardDomain = stored.WildcardDomain
	}
	if changed {
		return store.Save(ctx, stored)
	}
	return nil
}

func SaveSelfHostedOptions(ctx context.Context, cfg *models.Config, store SelfHostedStore, domain, tlsEmail string) error {
	if cfg.Cloud.Enabled {
		return errors.New("self-hosted setup is unavailable in cloud mode")
	}
	stored, err := store.Load(ctx)
	if err != nil {
		return err
	}
	stored.WildcardDomain = domain
	stored.TLSEmail = tlsEmail
	return store.Save(ctx, stored)
}
