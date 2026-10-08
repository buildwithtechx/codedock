package repositories

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Vault interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

func OpenDatabase(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open(DriverPostgres, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open Postgres database: %w", err)
	}
	db.SetMaxOpenConns(25)
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping Postgres database: %w", err)
	}
	return db, nil
}
