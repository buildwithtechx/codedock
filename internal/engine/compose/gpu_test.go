package compose

import (
	"testing"

	"codedock/internal/models"
)

func TestIsNvidiaGpuID(t *testing.T) {
	valid := []string{
		"all",
		"0",
		"1",
		"0:1",
		"GPU-1234abcd-56ef-7890",
		"MIG-GPU-1234abcd/0/1",
		"nvidia.com/gpu=0",
		"nvidia.com/gpu=GPU-abc-123",
	}
	for _, id := range valid {
		if !IsNvidiaGpuID(id) {
			t.Errorf("expected %q to be valid NVIDIA GPU ID", id)
		}
	}

	invalid := []string{
		"",
		"amd-gpu-1",
		"tpu-0",
		"GPU:invalid:syntax",
		"nvidia.com/other=0",
	}
	for _, id := range invalid {
		if IsNvidiaGpuID(id) {
			t.Errorf("expected %q to be invalid NVIDIA GPU ID", id)
		}
	}
}

func TestNormalizeGPURequests(t *testing.T) {
	t.Run("defaults to all with gpu capability", func(t *testing.T) {
		reqs := []models.GPURequest{{}}
		got, err := NormalizeGPURequests(reqs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("expected 1 request, got %d", len(got))
		}
		if got[0].Driver != "nvidia" || got[0].Count != "all" || len(got[0].Capabilities) != 1 || got[0].Capabilities[0] != "gpu" {
			t.Errorf("unexpected normalized request: %+v", got[0])
		}
	})

	t.Run("specific count", func(t *testing.T) {
		reqs := []models.GPURequest{{Count: 2}}
		got, err := NormalizeGPURequests(reqs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got[0].Count != 2 {
			t.Errorf("expected count 2, got %v", got[0].Count)
		}
	})

	t.Run("rejects count and device_ids combined", func(t *testing.T) {
		reqs := []models.GPURequest{{Count: 1, DeviceIDs: []string{"0"}}}
		_, err := NormalizeGPURequests(reqs)
		if err == nil {
			t.Fatalf("expected error combining count and device_ids")
		}
	})

	t.Run("rejects non-nvidia driver", func(t *testing.T) {
		reqs := []models.GPURequest{{Driver: "rocm"}}
		_, err := NormalizeGPURequests(reqs)
		if err == nil {
			t.Fatalf("expected error for non-nvidia driver")
		}
	})

	t.Run("specific device IDs", func(t *testing.T) {
		reqs := []models.GPURequest{{DeviceIDs: []string{"0", "GPU-test-1"}}}
		got, err := NormalizeGPURequests(reqs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got[0].DeviceIDs) != 2 {
			t.Errorf("expected 2 device IDs, got %d", len(got[0].DeviceIDs))
		}
	})
}

func TestToDockerDeviceRequests(t *testing.T) {
	t.Run("translates count all to -1", func(t *testing.T) {
		reqs := []models.GPURequest{{Count: "all"}}
		dockerReqs, err := ToDockerDeviceRequests(reqs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(dockerReqs) != 1 {
			t.Fatalf("expected 1 docker request, got %d", len(dockerReqs))
		}
		if dockerReqs[0].Driver != "nvidia" || dockerReqs[0].Count != -1 {
			t.Errorf("expected Driver nvidia and Count -1, got %+v", dockerReqs[0])
		}
	})

	t.Run("translates specific devices", func(t *testing.T) {
		reqs := []models.GPURequest{{DeviceIDs: []string{"0"}}}
		dockerReqs, err := ToDockerDeviceRequests(reqs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(dockerReqs[0].DeviceIDs) != 1 || dockerReqs[0].DeviceIDs[0] != "0" {
			t.Errorf("expected DeviceID 0, got %+v", dockerReqs[0].DeviceIDs)
		}
	})
}
