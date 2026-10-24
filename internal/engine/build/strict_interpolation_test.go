package build

import "testing"

func TestStrictBindingsResolveChainsAndRejectMissingOrCyclicSources(t *testing.T) {
	registry := map[string]map[string]string{"database": {"URL": "postgres://db/app"}, "api": {"DATABASE_URL": "${database.URL}"}}
	values, err := InterpolateEnvVarsStrict(map[string]string{"URL": "${api.DATABASE_URL}"}, registry)
	if err != nil || values["URL"] != "postgres://db/app" {
		t.Fatalf("binding chain failed: %v", err)
	}
	if _, err := InterpolateEnvVarsStrict(map[string]string{"URL": "${other.URL}"}, registry); err == nil {
		t.Fatal("missing source silently retained")
	}
	registry["database"]["URL"] = "${api.DATABASE_URL}"
	if _, err := InterpolateEnvVarsStrict(map[string]string{"URL": "${api.DATABASE_URL}"}, registry); err == nil {
		t.Fatal("cyclic binding accepted")
	}
}
