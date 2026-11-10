package vars

// VariableResponse wraps a single variable.
type VariableResponse struct {
	Variable *Variable `json:"variable"`
}

// VariableListResponse wraps a list of variables.
type VariableListResponse struct {
	Variables []*Variable `json:"variables"`
	Total     int         `json:"total"`
}
