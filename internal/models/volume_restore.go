package models

type VolumeRestoreTarget struct {
	RecordID        string `json:"recordId"`
	VolumeName      string `json:"volumeName"`
	ContainerID     string `json:"-"`
	VolumeCreatedAt string `json:"-"`
	TimeoutSeconds  int    `json:"timeoutSeconds"`
}

type VolumeRestoreRequest struct {
	VolumeName       string `json:"volumeName"`
	ConfirmOverwrite bool   `json:"confirmOverwrite"`
}
