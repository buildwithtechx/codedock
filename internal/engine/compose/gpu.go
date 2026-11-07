package compose

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/docker/docker/api/types/container"

	"codedock/internal/models"
)

var nvidiaGpuIDRegex = regexp.MustCompile(`^(?:all|\d+(?::\d+)?|GPU-[A-Za-z0-9-]+|MIG-[A-Za-z0-9/-]+)$`)

func IsNvidiaGpuID(id string) bool {
	id = strings.TrimSpace(id)
	selector := strings.TrimPrefix(id, "nvidia.com/gpu=")
	return nvidiaGpuIDRegex.MatchString(selector)
}

func NormalizeComposeGpuRequests(value any) ([]models.GPURequest, error) {
	if value == nil {
		return nil, nil
	}
	switch v := value.(type) {
	case string:
		trimmed := strings.ToLower(strings.TrimSpace(v))
		if trimmed == "all" {
			return []models.GPURequest{{Driver: "nvidia", Count: "all", Capabilities: []string{"gpu"}}}, nil
		}
		return nil, errors.New("gpus must be all or a list of device requests")
	case int:
		if v < 0 || v > 64 {
			return nil, errors.New("GPU count must be all or an integer between 0 and 64")
		}
		return []models.GPURequest{{Driver: "nvidia", Count: v, Capabilities: []string{"gpu"}}}, nil
	case []models.GPURequest:
		return NormalizeGPURequests(v)
	case models.GPURequest:
		return NormalizeGPURequests([]models.GPURequest{v})
	case []any:
		reqs := make([]models.GPURequest, 0, len(v))
		for _, item := range v {
			if m, ok := item.(map[string]any); ok {
				var req models.GPURequest
				if drv, ok := m["driver"].(string); ok {
					req.Driver = drv
				}
				if cnt, ok := m["count"]; ok {
					req.Count = cnt
				}
				if caps, ok := m["capabilities"].([]any); ok {
					for _, c := range caps {
						if cs, ok := c.(string); ok {
							req.Capabilities = append(req.Capabilities, cs)
						}
					}
				}
				if devs, ok := m["device_ids"].([]any); ok {
					for _, d := range devs {
						if ds, ok := d.(string); ok {
							req.DeviceIDs = append(req.DeviceIDs, ds)
						}
					}
				}
				reqs = append(reqs, req)
			} else {
				return nil, errors.New("each GPU request must be an object")
			}
		}
		return NormalizeGPURequests(reqs)
	case map[string]any:
		return NormalizeComposeGpuRequests([]any{v})
	default:
		return nil, errors.New("gpus must be all or a list of device requests")
	}
}

func NormalizeGPURequests(requests []models.GPURequest) ([]models.GPURequest, error) {
	if len(requests) == 0 {
		return nil, nil
	}
	if len(requests) > 16 {
		return nil, errors.New("at most 16 GPU request groups are supported")
	}

	normalized := make([]models.GPURequest, 0, len(requests))
	for _, req := range requests {
		driver := strings.TrimSpace(strings.ToLower(req.Driver))
		if driver == "" {
			driver = "nvidia"
		}
		if driver != "nvidia" {
			return nil, errors.New("GPU requests support only the NVIDIA driver")
		}

		out := models.GPURequest{
			Driver: "nvidia",
		}

		hasCount := req.Count != nil
		hasDevices := len(req.DeviceIDs) > 0

		if hasCount && hasDevices {
			return nil, errors.New("GPU count and device_ids cannot be combined")
		}

		if hasCount {
			switch v := req.Count.(type) {
			case string:
				trimmed := strings.TrimSpace(strings.ToLower(v))
				if trimmed != "all" {
					return nil, errors.New("GPU count must be all or an integer between 0 and 64")
				}
				out.Count = "all"
			case int:
				if v < 0 || v > 64 {
					return nil, errors.New("GPU count must be all or an integer between 0 and 64")
				}
				out.Count = v
			case float64:
				iv := int(v)
				if float64(iv) != v || iv < 0 || iv > 64 {
					return nil, errors.New("GPU count must be all or an integer between 0 and 64")
				}
				out.Count = iv
			default:
				return nil, errors.New("GPU count must be all or an integer between 0 and 64")
			}
		}

		if hasDevices {
			if len(req.DeviceIDs) > 64 {
				return nil, errors.New("GPU device_ids must contain between 1 and 64 non-empty IDs")
			}
			outDevices := make([]string, 0, len(req.DeviceIDs))
			for _, id := range req.DeviceIDs {
				trimmed := strings.TrimSpace(id)
				if trimmed == "" || len(trimmed) > 128 {
					return nil, errors.New("GPU device_ids must contain between 1 and 64 non-empty IDs")
				}
				if !IsNvidiaGpuID(trimmed) {
					return nil, errors.New("GPU device_ids must select NVIDIA GPU indexes, UUIDs or NVIDIA GPU CDI names")
				}
				outDevices = append(outDevices, trimmed)
			}
			out.DeviceIDs = outDevices
		}

		if !hasCount && !hasDevices {
			out.Count = "all"
		}

		if len(req.Capabilities) > 0 {
			if len(req.Capabilities) > 16 {
				return nil, errors.New("GPU capabilities must contain between 1 and 16 non-empty names")
			}
			outCaps := make([]string, 0, len(req.Capabilities))
			for _, cap := range req.Capabilities {
				trimmed := strings.TrimSpace(cap)
				if trimmed == "" || len(trimmed) > 64 {
					return nil, errors.New("GPU capabilities must contain between 1 and 16 non-empty names")
				}
				outCaps = append(outCaps, trimmed)
			}
			out.Capabilities = outCaps
		} else {
			out.Capabilities = []string{"gpu"}
		}

		if len(req.Options) > 0 {
			if len(req.Options) > 16 {
				return nil, errors.New("GPU driver options must be a string map with at most 16 entries")
			}
			outOptions := make(map[string]string)
			for k, val := range req.Options {
				k = strings.TrimSpace(k)
				if k == "" || len(k) > 64 || len(val) > 1000 {
					return nil, errors.New("GPU driver options must be a string map with at most 16 entries")
				}
				outOptions[k] = val
			}
			out.Options = outOptions
		}

		normalized = append(normalized, out)
	}

	return normalized, nil
}

func ToDockerDeviceRequests(requests []models.GPURequest) ([]container.DeviceRequest, error) {
	normalized, err := NormalizeGPURequests(requests)
	if err != nil {
		return nil, fmt.Errorf("normalize GPU requests: %w", err)
	}

	deviceReqs := make([]container.DeviceRequest, 0, len(normalized))
	for _, req := range normalized {
		count := 0
		if req.Count == "all" {
			count = -1
		} else if c, ok := req.Count.(int); ok {
			count = c
		}

		deviceReqs = append(deviceReqs, container.DeviceRequest{
			Driver:       "nvidia",
			Count:        count,
			DeviceIDs:    req.DeviceIDs,
			Capabilities: [][]string{req.Capabilities},
			Options:      req.Options,
		})
	}

	return deviceReqs, nil
}
