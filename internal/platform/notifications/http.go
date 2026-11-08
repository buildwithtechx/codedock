package notifications

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func postJSON(url string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
