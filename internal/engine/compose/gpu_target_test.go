package compose

import (
	"strings"
	"testing"
)

func TestEffectiveGpu(t *testing.T) {
	one := 1
	two := 2

	t.Run("probe detected with runtime", func(t *testing.T) {
		res := EffectiveGpu(&HostGpuCapacity{GpuRuntime: true}, "")
		if !res.HasGpu || res.Source != "probe" {
			t.Fatalf("expected GPU detected from probe, got %+v", res)
		}
	})

	t.Run("probe detected with count", func(t *testing.T) {
		res := EffectiveGpu(&HostGpuCapacity{GpuCount: &two}, "")
		if !res.HasGpu || res.Source != "probe" || *res.Count != 2 {
			t.Fatalf("expected GPU detected with count 2, got %+v", res)
		}
	})

	t.Run("override yes wins over no probe", func(t *testing.T) {
		res := EffectiveGpu(&HostGpuCapacity{GpuRuntime: false}, GpuOverrideYes)
		if !res.HasGpu || res.Source != "override" {
			t.Fatalf("expected GPU true via override, got %+v", res)
		}
	})

	t.Run("override no hides card", func(t *testing.T) {
		res := EffectiveGpu(&HostGpuCapacity{GpuRuntime: true, GpuCount: &one}, GpuOverrideNo)
		if res.HasGpu || res.Source != "override" || res.Count != nil {
			t.Fatalf("expected GPU hidden via override no, got %+v", res)
		}
	})

	t.Run("unknown capacity without override", func(t *testing.T) {
		res := EffectiveGpu(nil, "")
		if res.HasGpu || res.Source != "none" {
			t.Fatalf("expected no GPU, got %+v", res)
		}
	})
}

func TestCheckGpuTarget(t *testing.T) {
	one := 1

	t.Run("no gpu services returns nil", func(t *testing.T) {
		if res := CheckGpuTarget(nil, nil, ""); res != nil {
			t.Fatalf("expected nil check, got %+v", res)
		}
	})

	t.Run("has gpu passes", func(t *testing.T) {
		res := CheckGpuTarget([]string{"ollama"}, &HostGpuCapacity{GpuCount: &one}, "")
		if res == nil || res.Status != "pass" {
			t.Fatalf("expected pass status, got %+v", res)
		}
	})

	t.Run("server marked no gpu warns", func(t *testing.T) {
		res := CheckGpuTarget([]string{"ollama"}, &HostGpuCapacity{GpuCount: &one}, GpuOverrideNo)
		if res == nil || res.Status != "warn" || !strings.Contains(res.Message, "marked \"no GPU\"") {
			t.Fatalf("expected warn for marked no gpu, got %+v", res)
		}
	})

	t.Run("no nvidia runtime warns", func(t *testing.T) {
		res := CheckGpuTarget([]string{"ollama", "vllm"}, &HostGpuCapacity{GpuRuntime: false}, "")
		if res == nil || res.Status != "warn" || !strings.Contains(res.Message, "no NVIDIA runtime was detected") {
			t.Fatalf("expected warn for missing runtime, got %+v", res)
		}
	})
}
