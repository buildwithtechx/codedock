package models

type ExampleApp struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Repo        string `json:"repo"`
	Icon        string `json:"icon,omitempty"`
	Logo        string `json:"logo,omitempty"`
}
