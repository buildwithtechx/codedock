package build

import (
	"fmt"
	"path/filepath"
	"strings"
)

func ResolveSourceRoot(source, root string) (string, error) {
	base, err := filepath.Abs(source)
	if err != nil {
		return "", err
	}
	base, err = filepath.EvalSymlinks(base)
	if err != nil {
		return "", fmt.Errorf("repository checkout unavailable: %w", err)
	}
	root = strings.TrimPrefix(strings.ReplaceAll(root, "\\", "/"), "/")
	resolved, err := filepath.EvalSymlinks(filepath.Join(base, filepath.FromSlash(root)))
	if err != nil {
		return "", fmt.Errorf("repository root unavailable: %w", err)
	}
	relative, err := filepath.Rel(base, resolved)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("repository root must stay within the checkout")
	}
	return resolved, nil
}
