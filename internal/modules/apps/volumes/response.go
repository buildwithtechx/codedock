package volumes

// VolumeResponse wraps a single service volume.
type VolumeResponse struct {
	Volume *ServiceVolume `json:"volume"`
}

// VolumeListResponse wraps a list of service volumes.
type VolumeListResponse struct {
	Volumes []ServiceVolume `json:"volumes"`
	Total   int             `json:"total"`
}
