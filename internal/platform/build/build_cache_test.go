package build

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestBuildStorageReserve(t *testing.T) {
	t.Run("small disk calculates 10 percent", func(t *testing.T) {
		storage := BuildStorage{TotalBytes: 10 * 1024 * 1024 * 1024} // 10 GiB
		reserve := BuildStorageReserve(storage)
		expected := uint64(1 * 1024 * 1024 * 1024) // 1 GiB
		if reserve != expected {
			t.Fatalf("expected %d, got %d", expected, reserve)
		}
	})

	t.Run("large disk capped at 2 GiB", func(t *testing.T) {
		storage := BuildStorage{TotalBytes: 100 * 1024 * 1024 * 1024} // 100 GiB
		reserve := BuildStorageReserve(storage)
		if reserve != MaxBuildStorageReserveBytes {
			t.Fatalf("expected %d, got %d", MaxBuildStorageReserveBytes, reserve)
		}
	})
}

func TestAssertBuildStorage(t *testing.T) {
	storage := BuildStorage{
		TotalBytes: 50 * 1024 * 1024 * 1024,
		FreeBytes:  3 * 1024 * 1024 * 1024, // 3 GiB free (reserve is 2 GiB)
	}

	if err := AssertBuildStorage(storage); err != nil {
		t.Fatalf("expected healthy storage to pass, got %v", err)
	}

	storage.FreeBytes = 1 * 1024 * 1024 * 1024 // 1 GiB free (below 2 GiB reserve)
	if err := AssertBuildStorage(storage); err != ErrBuildStorageLow {
		t.Fatalf("expected ErrBuildStorageLow, got %v", err)
	}

	storage.FreeBytes = 3 * 1024 * 1024 * 1024
	zeroInodes := int64(0)
	storage.FreeInodes = &zeroInodes
	if err := AssertBuildStorage(storage); err != ErrBuildStorageLow {
		t.Fatalf("expected ErrBuildStorageLow on exhausted inodes, got %v", err)
	}
}

func TestParseBuildStorage(t *testing.T) {
	raw := `
CODEDOCK_STORAGE 52428800 20971520
CODEDOCK_INODES 3500000
`
	storage, err := ParseBuildStorage(raw)
	if err != nil {
		t.Fatalf("failed to parse build storage: %v", err)
	}
	if storage.TotalBytes != 52428800*1024 {
		t.Fatalf("unexpected totalBytes: %d", storage.TotalBytes)
	}
	if storage.FreeBytes != 20971520*1024 {
		t.Fatalf("unexpected freeBytes: %d", storage.FreeBytes)
	}
	if storage.FreeInodes == nil || *storage.FreeInodes != 3500000 {
		t.Fatalf("unexpected freeInodes: %v", storage.FreeInodes)
	}

	codedockRaw := `
CODEDOCK_STORAGE 10485760 5242880
CODEDOCK_INODES 1200000
`
	storage, err = ParseBuildStorage(codedockRaw)
	if err != nil {
		t.Fatalf("failed to parse codedock storage: %v", err)
	}
	if storage.TotalBytes != 10485760*1024 {
		t.Fatalf("unexpected totalBytes: %d", storage.TotalBytes)
	}

	btrfsRaw := `
CODEDOCK_STORAGE 524288000 301989888
CODEDOCK_INODES 0 0
`
	storage, err = ParseBuildStorage(btrfsRaw)
	if err != nil {
		t.Fatalf("failed to parse btrfs storage: %v", err)
	}
	if storage.FreeInodes != nil {
		t.Fatalf("expected nil freeInodes for btrfs 0 0, got %v", storage.FreeInodes)
	}
	if err := AssertBuildStorage(storage); err != nil {
		t.Fatalf("expected btrfs storage to pass build assertion, got %v", err)
	}

	invalid := `not a valid storage output`
	if _, err := ParseBuildStorage(invalid); err != ErrBuildStorageUnknown {
		t.Fatalf("expected ErrBuildStorageUnknown, got %v", err)
	}
}

func TestIsBuildCacheMountsEnabled(t *testing.T) {
	if !IsBuildCacheMountsEnabled(nil) {
		t.Fatal("expected default to be enabled")
	}

	os.Setenv("CODEDOCK_BUILD_CACHE_MOUNTS", "false")
	defer os.Unsetenv("CODEDOCK_BUILD_CACHE_MOUNTS")
	if IsBuildCacheMountsEnabled(nil) {
		t.Fatal("expected disabled when CODEDOCK_BUILD_CACHE_MOUNTS=false")
	}

	os.Unsetenv("CODEDOCK_BUILD_CACHE_MOUNTS")
	if !IsBuildCacheMountsEnabled(map[string]string{"FOO": "bar"}) {
		t.Fatal("expected enabled with generic env vars")
	}

	if IsBuildCacheMountsEnabled(map[string]string{"CODEDOCK_BUILD_CACHE_MOUNTS": "false"}) {
		t.Fatal("expected disabled when env map specifies false")
	}
}

func TestGetPackageCacheDirs(t *testing.T) {
	npmDirs := GetPackageCacheDirs("npm")
	if len(npmDirs) != 1 || npmDirs[0] != "/root/.npm" {
		t.Fatalf("unexpected npm dirs: %v", npmDirs)
	}

	goDirs := GetPackageCacheDirs("go")
	if len(goDirs) != 2 || goDirs[0] != "/go/pkg/mod" {
		t.Fatalf("unexpected go dirs: %v", goDirs)
	}

	if GetPackageCacheDirs("unknown") != nil {
		t.Fatal("expected nil for unknown package manager")
	}
}

func TestFormatCacheMount(t *testing.T) {
	mount, err := FormatCacheMount("test-id", "/root/.cache/pip")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "--mount=type=cache,id=codedock-test-id-${TARGETPLATFORM},target=/root/.cache/pip,sharing=locked"
	if mount != expected {
		t.Fatalf("expected %s, got %s", expected, mount)
	}

	if _, err := FormatCacheMount("test-id", "relative/path"); err == nil {
		t.Fatal("expected error for relative path")
	}

	if _, err := FormatCacheMount("test-id", "/root/unsafe;rm -rf /"); err == nil {
		t.Fatal("expected error for unsafe path with semicolon")
	}
}

func TestBuildCacheKeepBytes(t *testing.T) {
	const gib = 1024 * 1024 * 1024
	storage := BuildStorage{
		TotalBytes: 100 * gib,
		FreeBytes:  50 * gib,
	}

	t.Run("automatic budget calculates 20 percent of disk", func(t *testing.T) {
		keep, err := BuildCacheKeepBytes(storage, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := uint64(20 * gib)
		if keep != expected {
			t.Fatalf("expected %d, got %d", expected, keep)
		}
	})

	t.Run("custom valid budget within disk bounds", func(t *testing.T) {
		budget := int64(24 * gib)
		keep, err := BuildCacheKeepBytes(storage, &budget)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if keep != uint64(budget) {
			t.Fatalf("expected %d, got %d", budget, keep)
		}
	})

	t.Run("custom budget exceeding disk minus reserve is clamped", func(t *testing.T) {
		budget := int64(500 * gib)
		keep, err := BuildCacheKeepBytes(storage, &budget)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Reserve for 100 GiB is 2 GiB (min(2 GiB, 10 GiB)).
		// Max budget is 100 - 2 = 98 GiB.
		expected := uint64(98 * gib)
		if keep != expected {
			t.Fatalf("expected %d, got %d", expected, keep)
		}
	})

	t.Run("invalid budget returns error", func(t *testing.T) {
		invalid := int64(0)
		if _, err := BuildCacheKeepBytes(storage, &invalid); err == nil {
			t.Fatal("expected error for budget <= 0")
		}
		negative := int64(-10)
		if _, err := BuildCacheKeepBytes(storage, &negative); err == nil {
			t.Fatal("expected error for negative budget")
		}
	})
}

func TestPressureAdjustedKeepBytes(t *testing.T) {
	const gib = 1024 * 1024 * 1024
	storage := BuildStorage{
		TotalBytes: 100 * gib,
		FreeBytes:  50 * gib, // Target is 20 GiB (20%), so plenty of free space
	}

	t.Run("unconstrained when plenty of free space", func(t *testing.T) {
		budget := uint64(24 * gib)
		adjusted := PressureAdjustedKeepBytes(storage, 15*gib, budget)
		if adjusted != budget {
			t.Fatalf("expected %d, got %d", budget, adjusted)
		}
	})

	t.Run("shrinks under disk pressure towards floor", func(t *testing.T) {
		underPressure := BuildStorage{
			TotalBytes: 100 * gib,
			FreeBytes:  10 * gib, // Target is 20 GiB, 10 GiB below target
		}
		// Floor is min(5 GiB, 10 GiB) = 5 GiB.
		// idle (8 GiB) + free (10 GiB) - target (20 GiB) = -2 GiB -> clamped to floor 5 GiB.
		adjusted := PressureAdjustedKeepBytes(underPressure, 8*gib, 24*gib)
		expected := uint64(5 * gib)
		if adjusted != expected {
			t.Fatalf("expected %d, got %d", expected, adjusted)
		}
	})
}

func TestBuildKitCacheConfig(t *testing.T) {
	const gib = 1024 * 1024 * 1024
	storage := BuildStorage{
		TotalBytes: 100 * gib,
		FreeBytes:  50 * gib,
	}

	cfg, err := BuildKitCacheConfig(storage, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(cfg, fmt.Sprintf("maxUsedSpace = %d", 20*gib)) {
		t.Fatalf("expected maxUsedSpace in config, got:\n%s", cfg)
	}
	if !strings.Contains(cfg, `minFreeSpace = "20%"`) {
		t.Fatalf("expected minFreeSpace in config, got:\n%s", cfg)
	}
}
