package models

import "time"

type RuntimeImageStatus string

const (
	RuntimeImageStatusCurrent         RuntimeImageStatus = "current"
	RuntimeImageStatusUpdateAvailable RuntimeImageStatus = "update_available"
	RuntimeImageStatusPinned          RuntimeImageStatus = "pinned"
	RuntimeImageStatusUnknown         RuntimeImageStatus = "unknown"
)

type RuntimeImageInfo struct {
	ServiceID     *string            `json:"serviceId"`
	ServiceName   *string            `json:"serviceName"`
	Role          string             `json:"role"`
	Ref           string             `json:"ref"`
	CurrentDigest *string            `json:"currentDigest"`
	LatestDigest  *string            `json:"latestDigest"`
	Status        RuntimeImageStatus `json:"status"`
}

type RuntimeImageCheckResult struct {
	CheckedAt time.Time          `json:"checkedAt"`
	Images    []RuntimeImageInfo `json:"images"`
}

type UpdateRuntimeRequest struct {
	ServiceIDs []string `json:"serviceIds,omitempty"`
}
