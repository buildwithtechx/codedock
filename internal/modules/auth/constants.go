package auth

type UserRole string

const (
	UserRoleOwner  UserRole = "owner"
	UserRoleAdmin  UserRole = "admin"
	UserRoleMember UserRole = "member"
)

type MemberPermission string
type MemberStatus string

const (
	MemberPermissionOwner      MemberPermission = "owner"
	MemberPermissionAdmin      MemberPermission = "admin"
	MemberPermissionMember     MemberPermission = "member"
	MemberPermissionRestricted MemberPermission = "restricted"

	MemberStatusPending  MemberStatus = "pending"
	MemberStatusActive   MemberStatus = "active"
	MemberStatusAccepted MemberStatus = "accepted"
)

type DeviceSessionStatus string

const (
	DeviceStatusPending  DeviceSessionStatus = "pending"
	DeviceStatusApproved DeviceSessionStatus = "approved"
	DeviceStatusDenied   DeviceSessionStatus = "denied"
	DeviceStatusExpired  DeviceSessionStatus = "expired"
)
