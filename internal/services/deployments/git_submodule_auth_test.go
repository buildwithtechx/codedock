package deployments

import (
	"testing"
)

func TestParseSameOwnerSubmoduleRepos(t *testing.T) {
	gitmodules := `
[submodule "shared-lib"]
	path = lib/shared
	url = ../shared-lib.git
[submodule "nested"]
	path = sub/nested
	url = https://github.com/my-org/nested-repo.git
[submodule "ssh-repo"]
	path = sub/ssh
	url = git@github.com:my-org/ssh-repo.git
[submodule "different-owner"]
	path = vendor/other
	url = https://github.com/other-org/other-repo.git
[submodule "different-host"]
	path = vendor/gitlab
	url = https://gitlab.com/my-org/gitlab-repo.git
[submodule "self-reference"]
	path = self
	url = ../my-parent.git
`

	repos := ParseSameOwnerSubmoduleRepos(gitmodules, "my-org", "my-parent", "github.com")

	expected := []string{"shared-lib", "nested-repo", "ssh-repo"}
	if len(repos) != len(expected) {
		t.Fatalf("expected %d repos, got %d: %v", len(expected), len(repos), repos)
	}

	for i, name := range expected {
		if repos[i] != name {
			t.Errorf("expected repo %d to be %q, got %q", i, name, repos[i])
		}
	}
}

func TestParseSameOwnerSubmoduleReposCaseInsensitive(t *testing.T) {
	gitmodules := `
[submodule "lib"]
	url = https://github.com/My-Org/LIB-Core.git
`
	repos := ParseSameOwnerSubmoduleRepos(gitmodules, "my-org", "parent", "github.com")
	if len(repos) != 1 || repos[0] != "lib-core" {
		t.Fatalf("expected ['lib-core'], got %v", repos)
	}
}
