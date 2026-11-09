package compose

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateStackResults(t *testing.T) {
	config := `{"services":{"web":{},"worker":{}}}`
	results := []ComposeServiceResult{
		{Name: "web", State: "running", Health: "healthy"},
		{Name: "worker", State: "running", Health: ""},
	}
	assert.NoError(t, validateStackResults(config, results))

	unhealthyResults := []ComposeServiceResult{
		{Name: "web", State: "running", Health: "unhealthy"},
		{Name: "worker", State: "running", Health: ""},
	}
	assert.Error(t, validateStackResults(config, unhealthyResults))
}
