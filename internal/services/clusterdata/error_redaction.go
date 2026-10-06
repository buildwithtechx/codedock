package clusterdata

import (
	"codedock.run/codedock/internal/models"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

func redactDataError(plan *models.ClusterDataPlan, cause error) error {
	if cause == nil {
		return nil
	}
	values := []string{plan.Password}
	var document struct {
		Items []map[string]any `json:"items"`
	}
	if json.Unmarshal([]byte(plan.Manifest), &document) == nil {
		for _, object := range document.Items {
			if object["kind"] != "Secret" {
				continue
			}
			data, _ := object["stringData"].(map[string]any)
			for _, raw := range data {
				value, _ := raw.(string)
				values = append(values, value)
			}
		}
	}
	message := cause.Error()
	for _, value := range values {
		if value != "" {
			message = strings.ReplaceAll(message, value, "[redacted]")
			message = strings.ReplaceAll(message, base64.StdEncoding.EncodeToString([]byte(value)), "[redacted]")
		}
	}
	return errors.New(message)
}
