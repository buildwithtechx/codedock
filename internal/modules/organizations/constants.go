package organizations

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
