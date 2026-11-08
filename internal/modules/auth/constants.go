package auth

type DeviceSessionStatus string

const (
	DeviceStatusPending  DeviceSessionStatus = "pending"
	DeviceStatusApproved DeviceSessionStatus = "approved"
	DeviceStatusDenied   DeviceSessionStatus = "denied"
	DeviceStatusExpired  DeviceSessionStatus = "expired"
)
