package models

import (
	"time"
)

type MemberPermission string
type MemberStatus string

const (
	MemberPermissionOwner      MemberPermission = "owner"
	MemberPermissionAdmin      MemberPermission = "admin"
	MemberPermissionMember     MemberPermission = "member"
	MemberPermissionRestricted MemberPermission = "restricted"
	MemberStatusPending        MemberStatus     = "pending"
	MemberStatusActive         MemberStatus     = "active"
	MemberStatusAccepted       MemberStatus     = "accepted"
)

type Organization struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type OrganizationMember struct {
	ID             string           `json:"id" db:"id"`
	OrganizationID string           `json:"organizationId" db:"organization_id"`
	UserID         string           `json:"userId,omitempty" db:"user_id"`
	Email          string           `json:"email" db:"email"`
	Permission     MemberPermission `json:"permission" db:"permission"`
	Status         MemberStatus     `json:"status" db:"status"`
	InvitedAt      time.Time        `json:"invitedAt" db:"invited_at"`
	AcceptedAt     time.Time        `json:"acceptedAt" db:"accepted_at"`
	InviteLink     string           `json:"inviteLink,omitempty" db:"-"`
}

type CreateOrganizationRequest struct {
	Name string `json:"name"`
}

type UpdateOrganizationRequest struct {
	Name string `json:"name"`
}

type InviteOrganizationMemberRequest struct {
	Email      string           `json:"email"`
	Permission MemberPermission `json:"permission"`
	Delivery   string           `json:"delivery,omitempty"`
	Projects   []string         `json:"projects,omitempty"`
}
