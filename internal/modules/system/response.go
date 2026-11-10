package system

// SystemStatsResponse wraps host and docker statistics.
type SystemStatsResponse struct {
	Stats *SystemStats `json:"stats"`
}

// LocalCLISessionResponse returns authentication credentials for local CLI daemon handshakes.
type LocalCLISessionResponse struct {
	Token string `json:"token"`
	User  any    `json:"user"`
}

// UpdateStatusResponse wraps update daemon status.
type UpdateStatusResponse struct {
	Status *UpdateStatus `json:"status"`
}
