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
