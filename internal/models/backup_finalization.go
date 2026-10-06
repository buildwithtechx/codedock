package models

type BackupFinalization struct {
	Record    *BackupRecord `json:"record"`
	FilePath  string        `json:"-"`
	S3URL     string        `json:"s3Url"`
	ExecLogs  string        `json:"execLogs"`
	SizeBytes int64         `json:"sizeBytes"`
}
