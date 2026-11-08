package organizations

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

type UpdateOrganizationMemberRequest struct {
	Permission MemberPermission `json:"permission"`
}
