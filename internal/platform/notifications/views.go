package notifications

import (
	"embed"
	"fmt"
	"html/template"
)

//go:embed emails/*.tmpl
var emailFS embed.FS

var HTMLTemplates *template.Template

func LoadTemplates() error {
	var err error
	HTMLTemplates, err = template.ParseFS(emailFS, "emails/*.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse email templates: %w", err)
	}
	return nil
}
