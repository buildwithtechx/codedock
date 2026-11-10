package system

// LocalCLISessionRequest defines optional parameters for requesting a local CLI session token.
type LocalCLISessionRequest struct {
	ClientVersion string `json:"clientVersion,omitempty"`
}

// SystemMaintenanceRequest defines parameters for triggering cleanup routines.
type SystemMaintenanceRequest struct {
	PruneImages   bool `json:"pruneImages"`
	PruneVolumes  bool `json:"pruneVolumes"`
	PruneNetworks bool `json:"pruneNetworks"`
}
