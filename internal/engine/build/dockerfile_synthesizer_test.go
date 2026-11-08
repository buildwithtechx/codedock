package build

import (
	"strings"
	"testing"

	"codedock/internal/models"
)

func TestSynthesizeDockerfileStaticSPAFallback(t *testing.T) {
	trueVal := true
	falseVal := false

	t.Run("default (nil) fallback renders index.html", func(t *testing.T) {
		res, err := synthesizeDockerfile("Static", BuildOptions{
			AppConfig: &models.AppService{},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(res, "try_files $uri $uri/ /index.html;") {
			t.Errorf("expected default static dockerfile to fallback to /index.html, got:\n%s", res)
		}
	})

	t.Run("explicit true fallback renders index.html", func(t *testing.T) {
		res, err := synthesizeDockerfile("Static", BuildOptions{
			AppConfig: &models.AppService{
				SPAFallback: &trueVal,
			},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(res, "try_files $uri $uri/ /index.html;") {
			t.Errorf("expected explicit true fallback to render /index.html, got:\n%s", res)
		}
	})

	t.Run("explicit false fallback renders real 404 and 404.html", func(t *testing.T) {
		res, err := synthesizeDockerfile("Static", BuildOptions{
			AppConfig: &models.AppService{
				SPAFallback: &falseVal,
			},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(res, "try_files $uri $uri/ =404;") {
			t.Errorf("expected false fallback to render =404, got:\n%s", res)
		}
		if !strings.Contains(res, "error_page 404 /404.html;") {
			t.Errorf("expected false fallback to configure error_page 404 /404.html, got:\n%s", res)
		}
	})
}

func TestReinterpolateCommand(t *testing.T) {
	env := map[string]string{
		"PORT": "8080",
		"HOST": "0.0.0.0",
		"MODE": "production",
	}

	tests := []struct {
		name     string
		cmd      string
		expected string
	}{
		{
			name:     "brace syntax",
			cmd:      "serve -p ${PORT} -h ${HOST}",
			expected: "serve -p 8080 -h 0.0.0.0",
		},
		{
			name:     "dollar syntax",
			cmd:      "node server.js --port=$PORT --env=$MODE",
			expected: "node server.js --port=8080 --env=production",
		},
		{
			name:     "unset variable kept intact",
			cmd:      "echo ${UNSET_VAR}",
			expected: "echo ${UNSET_VAR}",
		},
		{
			name:     "empty command",
			cmd:      "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReinterpolateCommand(tt.cmd, env)
			if got != tt.expected {
				t.Errorf("ReinterpolateCommand(%q) = %q, want %q", tt.cmd, got, tt.expected)
			}
		})
	}
}

func TestResolveBunToolchain(t *testing.T) {
	t.Run("requires exact stable version", func(t *testing.T) {
		_, err := ResolveBunToolchain("^1.4.4", "")
		if err == nil {
			t.Fatal("expected error for non-exact semver")
		}
	})

	t.Run("refuses downgrade when repo requires newer version", func(t *testing.T) {
		_, err := ResolveBunToolchain("1.4.4", "bun@1.5.0")
		if err == nil {
			t.Fatal("expected error for downgrade attempt")
		}
		if !strings.Contains(err.Error(), "the build will not downgrade bun") {
			t.Errorf("expected downgrade error message, got: %v", err)
		}
	})

	t.Run("refuses discarding integrity requirement", func(t *testing.T) {
		_, err := ResolveBunToolchain("1.4.4", "bun@1.3.1+sha256:abcd1234ef5678")
		if err == nil {
			t.Fatal("expected error for discarding integrity requirement")
		}
		if !strings.Contains(err.Error(), "integrity requirement") {
			t.Errorf("expected integrity error message, got: %v", err)
		}
	})

	t.Run("allows valid toolchain override upgrade", func(t *testing.T) {
		pin, err := ResolveBunToolchain("1.4.4", "bun@1.3.1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if pin != "bun@1.4.4" {
			t.Errorf("pin = %q, want bun@1.4.4", pin)
		}
	})
}

func TestSynthesizeDockerfileBunToolchainAndImageOverride(t *testing.T) {
	t.Run("replaces FROM with BuildImage and synthesizes bun Dockerfile", func(t *testing.T) {
		res, err := synthesizeDockerfile("Bun", BuildOptions{
			AppConfig: &models.AppService{
				BuildImage: "oven/bun:1.4.4-alpine",
			},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.HasPrefix(res, "FROM oven/bun:1.4.4-alpine\n") {
			t.Errorf("expected Dockerfile to start with custom BuildImage, got:\n%s", res)
		}
		if !strings.Contains(res, "bun install --frozen-lockfile") {
			t.Errorf("expected bun Dockerfile template, got:\n%s", res)
		}
	})

	t.Run("rejects synthesis if runtime image downgrades repository pin", func(t *testing.T) {
		_, err := synthesizeDockerfile("Bun", BuildOptions{
			AppConfig: &models.AppService{
				BuildImage: "oven/bun:1.4.4-alpine",
			},
			EnvVars: map[string]string{
				"PACKAGE_MANAGER": "bun@1.5.0",
			},
		})
		if err == nil {
			t.Fatal("expected error for build downgrade attempt")
		}
		if !strings.Contains(err.Error(), "the build will not downgrade bun") {
			t.Errorf("expected downgrade error, got: %v", err)
		}
	})
}
