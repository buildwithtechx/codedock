package repositories

import (
	"database/sql"

	"codedock/internal/platform/database"
)

type Vault interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

func OpenDatabase(databaseURL string) (*sql.DB, error) {
	return database.OpenDatabase(databaseURL)
}
