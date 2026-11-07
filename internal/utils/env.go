package utils

import "codedock/internal/config"

func IsDryRun() bool {
	return config.Get().Docker.DryRun
}
