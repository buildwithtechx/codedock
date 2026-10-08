package system

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func exampleStubServer(contents, manifest string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/manifest.json" {
			w.Write([]byte(manifest))
			return
		}
		w.Write([]byte(contents))
	}))
}

func TestListExamplesMergesManifest(t *testing.T) {
	contents := `[{"name":"node-express","type":"dir","html_url":"https://github.com/org/examples/tree/main/node-express"},{"name":"not-yet-described","type":"dir","html_url":"https://github.com/org/examples/tree/main/not-yet-described"},{"name":"README.md","type":"file","html_url":"https://github.com/org/examples/blob/main/README.md"}]`
	manifest := `{"templates":[{"id":"node-express","name":"Node Express","description":"Minimal Node.js API built on Express","logo":"logos/node-express.svg"}]}`
	server := exampleStubServer(contents, manifest)
	defer server.Close()

	service := NewExampleServiceWithURLs(server.URL+"/contents", server.URL+"/manifest.json", server.URL+"/raw/")
	examples, err := service.ListExamples()
	if err != nil {
		t.Fatal(err)
	}
	if len(examples) != 2 {
		t.Fatalf("expected directories only, got %+v", examples)
	}
	described := examples[0]
	if described.ID != "node-express" || described.Name != "Node Express" {
		t.Errorf("expected manifest name, got %+v", described)
	}
	if described.Description != "Minimal Node.js API built on Express" {
		t.Errorf("expected manifest description, got %+v", described)
	}
	if described.Logo != server.URL+"/raw/logos/node-express.svg" {
		t.Errorf("expected raw logo URL, got %+v", described)
	}
	fallback := examples[1]
	if fallback.Name != "Not Yet Described" || fallback.Logo != "" {
		t.Errorf("expected formatted fallback without logo, got %+v", fallback)
	}
}

func TestListExamplesSurvivesBrokenManifest(t *testing.T) {
	contents := `[{"name":"go-fiber","type":"dir","html_url":"https://github.com/org/examples/tree/main/go-fiber"}]`
	server := exampleStubServer(contents, `not json`)
	defer server.Close()

	service := NewExampleServiceWithURLs(server.URL+"/contents", server.URL+"/manifest.json", server.URL+"/raw/")
	examples, err := service.ListExamples()
	if err != nil {
		t.Fatal(err)
	}
	if len(examples) != 1 || examples[0].Name != "Go Fiber" || examples[0].Logo != "" {
		t.Fatalf("expected listing without manifest enrichment, got %+v", examples)
	}
}

func TestListExamplesFailsWithoutContentsAndCache(t *testing.T) {
	server := exampleStubServer(`not json`, `{}`)
	defer server.Close()

	service := NewExampleServiceWithURLs(server.URL+"/contents", server.URL+"/manifest.json", server.URL+"/raw/")
	if _, err := service.ListExamples(); err == nil {
		t.Error("expected contents error without cache")
	}
}
