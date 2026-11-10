package templates

// TemplateResponse wraps a single template item.
type TemplateResponse struct {
	App *OneClickApp `json:"app"`
}

// TemplateListResponse wraps a list of one-click templates.
type TemplateListResponse struct {
	Apps  []OneClickApp `json:"apps"`
	Total int           `json:"total"`
}
