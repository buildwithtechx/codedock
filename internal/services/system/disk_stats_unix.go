//go:build !windows

package system

import (
	"math"
	"syscall"

	"codedock.run/codedock/internal/models"
)

func getDiskStats() models.DiskStats {
	var stat syscall.Statfs_t
	if err := syscall.Statfs("/", &stat); err != nil {
		return models.DiskStats{}
	}

	total := int64(stat.Blocks * uint64(stat.Bsize) / (1024 * 1024 * 1024))
	free := int64(stat.Bavail * uint64(stat.Bsize) / (1024 * 1024 * 1024))
	used := total - free

	var percent float64
	if total > 0 {
		percent = float64(used) / float64(total) * 100
		percent = math.Round(percent*10) / 10
	}

	return models.DiskStats{
		TotalGB: total,
		UsedGB:  used,
		FreeGB:  free,
		Percent: percent,
	}
}
