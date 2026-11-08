package organizations

import "time"

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
