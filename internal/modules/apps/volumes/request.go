package volumes

// CreateVolumeRequest defines parameters for mounting a persistent volume.
type CreateVolumeRequest struct {
	HostPath      string `json:"hostPath"`
	ContainerPath string `json:"containerPath"`
}
