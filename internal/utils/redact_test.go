package utils

import "testing"

func TestHasRedactedConfiguration(t *testing.T) {
	cases := []struct {
		name     string
		values   []string
		expected bool
	}{
		{
			name:     "valid clean configuration",
			values:   []string{"DATABASE_URL=postgres://user:pass@db:5432/main", "LOG_LEVEL=info"},
			expected: false,
		},
		{
			name:     "scrubbed dsn with asterisks",
			values:   []string{"postgres://user:***@localhost:5432/mydb"},
			expected: true,
		},
		{
			name:     "bare asterisk redaction",
			values:   []string{"***"},
			expected: true,
		},
		{
			name:     "quoted secret mask",
			values:   []string{`SECRET="*****"`},
			expected: true,
		},
		{
			name:     "assignment mask",
			values:   []string{"password: ***"},
			expected: true,
		},
		{
			name:     "empty values",
			values:   []string{"", "   "},
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := HasRedactedConfiguration(tc.values...)
			if got != tc.expected {
				t.Errorf("HasRedactedConfiguration(%v) = %v; want %v", tc.values, got, tc.expected)
			}
		})
	}
}
