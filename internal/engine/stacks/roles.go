package stacks

import "strings"

type RoleKind string

const (
	RoleKindWeb       RoleKind = "web"
	RoleKindWorker    RoleKind = "worker"
	RoleKindScheduler RoleKind = "scheduler"
)

type RoleGates struct {
	Deps  []string
	Files []string
}

type StackRole struct {
	Name            string
	Kind            RoleKind
	Command         string
	Replicas        int
	Singleton       bool
	StopSignal      string
	StopGracePeriod string
	When            *RoleGates
	Unless          *RoleGates
}

func matchesGates(role StackRole, deps, files map[string]bool) bool {
	if role.When != nil {
		for _, dep := range role.When.Deps {
			if !deps[strings.ToLower(dep)] {
				return false
			}
		}
		for _, file := range role.When.Files {
			if !files[file] {
				return false
			}
		}
	}
	if role.Unless != nil {
		for _, dep := range role.Unless.Deps {
			if deps[strings.ToLower(dep)] {
				return false
			}
		}
		for _, file := range role.Unless.Files {
			if files[file] {
				return false
			}
		}
	}
	return true
}

func ResolveStackRoles(configRoles *[]StackRole, defaultRoles []StackRole, deps, files []string) []StackRole {
	if configRoles != nil {
		out := make([]StackRole, len(*configRoles))
		copy(out, *configRoles)
		return out
	}

	depSet := make(map[string]bool, len(deps))
	for _, dep := range deps {
		depSet[strings.ToLower(dep)] = true
	}
	fileSet := make(map[string]bool, len(files))
	for _, file := range files {
		fileSet[file] = true
	}

	matched := make([]StackRole, 0, len(defaultRoles))
	for _, role := range defaultRoles {
		if matchesGates(role, depSet, fileSet) {
			matched = append(matched, role)
		}
	}

	counts := make(map[RoleKind]int, len(matched))
	for _, role := range matched {
		counts[role.Kind]++
	}
	out := make([]StackRole, 0, len(matched))
	for _, role := range matched {
		if counts[role.Kind] == 1 {
			out = append(out, role)
		}
	}
	return out
}

func ResolveWebCommand(command string, optimize bool, optimizeCommand string) string {
	if !optimize || optimizeCommand == "" {
		return command
	}
	return optimizeCommand + " && " + command
}
