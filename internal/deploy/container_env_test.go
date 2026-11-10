package deploy

import (
	"reflect"
	"strings"
	"testing"
)

func TestWorkerContainerUnsetsPortEnv(t *testing.T) {
	inputEnvs := []string{
		"APP_ENV=production",
		"PORT=8080",
		"DATABASE_URL=postgres://...",
		"PORT_OFFSET=10",
	}

	filtered := make([]string, 0, len(inputEnvs)+1)
	for _, e := range inputEnvs {
		if !strings.HasPrefix(e, "PORT=") {
			filtered = append(filtered, e)
		}
	}
	filtered = append(filtered, "PORT")

	expected := []string{
		"APP_ENV=production",
		"DATABASE_URL=postgres://...",
		"PORT_OFFSET=10",
		"PORT",
	}

	if !reflect.DeepEqual(filtered, expected) {
		t.Errorf("got %v, want %v", filtered, expected)
	}
}
