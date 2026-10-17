package utils

import (
	"reflect"
	"testing"
)

func TestInspectCommandEnvironment(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantVars   map[string]string
		wantKeys   []string
		wantCmd    string
		wantMigrat bool
	}{
		{
			name:       "simple single prefix",
			input:      "NODE_ENV=production node server.js",
			wantVars:   map[string]string{"NODE_ENV": "production"},
			wantKeys:   []string{"NODE_ENV"},
			wantCmd:    "node server.js",
			wantMigrat: true,
		},
		{
			name:       "cross-env prefix with multiple quotes",
			input:      `cross-env PORT=8080 API_KEY="secret-123" npm start`,
			wantVars:   map[string]string{"PORT": "8080", "API_KEY": "secret-123"},
			wantKeys:   []string{"PORT", "API_KEY"},
			wantCmd:    "npm start",
			wantMigrat: true,
		},
		{
			name:       "export syntax with &&",
			input:      "export FOO=bar BAZ=qux && ./run.sh",
			wantVars:   map[string]string{"FOO": "bar", "BAZ": "qux"},
			wantKeys:   []string{"FOO", "BAZ"},
			wantCmd:    "./run.sh",
			wantMigrat: true,
		},
		{
			name:       "plain command with no vars",
			input:      "npm run build",
			wantVars:   map[string]string{},
			wantKeys:   nil,
			wantCmd:    "npm run build",
			wantMigrat: false,
		},
		{
			name:       "chained assignments require manual edit",
			input:      "VAR1=val1 cmd1 && VAR2=val2 cmd2",
			wantVars:   map[string]string{"VAR1": "val1"},
			wantKeys:   []string{"VAR1", "VAR2"},
			wantCmd:    "cmd1 && VAR2=val2 cmd2",
			wantMigrat: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InspectCommandEnvironment(tt.input)
			if !reflect.DeepEqual(got.Variables, tt.wantVars) {
				t.Errorf("Variables = %+v, want %+v", got.Variables, tt.wantVars)
			}
			if len(got.Keys) != len(tt.wantKeys) {
				t.Errorf("Keys = %+v, want %+v", got.Keys, tt.wantKeys)
			}
			if got.Command != tt.wantCmd {
				t.Errorf("Command = %q, want %q", got.Command, tt.wantCmd)
			}
			if got.Migratable != tt.wantMigrat {
				t.Errorf("Migratable = %v, want %v", got.Migratable, tt.wantMigrat)
			}
		})
	}
}

func TestCommandEnvironmentWarning(t *testing.T) {
	w := CommandEnvironmentWarning("PORT=3000 npm start", "Start command")
	if w == "" {
		t.Fatal("expected warning, got empty")
	}
	if !reflect.DeepEqual(CommandEnvironmentWarning("npm start", "Start command"), "") {
		t.Fatal("expected no warning for plain command")
	}
}
