package models

type RollbackCapacity struct {
	Window             int     `json:"window"`
	Source             string  `json:"source"`
	Explicit           *int    `json:"explicit"`
	SnapshotSizeBytes  *int64  `json:"snapshotSizeBytes"`
	MeasuredAt         *string `json:"measuredAt"`
	DiskFreeBytes      *int64  `json:"diskFreeBytes"`
	DiskTotalBytes     *int64  `json:"diskTotalBytes"`
	MaxWindow          int     `json:"maxWindow"`
	DiskBudgetFraction float64 `json:"diskBudgetFraction"`
	Strategy           string  `json:"strategy"`
}

type UpdateRollbackSettingsRequest struct {
	DefaultRollbackStrategy *string `json:"defaultRollbackStrategy,omitempty"`
	RollbackWindow          *int    `json:"rollbackWindow"`
}

type RetentionReclaimSummary struct {
	Purged int `json:"purged"`
	Failed int `json:"failed"`
}

type PinDeploymentResponse struct {
	ID     string `json:"id"`
	Pinned bool   `json:"pinned"`
}
