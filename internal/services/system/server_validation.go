package system

import (
	"fmt"
	"golang.org/x/crypto/ssh"
)

func validateServerConnection(local bool, port int, transport, jumpHost, key, password string) error {
	if port < 0 || port > 65535 {
		return fmt.Errorf("SSH port must be between 1 and 65535")
	}
	if local {
		return nil
	}
	if transport != "" && transport != "direct" {
		return fmt.Errorf("only direct SSH transport is currently supported")
	}
	if jumpHost != "" {
		return fmt.Errorf("SSH jump hosts are currently unavailable")
	}
	if key == "" && password == "" {
		return fmt.Errorf("an SSH private key or password is required")
	}
	if key != "" {
		if _, err := ssh.ParsePrivateKey([]byte(key)); err != nil {
			return fmt.Errorf("invalid SSH private key: %w", err)
		}
	}
	return nil
}
