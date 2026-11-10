package compose

// StackResponse wraps a single compose stack payload.
type StackResponse struct {
	Stack *ComposeStack `json:"stack"`
}

// StackListResponse wraps a list of compose stacks.
type StackListResponse struct {
	Stacks []ComposeStack `json:"stacks"`
	Total  int            `json:"total"`
}

// StackConfigResponse wraps the decrypted raw configuration of a compose stack.
type StackConfigResponse struct {
	Config string `json:"config"`
}
