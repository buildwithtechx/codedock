package operations

// OperationResponse wraps an operation details payload.
type OperationResponse struct {
	Operation *Operation `json:"operation"`
}

// OperationListResponse wraps a slice of operations.
type OperationListResponse struct {
	Operations []Operation `json:"operations"`
	Total      int         `json:"total"`
}
