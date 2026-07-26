package models

import (
	"codedock.run/codedock/pkg/types"
)

type MemberPermission = types.MemberPermission
type MemberStatus = types.MemberStatus

const (
	MemberPermissionOwner  = types.MemberPermissionOwner
	MemberPermissionAdmin  = types.MemberPermissionAdmin
	MemberPermissionMember = types.MemberPermissionMember
	MemberStatusPending    = types.MemberStatusPending
	MemberStatusActive     = types.MemberStatusActive
	MemberStatusAccepted   = types.MemberStatusAccepted
)

type Organization = types.Organization
type OrganizationMember = types.OrganizationMember
type CreateOrganizationRequest = types.CreateOrganizationRequest
type UpdateOrganizationRequest = types.UpdateOrganizationRequest
