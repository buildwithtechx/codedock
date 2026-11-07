package projects

import "testing"

func TestNormalizeRuntimeVersion(t *testing.T) {
	cases := []struct {
		name    string
		raw     string
		want    string
		wantErr bool
	}{
		{"empty passes through", "", "", false},
		{"trims whitespace", "  8.4 ", "8.4", false},
		{"node semver", "22.11.0", "22.11.0", false},
		{"rejects spaces inside", "8.4 beta", "", true},
		{"rejects shell injection", "8.4; rm -rf /", "", true},
		{"rejects leading dash", "-8.4", "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := normalizeRuntimeVersion(tc.raw)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error for %q", tc.raw)
			}
			if !tc.wantErr {
				if err != nil {
					t.Fatalf("unexpected error for %q: %v", tc.raw, err)
				}
				if got != tc.want {
					t.Fatalf("got %q, want %q", got, tc.want)
				}
			}
		})
	}
}
