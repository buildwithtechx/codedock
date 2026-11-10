package stacks

import "testing"

func TestResolveStackRolesExplicitConfigWins(t *testing.T) {
	explicit := []StackRole{{Name: "custom", Kind: RoleKindWorker, Command: "run"}}
	got := ResolveStackRoles(&explicit, Registry["laravel"].DefaultRoles, []string{"laravel/horizon"}, nil)
	if len(got) != 1 || got[0].Name != "custom" {
		t.Fatalf("explicit roles should win, got %#v", got)
	}
	empty := []StackRole{}
	if got := ResolveStackRoles(&empty, Registry["laravel"].DefaultRoles, nil, nil); len(got) != 0 {
		t.Fatalf("explicit empty should win, got %#v", got)
	}
}

func TestResolveStackRolesLaravelHorizonReplacesQueue(t *testing.T) {
	laravel := Registry["laravel"].DefaultRoles

	without := roleNames(ResolveStackRoles(nil, laravel, nil, nil))
	if !contains(without, "queue") || contains(without, "horizon") || !contains(without, "scheduler") {
		t.Fatalf("expected queue + scheduler without horizon, got %v", without)
	}

	with := roleNames(ResolveStackRoles(nil, laravel, []string{"laravel/horizon"}, nil))
	if contains(with, "queue") || !contains(with, "horizon") || !contains(with, "scheduler") {
		t.Fatalf("expected horizon + scheduler instead of queue, got %v", with)
	}
}

func TestResolveStackRolesDropsAmbiguousKind(t *testing.T) {
	got := ResolveStackRoles(nil, Registry["rails"].DefaultRoles, []string{"solid_queue", "sidekiq"}, []string{"bin/jobs"})
	if len(got) != 0 {
		t.Fatalf("two workers match, both should be dropped, got %#v", got)
	}
}

func TestResolveStackRolesGates(t *testing.T) {
	rails := Registry["rails"].DefaultRoles

	got := roleNames(ResolveStackRoles(nil, rails, []string{"sidekiq"}, nil))
	if len(got) != 1 || got[0] != "sidekiq" {
		t.Fatalf("expected only sidekiq, got %v", got)
	}
	if got := ResolveStackRoles(nil, rails, []string{"solid_queue"}, nil); len(got) != 0 {
		t.Fatalf("jobs requires bin/jobs too, got %#v", got)
	}
	got = roleNames(ResolveStackRoles(nil, rails, []string{"solid_queue"}, []string{"bin/jobs"}))
	if len(got) != 1 || got[0] != "jobs" {
		t.Fatalf("expected jobs with both gate signals, got %v", got)
	}
}

func TestResolveWebCommand(t *testing.T) {
	if got := ResolveWebCommand("serve", false, "php artisan optimize"); got != "serve" {
		t.Errorf("optimize off should keep command, got %q", got)
	}
	if got := ResolveWebCommand("serve", true, "php artisan optimize"); got != "php artisan optimize && serve" {
		t.Errorf("optimize on should prefix, got %q", got)
	}
	if got := ResolveWebCommand("serve", true, ""); got != "serve" {
		t.Errorf("missing optimize command should keep command, got %q", got)
	}
}

func roleNames(roles []StackRole) []string {
	names := make([]string, len(roles))
	for i, role := range roles {
		names[i] = role.Name
	}
	return names
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
