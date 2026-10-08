package systemdb

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ReadPassword(dataDir string) (string, error) {
	raw, err := os.ReadFile(filepath.Join(dataDir, "postgres", "password"))
	if err != nil {
		return "", fmt.Errorf("failed to read postgres password: %w", err)
	}
	password := strings.TrimSpace(string(raw))
	if len(password) < 32 {
		return "", fmt.Errorf("postgres password is missing or invalid")
	}
	return password, nil
}

func loadOrGeneratePassword(dataDir string) (string, error) {
	dir := filepath.Join(dataDir, "postgres")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", fmt.Errorf("failed to create postgres dir: %w", err)
	}
	path := filepath.Join(dir, "password")
	if raw, err := os.ReadFile(path); err == nil {
		password := strings.TrimSpace(string(raw))
		if len(password) >= 32 {
			return password, nil
		}
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("failed to read postgres password: %w", err)
	}
	password, err := generatePassword()
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(path, []byte(password+"\n"), 0o600); err != nil {
		return "", fmt.Errorf("failed to write postgres password: %w", err)
	}
	return password, nil
}

func generatePassword() (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("failed to generate postgres password: %w", err)
	}
	return hex.EncodeToString(raw), nil
}
