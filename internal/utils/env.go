package utils

import "codedock.run/codedock/internal/config"

func IsDryRun() bool {
	return config.Get().Docker.DryRun
}
