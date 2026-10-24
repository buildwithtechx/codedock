package systemdb

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func resolvePostgresPort(dataDir string, preferred int) (int, error) {
	if stored, err := loadStoredPort(dataDir); err == nil {
		return stored, nil
	}
	return refreshPostgresPort(dataDir, preferred)
}

func refreshPostgresPort(dataDir string, preferred int) (int, error) {
	for port := preferred; port < preferred+100; port++ {
		if port < 1 || port > 65535 {
			continue
		}
		ln, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(port))
		if err != nil {
			continue
		}
		ln.Close()
		if err := storePort(dataDir, port); err != nil {
			return 0, err
		}
		return port, nil
	}
	return 0, fmt.Errorf("no free postgres port near %d", preferred)
}

func loadStoredPort(dataDir string) (int, error) {
	raw, err := os.ReadFile(filepath.Join(dataDir, "postgres", "port"))
	if err != nil {
		return 0, err
	}
	port, err := strconv.Atoi(strings.TrimSpace(string(raw)))
	if err != nil || port < 1 || port > 65535 {
		return 0, fmt.Errorf("stored postgres port is invalid")
	}
	return port, nil
}

func storePort(dataDir string, port int) error {
	dir := filepath.Join(dataDir, "postgres")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("failed to create postgres dir: %w", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "port"), []byte(strconv.Itoa(port)+"\n"), 0o600); err != nil {
		return fmt.Errorf("failed to write postgres port: %w", err)
	}
	return nil
}

func isPortConflict(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "ports are not available") ||
		strings.Contains(msg, "port is already allocated") ||
		strings.Contains(msg, "address already in use")
}
