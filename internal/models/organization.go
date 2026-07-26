package models

import (
	"codedock.run/codedock/pkg/types"
)

type MemberPermission = types.MemberPermission
type MemberStatus = types.MemberStatus

const (
	MemberPermissionAdmin  = types.MemberPermissionAdmin
	MemberPermissionMember = types.MemberPermissionMember
	MemberStatusPending    = types.MemberStatusPending
	MemberStatusActive     = types.MemberStatusActive
)

type Organization = types.Organization
type OrganizationMember = types.OrganizationMember
type CreateOrganizationRequest = types.CreateOrganizationRequest
type UpdateOrganizationRequest = types.UpdateOrganizationRequest
