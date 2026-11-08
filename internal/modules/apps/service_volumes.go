package apps

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

func (s *Service) CreateVolume(ctx context.Context, vol *ServiceVolume) (*ServiceVolume, error) {
	if vol.ServiceID == "" || vol.HostPath == "" || vol.ContainerPath == "" {
		return nil, errors.New("serviceId, hostPath and containerPath are required")
	}
	if vol.ID == "" {
		vol.ID = uuid.New().String()
	}
	vol.CreatedAt = time.Now().UTC()
	if err := s.volumeRepo.Create(ctx, vol); err != nil {
		return nil, err
	}
	return vol, nil
}

func (s *Service) ListVolumes(ctx context.Context, serviceID string) ([]ServiceVolume, error) {
	if serviceID == "" {
		return nil, errors.New("serviceId required")
	}
	return s.volumeRepo.ListByService(ctx, serviceID)
}

func (s *Service) DeleteVolume(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id required")
	}
	return s.volumeRepo.Delete(ctx, id)
}

func (s *Service) EmptyVolume(ctx context.Context, serviceID, volumeID string) (*ServiceVolume, string, error) {
	if serviceID == "" || volumeID == "" {
		return nil, "", errors.New("serviceId and volumeId are required")
	}
	vol, err := s.volumeRepo.GetByID(ctx, volumeID)
	if err != nil {
		return nil, "", err
	}
	if vol == nil || vol.ServiceID != serviceID {
		return nil, "", errors.New("volume not found for service")
	}
	previousHostPath := vol.HostPath
	vol.HostPath = "codedock-fresh-" + uuid.New().String()
	if err := s.volumeRepo.Update(ctx, vol); err != nil {
		return nil, "", err
	}
	return vol, previousHostPath, nil
}
