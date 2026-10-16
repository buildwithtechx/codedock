package build

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	MaxBuildStorageReserveBytes uint64 = 2 * 1024 * 1024 * 1024
	DefaultBuildCacheKeepBytes  uint64 = 5 * 1024 * 1024 * 1024
)

var (
	ErrBuildStorageLow     = errors.New("the build host has too little free storage; unused build caches were cleaned; free space or increase disk capacity before retrying")
	ErrBuildStorageUnknown = errors.New("could not measure the build host's free storage")
)

type BuildStorage struct {
	TotalBytes uint64 `json:"totalBytes"`
	FreeBytes  uint64 `json:"freeBytes"`
	FreeInodes *int64 `json:"freeInodes,omitempty"`
}

var packageCaches = map[string][]string{
	"npm":      {"/root/.npm"},
	"bun":      {"/root/.bun/install/cache"},
	"pnpm":     {"/root/.local/share/pnpm/store"},
	"pip":      {"/root/.cache/pip"},
	"pipenv":   {"/root/.cache/pip"},
	"poetry":   {"/root/.cache/pypoetry/cache", "/root/.cache/pypoetry/artifacts"},
	"uv":       {"/root/.cache/uv"},
	"composer": {"/root/.cache/composer"},
	"bundler":  {"/usr/local/bundle/cache"},
	"cargo":    {"/usr/local/cargo/registry", "/usr/local/cargo/git"},
	"go":       {"/go/pkg/mod", "/root/.cache/go-build"},
	"maven":    {"/root/.m2/repository"},
	"gradle":   {"/root/.gradle/caches", "/root/.gradle/wrapper"},
	"dotnet":   {"/root/.nuget/packages"},
	"mix":      {"/root/.hex/packages"},
}

func GetPackageCacheDirs(manager string) []string {
	if dirs, ok := packageCaches[strings.ToLower(strings.TrimSpace(manager))]; ok {
		res := make([]string, len(dirs))
		copy(res, dirs)
		return res
	}
	return nil
}

func BuildStorageReserve(storage BuildStorage) uint64 {
	tenPercent := (storage.TotalBytes + 9) / 10
	if tenPercent < MaxBuildStorageReserveBytes {
		return tenPercent
	}
	return MaxBuildStorageReserveBytes
}

func AssertBuildStorage(storage BuildStorage) error {
	if storage.FreeBytes < BuildStorageReserve(storage) || (storage.FreeInodes != nil && *storage.FreeInodes == 0) {
		return ErrBuildStorageLow
	}
	return nil
}

func ParseBuildStorage(raw string) (BuildStorage, error) {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	var storageLine, inodesLine string
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if strings.HasPrefix(trimmed, "CODEDOCK_STORAGE ") || strings.HasPrefix(trimmed, "OPENSHIP_STORAGE ") {
			parts := strings.SplitN(trimmed, " ", 2)
			if len(parts) == 2 {
				storageLine = parts[1]
			}
		} else if strings.HasPrefix(trimmed, "CODEDOCK_INODES ") || strings.HasPrefix(trimmed, "OPENSHIP_INODES ") {
			parts := strings.SplitN(trimmed, " ", 2)
			if len(parts) == 2 {
				inodesLine = parts[1]
			}
		}
	}

	var totalKb, freeKb uint64
	var freeInodes *int64

	if storageLine != "" {
		fields := strings.Fields(storageLine)
		if len(fields) < 2 {
			return BuildStorage{}, ErrBuildStorageUnknown
		}
		t, err1 := strconv.ParseUint(fields[0], 10, 64)
		f, err2 := strconv.ParseUint(fields[1], 10, 64)
		if err1 != nil || err2 != nil || t == 0 || f > t {
			return BuildStorage{}, ErrBuildStorageUnknown
		}
		totalKb = t
		freeKb = f

		if inodesLine != "" {
			inFields := strings.Fields(inodesLine)
			if len(inFields) > 0 {
				if inVal, err := strconv.ParseInt(inFields[0], 10, 64); err == nil && inVal >= 0 {
					freeInodes = &inVal
				}
			}
		}
	} else {
		fields := strings.Fields(strings.TrimSpace(raw))
		if len(fields) < 2 {
			return BuildStorage{}, ErrBuildStorageUnknown
		}
		t, err1 := strconv.ParseUint(fields[0], 10, 64)
		f, err2 := strconv.ParseUint(fields[1], 10, 64)
		if err1 != nil || err2 != nil || t == 0 || f > t {
			return BuildStorage{}, ErrBuildStorageUnknown
		}
		totalKb = t
		freeKb = f
		if len(fields) >= 3 {
			if inVal, err := strconv.ParseInt(fields[2], 10, 64); err == nil && inVal >= 0 {
				freeInodes = &inVal
			}
		}
	}

	return BuildStorage{
		TotalBytes: totalKb * 1024,
		FreeBytes:  freeKb * 1024,
		FreeInodes: freeInodes,
	}, nil
}

func IsBuildCacheMountsEnabled(envVars map[string]string) bool {
	if strings.EqualFold(os.Getenv("CODEDOCK_BUILD_CACHE_MOUNTS"), "false") ||
		strings.EqualFold(os.Getenv("OPENSHIP_BUILD_CACHE_MOUNTS"), "false") {
		return false
	}
	if envVars != nil {
		if strings.EqualFold(envVars["CODEDOCK_BUILD_CACHE_MOUNTS"], "false") ||
			strings.EqualFold(envVars["OPENSHIP_BUILD_CACHE_MOUNTS"], "false") {
			return false
		}
	}
	return true
}

var safeCacheMountPath = regexp.MustCompile(`^/[A-Za-z0-9._/-]+$`)

func FormatCacheMount(id, targetPath string) (string, error) {
	if !safeCacheMountPath.MatchString(targetPath) {
		return "", fmt.Errorf("invalid cache mount target path: %s", targetPath)
	}
	return fmt.Sprintf("--mount=type=cache,id=codedock-%s-${TARGETPLATFORM},target=%s,sharing=locked", id, targetPath), nil
}
