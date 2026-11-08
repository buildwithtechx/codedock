package backups

import "errors"

var (
	ErrBackupConfigNotFound         = errors.New("backup config not found")
	ErrBackupRecordNotFound         = errors.New("backup record not found")
	ErrS3DestinationNotFound        = errors.New("s3 destination not found")
	ErrSFTPDestinationNotFound      = errors.New("sftp destination not found")
	ErrPolicyBatchNotFound          = errors.New("policy batch not found")
	ErrInvalidPayload               = errors.New("invalid payload")
	ErrMissingID                    = errors.New("missing id parameter")
	ErrBackupIDRequired             = errors.New("backup id is required")
	ErrUnauthorized                 = errors.New("unauthorized")
	ErrInsufficientPermissions      = errors.New("insufficient permissions")
	ErrInsufficientAdminPermissions = errors.New("insufficient admin permissions")
	ErrProducerImmutable            = errors.New("backup producer is immutable; create a separate policy")
	ErrValidBackupConfigRequired    = errors.New("valid backup config required")
	ErrProtectionUnavailable        = errors.New("backup protection storage unavailable")
	ErrRestoreRuntimeUnavailable    = errors.New("restore runtime unavailable")
	ErrBatchStorageUnavailable      = errors.New("policy batch storage unavailable")
	ErrSFTPStorageUnavailable       = errors.New("sftp storage unavailable")
	ErrBatchNameAndProjectRequired  = errors.New("batch name and project are required")
	ErrInvalidSchedule              = errors.New("invalid backup schedule")
	ErrProtectionExpiryTooFar       = errors.New("protection expiry must be within ten years")
	ErrOnlyCompletedCanBeProtected  = errors.New("only completed backup records can be protected")
)
