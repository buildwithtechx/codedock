//go:build windows

package deployments

import (
	"syscall"
	"unsafe"
)

func readHostDiskSpace() (totalBytes int64, freeBytes int64, err error) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getDiskFreeSpaceEx := kernel32.NewProc("GetDiskFreeSpaceExW")
	pPath, err := syscall.UTF16PtrFromString("C:\\")
	if err != nil {
		return 0, 0, err
	}
	var freeBytesAvailable, totalNumberOfBytes, totalNumberOfFreeBytes int64
	r, _, callErr := getDiskFreeSpaceEx.Call(
		uintptr(unsafe.Pointer(pPath)),
		uintptr(unsafe.Pointer(&freeBytesAvailable)),
		uintptr(unsafe.Pointer(&totalNumberOfBytes)),
		uintptr(unsafe.Pointer(&totalNumberOfFreeBytes)),
	)
	if r == 0 {
		return 0, 0, callErr
	}
	return totalNumberOfBytes, freeBytesAvailable, nil
}
