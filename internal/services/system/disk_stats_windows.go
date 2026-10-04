//go:build windows

package system

import (
	"math"
	"syscall"
	"unsafe"

	"codedock.run/codedock/internal/models"
)

func getDiskStats() models.DiskStats {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getDiskFreeSpaceEx := kernel32.NewProc("GetDiskFreeSpaceExW")

	pPath, err := syscall.UTF16PtrFromString("C:\\")
	if err != nil {
		return models.DiskStats{}
	}

	var freeBytesAvailable, totalNumberOfBytes, totalNumberOfFreeBytes int64
	r, _, _ := getDiskFreeSpaceEx.Call(
		uintptr(unsafe.Pointer(pPath)),
		uintptr(unsafe.Pointer(&freeBytesAvailable)),
		uintptr(unsafe.Pointer(&totalNumberOfBytes)),
		uintptr(unsafe.Pointer(&totalNumberOfFreeBytes)),
	)
	if r == 0 {
		return models.DiskStats{}
	}

	total := totalNumberOfBytes / (1024 * 1024 * 1024)
	free := totalNumberOfFreeBytes / (1024 * 1024 * 1024)
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
