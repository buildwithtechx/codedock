package mail

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrMailNotFound = errors.New("not found")
	ErrMailExists   = errors.New("already exists")
)

func mapStatusError(status int, raw []byte) error {
	var payload struct {
		Error string `json:"error"`
	}
	message := strings.TrimSpace(string(raw))
	if err := json.Unmarshal(raw, &payload); err == nil && payload.Error != "" {
		message = payload.Error
	}
	switch status {
	case http.StatusNotFound:
		return ErrMailNotFound
	case http.StatusConflict:
		return ErrMailExists
	default:
		return fmt.Errorf("maild: status %d: %s", status, message)
	}
}

func mapRepoError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrMailNotFound
	}
	return err
}
