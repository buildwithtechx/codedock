package utils

import "codedock/internal/platform/vault"

type Vault = vault.Vault

func NewVault(dataDir string) (*Vault, error) {
	return vault.NewVault(dataDir)
}
