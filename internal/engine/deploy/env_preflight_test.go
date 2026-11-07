package deploy

import (
	"strings"
	"testing"
)

func TestFindMissingRequiredEnv(t *testing.T) {
	raw := map[string]string{
		"DATABASE_URL": "${DB_URL:?database connection is mandatory}",
		"REDIS_URL":    "${REDIS_HOST?must be set}",
		"PORT":         "${PORT:-8080}",
		"OPTIONAL":     "${EXTRA}",
	}

	// 1. When DB_URL and REDIS_HOST are missing
	resolved := map[string]string{
		"PORT": "8080",
	}

	missing := FindMissingRequiredEnv("web", raw, resolved)
	if len(missing) != 2 {
		t.Fatalf("expected 2 missing variables, got %d", len(missing))
	}

	err := FormatRequiredEnvError(missing)
	if err == nil {
		t.Fatal("expected formatted error, got nil")
	}

	errMsg := err.Error()
	expectedPrefix := "Required Compose environment variables are not configured. Service \"web\":"
	if !strings.HasPrefix(errMsg, expectedPrefix) {
		t.Fatalf("expected prefix %q, got %q", expectedPrefix, errMsg)
	}
	if !strings.Contains(errMsg, "DB_URL") || !strings.Contains(errMsg, "REDIS_HOST") {
		t.Fatalf("expected DB_URL and REDIS_HOST in message, got: %s", errMsg)
	}

	// 2. When DB_URL is empty with :?
	resolved["DB_URL"] = "   "
	resolved["REDIS_HOST"] = "localhost"
	missing2 := FindMissingRequiredEnv("web", raw, resolved)
	if len(missing2) != 1 || missing2[0].Variable != "DB_URL" {
		t.Fatalf("expected DB_URL to be missing due to empty value with :?, got %v", missing2)
	}

	// 3. When both are satisfied
	resolved["DB_URL"] = "postgres://user:pass@localhost:5432/db"
	missing3 := FindMissingRequiredEnv("web", raw, resolved)
	if len(missing3) != 0 {
		t.Fatalf("expected 0 missing variables, got %d", len(missing3))
	}
	if err := FormatRequiredEnvError(missing3); err != nil {
		t.Fatalf("expected nil error when none missing, got %v", err)
	}
}

func TestFormatRequiredEnvErrorMultipleServices(t *testing.T) {
	missing := []MissingEnvVariable{
		{ServiceName: "worker", Variable: "QUEUE_URL"},
		{ServiceName: "api", Variable: "SECRET_KEY"},
		{ServiceName: "api", Variable: "DATABASE_URL"},
	}

	err := FormatRequiredEnvError(missing)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expected := `Required Compose environment variables are not configured. Service "api": DATABASE_URL, SECRET_KEY; Service "worker": QUEUE_URL`
	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}
}
