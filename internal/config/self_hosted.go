package config

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/pkg/types"
)

func PrepareSelfHosted(cfg *types.Config) error {
	if cfg.Cloud.Enabled {
		return nil
	}
	stored, err := readSelfHostedConfig(cfg.Server.DataDir)
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
		return writeSelfHostedConfig(cfg.Server.DataDir, stored)
	}
	return nil
}

func SaveSelfHostedOptions(cfg *types.Config, domain, tlsEmail string) error {
	if cfg.Cloud.Enabled {
		return errors.New("self-hosted setup is unavailable in cloud mode")
	}
	stored, err := readSelfHostedConfig(cfg.Server.DataDir)
	if err != nil {
		return err
	}
	stored.WildcardDomain = domain
	stored.TLSEmail = tlsEmail
	return writeSelfHostedConfig(cfg.Server.DataDir, stored)
}

func readSelfHostedConfig(dataDir string) (*models.SelfHostedConfig, error) {
	path := filepath.Join(dataDir, "self-hosted.json")
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &models.SelfHostedConfig{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read self-hosted configuration: %w", err)
	}
	var stored models.SelfHostedConfig
	if err := json.Unmarshal(data, &stored); err != nil {
		return nil, fmt.Errorf("decode self-hosted configuration: %w", err)
	}
	if err := os.Chmod(path, 0o600); err != nil {
		return nil, fmt.Errorf("protect self-hosted configuration: %w", err)
	}
	return &stored, nil
}

func writeSelfHostedConfig(dataDir string, stored *models.SelfHostedConfig) error {
	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		return fmt.Errorf("create configuration directory: %w", err)
	}
	data, err := json.Marshal(stored)
	if err != nil {
		return fmt.Errorf("encode self-hosted configuration: %w", err)
	}
	file, err := os.CreateTemp(dataDir, ".self-hosted-*")
	if err != nil {
		return fmt.Errorf("create self-hosted configuration: %w", err)
	}
	defer func() {
		if err := os.Remove(file.Name()); err != nil && !errors.Is(err, os.ErrNotExist) {
			slog.Warn("remove temporary configuration", "err", err)
		}
	}()
	_, writeErr := file.Write(data)
	closeErr := file.Close()
	if err := errors.Join(writeErr, closeErr); err != nil {
		return fmt.Errorf("write self-hosted configuration: %w", err)
	}
	if err := os.Rename(file.Name(), filepath.Join(dataDir, "self-hosted.json")); err != nil {
		return fmt.Errorf("save self-hosted configuration: %w", err)
	}
	return nil
}
