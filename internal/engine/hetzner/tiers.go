package hetzner

import (
	"codedock/internal/models"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"
)

func Tiers() []models.ManagedTier {
	return []models.ManagedTier{
		{Name: "cx22", Cores: 2, MemoryGB: 4, DiskGB: 40, MonthlyPrice: 4.51},
		{Name: "cx32", Cores: 2, MemoryGB: 8, DiskGB: 80, MonthlyPrice: 8.75},
		{Name: "cx42", Cores: 4, MemoryGB: 16, DiskGB: 160, MonthlyPrice: 17.50},
		{Name: "cx52", Cores: 8, MemoryGB: 32, DiskGB: 320, MonthlyPrice: 34.99},
	}
}

func TierByName(name string) (models.ManagedTier, error) {
	for _, tier := range Tiers() {
		if tier.Name == name {
			return tier, nil
		}
	}
	return models.ManagedTier{}, fmt.Errorf("select a reviewed server tier")
}

func Regions() []string { return []string{"fsn1", "nbg1", "hel1", "ash", "hil", "sin"} }

func Images() []string { return []string{"ubuntu-24.04", "debian-12"} }

func GenerateKeyPair() (private, public string, err error) {
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", err
	}
	pemBlock, err := ssh.MarshalPrivateKey(privateKey, "")
	if err != nil {
		return "", "", err
	}
	publicKey, err := ssh.NewPublicKey(privateKey.Public())
	if err != nil {
		return "", "", err
	}
	return string(pem.EncodeToMemory(pemBlock)), strings.TrimSpace(string(ssh.MarshalAuthorizedKey(publicKey))), nil
}
