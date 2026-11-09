package monitoring

type AttentionActionRequest struct {
	Action string            `json:"action"`
	Params map[string]string `json:"params"`
}

type PathsToggleRequest struct {
	Enabled bool `json:"enabled"`
}
