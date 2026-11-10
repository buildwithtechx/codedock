package vars

type CreateVariableRequest struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	IsSecret bool   `json:"isSecret"`
}

type UpdateVariableRequest struct {
	Key      *string `json:"key,omitempty"`
	Value    *string `json:"value,omitempty"`
	IsSecret *bool   `json:"isSecret,omitempty"`
}
