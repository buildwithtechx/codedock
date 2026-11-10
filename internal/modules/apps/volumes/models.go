package volumes

import "time"

type ServiceVolume struct {
	ID            string    `json:"id" db:"id"`
	ServiceID     string    `json:"serviceId" db:"service_id"`
	HostPath      string    `json:"hostPath" db:"host_path"`
	ContainerPath string    `json:"containerPath" db:"container_path"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
}
