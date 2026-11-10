package app

import (
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"codedock/internal/config"
	"codedock/internal/http/apischema"
)

const openAPIGoldenPath = "../../docs/api/openapi.json"

func schemaTestRouter(t *testing.T) *echo.Echo {
	t.Helper()
	cfg := config.Get()
	previous := *cfg
	t.Cleanup(func() { *cfg = previous })
	cfg.Security.JWTSecret = "schema-test-secret"
	cfg.Cloud.Enabled = true

	services, err := initServices(Config{})
	if err != nil {
		t.Fatalf("failed to init services: %v", err)
	}

	app := &App{
		router:         echo.New(),
		authGuard:      services.authGuard,
		projectService: services.projectService,
		appService:     services.appService,
		mcpBridge:      services.bridge,
	}

	app.wireRoutes(services)
	return app.router
}

func TestEveryRouteHasSchema(t *testing.T) {
	router := schemaTestRouter(t)
	missing := []string{}
	for _, route := range router.Routes() {
		if route.Path == "/*" || route.Method == "echo_route_not_found" ||
			strings.HasPrefix(route.Path, "/.well-known/") ||
			strings.HasPrefix(route.Path, "/health") ||
			strings.HasPrefix(route.Path, "/api/health") ||
			strings.HasPrefix(route.Path, "/_codedock/") ||
			strings.HasPrefix(route.Path, "/api/v1/edge/") ||
			strings.HasPrefix(route.Path, "/api/mcp") ||
			strings.HasPrefix(route.Path, "/api/v1/mcp") ||
			strings.HasPrefix(route.Path, "/api/ai") ||
			strings.HasPrefix(route.Path, "/api/v1/ai") {
			continue
		}
		if _, ok := apischema.Find(route.Method, route.Path); !ok {
			missing = append(missing, route.Method+" "+route.Path)
		}
	}
	sort.Strings(missing)
	if len(missing) > 0 {
		t.Fatalf("routes without schema entries:\n%s", joinLines(missing))
	}
}

func TestSchemaHasNoOrphansOrDuplicates(t *testing.T) {
	router := schemaTestRouter(t)
	routes := map[string]bool{}
	for _, route := range router.Routes() {
		routes[route.Method+" "+route.Path] = true
	}
	seen := map[string]bool{}
	orphans := []string{}
	duplicates := []string{}
	for _, op := range apischema.All() {
		key := op.Method + " " + op.Path
		if seen[key] {
			duplicates = append(duplicates, key)
		}
		seen[key] = true
		if !routes[key] {
			orphans = append(orphans, key)
		}
	}
	if len(duplicates) > 0 {
		t.Errorf("duplicate schema entries:\n%s", joinLines(duplicates))
	}
	if len(orphans) > 0 {
		t.Errorf("schema entries without routes:\n%s", joinLines(orphans))
	}
}

func TestOpenAPIGolden(t *testing.T) {
	rendered, err := apischema.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if os.Getenv("UPDATE_OPENAPI") == "1" {
		if err := os.WriteFile(openAPIGoldenPath, rendered, 0600); err != nil {
			t.Fatal(err)
		}
		return
	}
	expected, err := os.ReadFile(openAPIGoldenPath)
	if err != nil {
		t.Fatalf("read golden file: %v (run with UPDATE_OPENAPI=1 to create it)", err)
	}
	if string(expected) != string(rendered) {
		t.Fatal("openapi document differs from docs/api/openapi.json (run with UPDATE_OPENAPI=1 to refresh)")
	}
}

func joinLines(lines []string) string {
	out := ""
	for _, line := range lines {
		out += "  " + line + "\n"
	}
	return out
}
