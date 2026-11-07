package compose

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/container"
)

const (
	MinComposeShmBytes int64 = 1
	MaxComposeShmBytes int64 = 1024 * 1024 * 1024 * 1024 // 1 TiB
)

var (
	ErrInvalidShmSize = errors.New("shm_size must be bytes or a size such as 2g")
	ErrShmSizeBounds  = errors.New("shared memory must be between 1 byte and 1 TiB")
	shmSizePattern    = regexp.MustCompile(`^(\d+(?:\.\d+)?)\s*([kmgt]?)b?$`)
)

func ParseComposeShmSize(value any) (int64, error) {
	if value == nil {
		return 0, nil
	}

	var bytes float64

	switch v := value.(type) {
	case int:
		bytes = float64(v)
	case int32:
		bytes = float64(v)
	case int64:
		bytes = float64(v)
	case uint:
		bytes = float64(v)
	case uint32:
		bytes = float64(v)
	case uint64:
		bytes = float64(v)
	case float64:
		bytes = v
	case string:
		trimmed := strings.ToLower(strings.TrimSpace(v))
		if trimmed == "" {
			return 0, nil
		}
		matches := shmSizePattern.FindStringSubmatch(trimmed)
		if len(matches) != 3 {
			return 0, ErrInvalidShmSize
		}

		num, err := strconv.ParseFloat(matches[1], 64)
		if err != nil || math.IsNaN(num) || math.IsInf(num, 0) {
			return 0, ErrInvalidShmSize
		}

		unit := matches[2]
		var multiplier float64 = 1
		switch unit {
		case "k":
			multiplier = 1024
		case "m":
			multiplier = 1024 * 1024
		case "g":
			multiplier = 1024 * 1024 * 1024
		case "t":
			multiplier = 1024 * 1024 * 1024 * 1024
		}
		bytes = num * multiplier
	default:
		return 0, ErrInvalidShmSize
	}

	if bytes < float64(MinComposeShmBytes) || bytes > float64(MaxComposeShmBytes) || math.IsNaN(bytes) || math.IsInf(bytes, 0) {
		return 0, ErrShmSizeBounds
	}

	intBytes := int64(math.Round(bytes))
	if intBytes < MinComposeShmBytes || intBytes > MaxComposeShmBytes {
		return 0, ErrShmSizeBounds
	}

	return intBytes, nil
}

func ApplyHostConfigShmSize(hostConfig *container.HostConfig, shmSize int64) {
	if hostConfig != nil && shmSize > 0 {
		hostConfig.ShmSize = shmSize
	}
}

func FormatShmSize(bytes int64) string {
	if bytes <= 0 {
		return ""
	}
	switch {
	case bytes >= 1024*1024*1024 && bytes%(1024*1024*1024) == 0:
		return fmt.Sprintf("%dgb", bytes/(1024*1024*1024))
	case bytes >= 1024*1024 && bytes%(1024*1024) == 0:
		return fmt.Sprintf("%dmb", bytes/(1024*1024))
	case bytes >= 1024 && bytes%1024 == 0:
		return fmt.Sprintf("%dkb", bytes/1024)
	default:
		return fmt.Sprintf("%db", bytes)
	}
}
