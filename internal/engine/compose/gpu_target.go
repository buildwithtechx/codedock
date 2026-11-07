package compose

import (
	"fmt"
	"strings"
)

type GpuOverride string

const (
	GpuOverrideYes GpuOverride = "yes"
	GpuOverrideNo  GpuOverride = "no"
)

type HostGpuCapacity struct {
	GpuRuntime bool `json:"gpuRuntime"`
	GpuCount   *int `json:"gpuCount,omitempty"`
}

type EffectiveGpuResult struct {
	HasGpu bool   `json:"hasGpu"`
	Count  *int   `json:"count,omitempty"`
	Source string `json:"source"`
}

type GpuTargetCheckResult struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func EffectiveGpu(capacity *HostGpuCapacity, override GpuOverride) EffectiveGpuResult {
	var count *int
	if capacity != nil {
		count = capacity.GpuCount
	}

	if override == GpuOverrideYes {
		return EffectiveGpuResult{HasGpu: true, Count: count, Source: "override"}
	}
	if override == GpuOverrideNo {
		return EffectiveGpuResult{HasGpu: false, Count: nil, Source: "override"}
	}

	seen := capacity != nil && (capacity.GpuRuntime || count != nil)
	if seen {
		return EffectiveGpuResult{HasGpu: true, Count: count, Source: "probe"}
	}
	return EffectiveGpuResult{HasGpu: false, Count: nil, Source: "none"}
}

func CheckGpuTarget(serviceNames []string, capacity *HostGpuCapacity, override GpuOverride) *GpuTargetCheckResult {
	if len(serviceNames) == 0 {
		return nil
	}

	eff := EffectiveGpu(capacity, override)
	if eff.HasGpu {
		return &GpuTargetCheckResult{Status: "pass"}
	}

	quoted := make([]string, len(serviceNames))
	for i, s := range serviceNames {
		quoted[i] = fmt.Sprintf("%q", s)
	}
	who := strings.Join(quoted, ", ")
	verb := "ask"
	if len(serviceNames) == 1 {
		verb = "asks"
	}

	if override == GpuOverrideNo {
		return &GpuTargetCheckResult{
			Status: "warn",
			Message: fmt.Sprintf("%s %s for a GPU, but this server is marked \"no GPU\". Pick a GPU server, or clear the mark on the server if that is wrong. You can continue; the container will start without a GPU or fail to start.",
				who, verb),
		}
	}

	if capacity == nil {
		return nil
	}

	return &GpuTargetCheckResult{
		Status: "warn",
		Message: fmt.Sprintf("%s %s for a GPU, but no NVIDIA runtime was detected on this server. Pick a GPU server, or mark this server as having a GPU if it does (an older Docker daemon cannot report it). You can continue; Docker will refuse to start the container if there is no GPU.",
			who, verb),
	}
}
