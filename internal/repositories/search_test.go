package repositories

import "testing"

func TestEscapeLikeTerm(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{name: "plain", input: "billing", want: "billing"},
		{name: "percent", input: "100%", want: `100\%`},
		{name: "underscore", input: "api_gateway", want: `api\_gateway`},
		{name: "backslash", input: `a\b`, want: `a\\b`},
		{name: "all wildcards", input: `%_%\`, want: `\%\_\%\\`},
		{name: "empty", input: "", want: ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := EscapeLikeTerm(tc.input); got != tc.want {
				t.Fatalf("EscapeLikeTerm(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestContainsLikePattern(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{name: "plain", input: "billing", want: "%billing%"},
		{name: "trims", input: "  billing  ", want: "%billing%"},
		{name: "wildcard literal", input: "100%_off", want: `%100\%\_off%`},
		{name: "empty", input: "   ", want: "%%"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ContainsLikePattern(tc.input); got != tc.want {
				t.Fatalf("ContainsLikePattern(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
