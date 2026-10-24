package utils

import (
	"testing"
)

func TestValidateEventTypes(t *testing.T) {
	validated, err := ValidateEventTypes([]string{" deploy ", "rollback"})
	if err != nil {
		t.Fatal(err)
	}
	if len(validated) != 2 || validated[0] != "deploy" || validated[1] != "rollback" {
		t.Fatalf("expected trimmed types, got %v", validated)
	}
	if _, err := ValidateEventTypes(nil); err != nil {
		t.Fatalf("expected nil to pass through, got %v", err)
	}
	for _, input := range [][]string{{""}, {"  "}, {"deploy,rollback"}} {
		if _, err := ValidateEventTypes(input); err == nil {
			t.Fatalf("expected %v to be rejected", input)
		}
	}
}
