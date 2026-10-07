package models

type BackupFinalization struct {
	Record         *BackupRecord `json:"record"`
	FilePath       string        `json:"-"`
	S3URL          string        `json:"s3Url"`
	SFTPURL        string        `json:"sftpUrl"`
	ParentRecordID string        `json:"parentRecordId"`
	ExecLogs       string        `json:"execLogs"`
	SizeBytes      int64         `json:"sizeBytes"`
}
