package models

type RuntimeObservation struct {
	Logs    string         `json:"logs"`
	Metrics map[string]any `json:"metrics"`
}
