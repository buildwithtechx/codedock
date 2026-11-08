package backups

import (
	"context"
	"fmt"

	"codedock/internal/modules/backups/providers"
)

func s3ToProviderConfig(dest *S3Destination) providers.S3DestinationConfig {
	return providers.S3DestinationConfig{
		Endpoint:        dest.Endpoint,
		Bucket:          dest.Bucket,
		Region:          dest.Region,
		AccessKeyID:     dest.AccessKeyID,
		SecretAccessKey: dest.SecretAccessKey,
		PathPrefix:      dest.PathPrefix,
	}
}

func sftpToProviderConfig(dest *SFTPDestination) providers.SFTPDestinationConfig {
	return providers.SFTPDestinationConfig{
		Host:       dest.Host,
		Port:       dest.Port,
		Username:   dest.Username,
		Password:   dest.Password,
		PrivateKey: dest.PrivateKey,
		PathPrefix: dest.PathPrefix,
	}
}

func (s *Service) CreateS3Destination(ctx context.Context, dest *S3Destination) error {
	if err := providers.EnsureS3Bucket(ctx, s3ToProviderConfig(dest)); err != nil {
		return err
	}
	return s.s3Repo.CreateS3Destination(ctx, dest)
}

func (s *Service) UpdateS3Destination(ctx context.Context, dest *S3Destination) error {
	if dest.SecretAccessKey != "" && dest.SecretAccessKey != "********" {
		if err := providers.EnsureS3Bucket(ctx, s3ToProviderConfig(dest)); err != nil {
			return err
		}
	}
	return s.s3Repo.UpdateS3Destination(ctx, dest)
}

func (s *Service) ListS3Destinations(ctx context.Context) ([]*S3Destination, error) {
	return s.s3Repo.ListS3Destinations(ctx)
}

func (s *Service) GetS3Destination(ctx context.Context, id string) (*S3Destination, error) {
	return s.s3Repo.GetS3Destination(ctx, id)
}

func (s *Service) DeleteS3Destination(ctx context.Context, id string) error {
	return s.s3Repo.DeleteS3Destination(ctx, id)
}

func (s *Service) SetDefaultS3Destination(ctx context.Context, id string) error {
	return s.s3Repo.SetDefaultDestination(ctx, id)
}

func (s *Service) VerifyS3Destination(ctx context.Context, id string) (bool, error) {
	dest, err := s.s3Repo.GetS3Destination(ctx, id)
	if err != nil {
		return false, err
	}
	verifyErr := providers.CheckS3Bucket(ctx, s3ToProviderConfig(dest))
	var errMsg string
	if verifyErr != nil {
		errMsg = verifyErr.Error()
	}
	_ = s.s3Repo.RecordVerificationResult(ctx, id, verifyErr == nil, errMsg)
	return verifyErr == nil, verifyErr
}

func (s *Service) VerifyS3Draft(ctx context.Context, dest *S3Destination) (bool, error) {
	err := providers.CheckS3Bucket(ctx, s3ToProviderConfig(dest))
	return err == nil, err
}

func (s *Service) CreateSFTPDestination(ctx context.Context, dest *SFTPDestination) error {
	if s.sftp == nil {
		return ErrSFTPStorageUnavailable
	}
	if dest.Name == "" || dest.Host == "" || dest.ProjectID == "" {
		return fmt.Errorf("name, host and project are required")
	}
	if dest.Port <= 0 {
		dest.Port = 22
	}
	return s.sftp.Create(ctx, dest)
}

func (s *Service) ListSFTPDestinations(ctx context.Context, projectID string) ([]SFTPDestination, error) {
	if s.sftp == nil {
		return nil, ErrSFTPStorageUnavailable
	}
	return s.sftp.ListByProject(ctx, projectID)
}

func (s *Service) GetSFTPDestination(ctx context.Context, id string) (*SFTPDestination, error) {
	if s.sftp == nil {
		return nil, ErrSFTPStorageUnavailable
	}
	return s.sftp.Get(ctx, id)
}

func (s *Service) DeleteSFTPDestination(ctx context.Context, id string) error {
	if s.sftp == nil {
		return ErrSFTPStorageUnavailable
	}
	return s.sftp.Delete(ctx, id)
}

func (s *Service) VerifySFTPDestination(ctx context.Context, id string) (bool, error) {
	if s.sftp == nil {
		return false, ErrSFTPStorageUnavailable
	}
	dest, err := s.sftp.Get(ctx, id)
	if err != nil {
		return false, err
	}
	err = providers.VerifySFTP(ctx, sftpToProviderConfig(dest))
	return err == nil, err
}
