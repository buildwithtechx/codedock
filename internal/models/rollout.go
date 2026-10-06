package models

type RolloutContainer struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Running bool   `json:"running"`
}

type RolloutJournal struct {
	ID          string             `json:"id"`
	PreviousApp AppService         `json:"previousApp"`
	Previous    []RolloutContainer `json:"previous"`
	Committed   bool               `json:"committed"`
}
