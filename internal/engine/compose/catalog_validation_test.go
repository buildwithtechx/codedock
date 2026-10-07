package compose

import (
	"strings"
	"testing"
)

func TestEmbeddedCatalogueValidates(t *testing.T) {
	manager, err := NewTemplateManager()
	if err != nil {
		t.Fatal(err)
	}
	ids := manager.ListTemplates()
	if len(ids) == 0 {
		t.Fatal("expected embedded catalogue templates")
	}
	for _, id := range ids {
		tmpl, err := manager.GetTemplate(id)
		if err != nil {
			t.Fatalf("load template %s: %v", id, err)
		}
		if err := ValidateTemplate(id, tmpl); err != nil {
			t.Errorf("template %s: %v", id, err)
		}
	}
}

func validCatalogueTemplate() ComposeTemplate {
	return ComposeTemplate{
		Services: map[string]ComposeService{
			"web": {
				Image: "example/web:1.0",
				Ports: []string{"8080:80"},
				Environment: []string{
					"APP_ENV=production",
					"APP_SECRET=${secret:APP_SECRET}",
				},
				Volumes: []string{"web-data:/data"},
			},
		},
		Volumes: map[string]any{"web-data": nil},
		XCodedock: &CodedockMetadata{
			IsOneClick:  true,
			Name:        "Example",
			Description: "Example application",
			Secrets:     []CodedockSecretSpec{{Var: "APP_SECRET", Length: 32}},
		},
	}
}

func TestValidateTemplateRejectsBrokenDocuments(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*ComposeTemplate)
		want   string
	}{
		{"no services", func(tmpl *ComposeTemplate) { tmpl.Services = nil }, "services are required"},
		{"no image", func(tmpl *ComposeTemplate) {
			service := tmpl.Services["web"]
			service.Image = ""
			tmpl.Services["web"] = service
		}, "needs an image"},
		{"bad port", func(tmpl *ComposeTemplate) {
			service := tmpl.Services["web"]
			service.Ports = []string{"99999:80"}
			tmpl.Services["web"] = service
		}, "invalid port mapping"},
		{"bind mount", func(tmpl *ComposeTemplate) {
			service := tmpl.Services["web"]
			service.Volumes = []string{"/host/path:/data"}
			tmpl.Services["web"] = service
		}, "declared named volume"},
		{"undeclared volume", func(tmpl *ComposeTemplate) {
			service := tmpl.Services["web"]
			service.Volumes = []string{"missing:/data"}
			tmpl.Services["web"] = service
		}, "not declared"},
		{"bad env", func(tmpl *ComposeTemplate) {
			service := tmpl.Services["web"]
			service.Environment = []string{"9BAD=value"}
			tmpl.Services["web"] = service
		}, "invalid environment entry"},
		{"no metadata", func(tmpl *ComposeTemplate) { tmpl.XCodedock = nil }, "x-codedock metadata is required"},
		{"no name", func(tmpl *ComposeTemplate) { tmpl.XCodedock.Name = "" }, "one-click name is required"},
		{"undeclared secret", func(tmpl *ComposeTemplate) { tmpl.XCodedock.Secrets = nil }, "undeclared secret"},
		{"unused secret", func(tmpl *ComposeTemplate) {
			tmpl.XCodedock.Secrets = append(tmpl.XCodedock.Secrets, CodedockSecretSpec{Var: "UNUSED"})
		}, "never referenced"},
		{"unknown dependency", func(tmpl *ComposeTemplate) {
			service := tmpl.Services["web"]
			service.DependsOn = []string{"missing"}
			tmpl.Services["web"] = service
		}, "unknown dependency"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tmpl := validCatalogueTemplate()
			tc.mutate(&tmpl)
			err := ValidateTemplate("example", tmpl)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("expected error containing %q, got %v", tc.want, err)
			}
		})
	}
}

func TestParseCatalogPort(t *testing.T) {
	cases := []struct {
		entry           string
		host, container int
		protocol        string
		expectError     bool
	}{
		{"80", 0, 80, "tcp", false},
		{"8080:80", 8080, 80, "tcp", false},
		{"8080:80/udp", 8080, 80, "udp", false},
		{"0:80", 0, 0, "", true},
		{"8080:0", 0, 0, "", true},
		{"8080:80/gre", 0, 0, "", true},
		{"1:2:3", 0, 0, "", true},
		{"http", 0, 0, "", true},
	}
	for _, tc := range cases {
		mapping, err := ParseCatalogPort(tc.entry)
		if tc.expectError {
			if err == nil {
				t.Errorf("entry %q: expected error", tc.entry)
			}
			continue
		}
		if err != nil {
			t.Errorf("entry %q: %v", tc.entry, err)
			continue
		}
		if mapping.Host != tc.host || mapping.Container != tc.container || mapping.Protocol != tc.protocol {
			t.Errorf("entry %q: got %+v", tc.entry, mapping)
		}
	}
}
