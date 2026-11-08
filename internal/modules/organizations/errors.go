package organizations

import "errors"

var (
	ErrUserIDRequired          = errors.New("user id is required")
	ErrNameRequired            = errors.New("name is required")
	ErrEmailRequired           = errors.New("email is required")
	ErrMemberIDRequired        = errors.New("member id is required")
	ErrInsufficientPermissions = errors.New("insufficient permissions to update organization")
	ErrAlreadyInvitedOrMember  = errors.New("user already invited or is a member")
	ErrTargetMemberNotFound    = errors.New("target member not found")
	ErrOwnerInviteOnly         = errors.New("only organization owners can invite new owners")
	ErrOwnerManageOnly         = errors.New("only organization owners can manage owner permissions")
	ErrCannotRemoveLastOwner   = errors.New("cannot remove the last owner of an organization")
	ErrCannotDemoteLastOwner   = errors.New("cannot demote the last owner of an organization")
)
