package utils

import (
	"fmt"
	"strings"

	"codedock.run/codedock/internal/config"
)

func DefaultDBMemoryMB() int {
	mb := int(config.Get().Defaults.DBMemoryMB)
	if mb <= 0 {
		return 1024
	}
	return mb
}

func DefaultDBCPURequest() float64 {
	cpu := config.Get().Defaults.DBCPU
	if cpu <= 0 {
		return 1.0
	}
	return cpu
}

func MegaBytesToBytes(mb int) int64 {
	if mb <= 0 {
		return 512 * 1024 * 1024
	}
	return int64(mb) * 1024 * 1024
}

func CPURequestToNanoCPUs(cores float64) int64 {
	if cores <= 0 {
		return 500_000_000
	}
	return int64(cores * 1_000_000_000)
}

func NormalizeContainerName(projectID string) string {
	return fmt.Sprintf("codedock-%s", strings.ToLower(strings.TrimSpace(projectID)))
}
