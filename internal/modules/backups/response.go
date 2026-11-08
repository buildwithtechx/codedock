package backups

func RedactBackupConfig(cfg *BackupConfig) *BackupConfig {
	if cfg == nil {
		return nil
	}
	cfg.DbPassword = "********"
	return cfg
}

func RedactS3Destination(dest *S3Destination) *S3Destination {
	if dest == nil {
		return nil
	}
	dest.SecretAccessKey = "********"
	return dest
}

func RedactSFTPDestination(dest *SFTPDestination) *SFTPDestination {
	if dest == nil {
		return nil
	}
	dest.Password = "********"
	dest.PrivateKey = "********"
	return dest
}

type VolumeRestoreTargetResponse struct {
	TargetServerID string `json:"targetServerId"`
	TargetVolume   string `json:"targetVolume"`
}
