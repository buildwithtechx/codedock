package deployments

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRepositoryInspectionDetectsFrameworkAndPackageManager(t *testing.T) {
	directory := t.TempDir()
	for name, value := range map[string]string{"package.json": `{"scripts":{"build":"vite build"},"devDependencies":{"vite":"1"}}`, "pnpm-lock.yaml": "lockfileVersion: 9"} {
		if err := os.WriteFile(filepath.Join(directory, name), []byte(value), 0600); err != nil {
			t.Fatal(err)
		}
	}
	result, err := inspectRepositoryFiles(directory)
	if err != nil || result.Framework != "Vite" || result.PackageManager != "pnpm" || result.BuildCommand != "pnpm run build" || result.StaticOutput != "dist" || result.InternalPort != 80 {
		t.Fatalf("incorrect defaults: %+v %v", result, err)
	}
}

func TestInspectionRootRejectsCheckoutEscape(t *testing.T) {
	directory := t.TempDir()
	if _, err := inspectionRoot(directory, "../"); err == nil {
		t.Fatal("root escaped checkout")
	}
}
