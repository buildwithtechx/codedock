package databases

type DatabaseQueryResponse struct {
	Columns         []string         `json:"columns"`
	Rows            []map[string]any `json:"rows"`
	RowCount        int              `json:"rowCount"`
	ExecutionTimeMs int64            `json:"executionTimeMs"`
	Result          any              `json:"result,omitempty"`
}
