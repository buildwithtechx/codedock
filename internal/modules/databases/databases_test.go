package databases

import (
	"testing"
)

func TestBuildSQLMutations(t *testing.T) {
	ins, err := buildInsertSQL("postgres", "users", map[string]any{"name": "O'Reilly"})
	if err != nil || ins == "" {
		t.Fatalf("expected valid insert SQL, got %q, err %v", ins, err)
	}

	upd, err := buildUpdateSQL("postgres", "users", map[string]any{"id": 1}, map[string]any{"name": "Alice"})
	if err != nil || upd == "" {
		t.Fatalf("expected valid update SQL, got %q, err %v", upd, err)
	}

	del, err := buildDeleteSQL("mysql", "users", map[string]any{"id": 1})
	if err != nil || del == "" {
		t.Fatalf("expected valid delete SQL, got %q, err %v", del, err)
	}

	if _, err := buildDeleteSQL("postgres", "users", nil); err == nil {
		t.Fatal("expected error when keys are empty")
	}
}
