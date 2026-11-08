package clusters

type OperationReviewResponse struct {
	OperationID  string `json:"operationId"`
	Confirmation string `json:"confirmation"`
	Summary      string `json:"summary,omitempty"`
}

type ClusterReviewResponse struct {
	Review *OperationReviewResponse `json:"review"`
	Plan   *ClusterPlan             `json:"plan"`
}

type StorageReviewResponse struct {
	Review *OperationReviewResponse `json:"review"`
	Plan   *StorageSetupPlan        `json:"plan"`
}

type ClusterDataReviewResponse struct {
	Review *OperationReviewResponse `json:"review"`
	Plan   *ClusterDataPlan         `json:"plan"`
}
