package volumes

import "context"

// ServiceVolumeRepository defines database persistence for Docker volume mounts.
type ServiceVolumeRepository interface {
	Create(ctx context.Context, volume *ServiceVolume) error
	GetByID(ctx context.Context, id string) (*ServiceVolume, error)
	ListByService(ctx context.Context, serviceID string) ([]ServiceVolume, error)
	Update(ctx context.Context, volume *ServiceVolume) error
	Delete(ctx context.Context, id string) error
}

// VolumeService provides business operations for managing persistent volumes.
type VolumeService interface {
	CreateVolume(ctx context.Context, vol *ServiceVolume) (*ServiceVolume, error)
	ListVolumes(ctx context.Context, serviceID string) ([]ServiceVolume, error)
	DeleteVolume(ctx context.Context, id string) error
	EmptyVolume(ctx context.Context, serviceID, volumeID string) (*ServiceVolume, string, error)
}

// AccessChecker validates user permissions for volume operations.
type AccessChecker interface {
	HasPermission(ctx context.Context, projectID, userID, role, permission string) bool
	AccessibleServiceIDs(ctx context.Context, userID, role string) (map[string]bool, error)
}
