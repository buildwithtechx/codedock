package http

import (
	"os"
	"sort"
	"testing"
	"time"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/config"
	"codedock.run/codedock/internal/http/apischema"
	"codedock.run/codedock/internal/http/middleware"
	authservices "codedock.run/codedock/internal/services/auth"
)

const openAPIGoldenPath = "../../docs/api/openapi.json"

func schemaTestRouter(t *testing.T) *echo.Echo {
	t.Helper()
	cfg := config.Get()
	previous := *cfg
	t.Cleanup(func() { *cfg = previous })
	cfg.Security.JWTSecret = "schema-test-secret"
	tokens, err := authservices.NewTokenService()
	if err != nil {
		t.Fatal(err)
	}
	router := echo.New()
	server := &Server{
		router:          router,
		authGuard:       middleware.NewAuthGuard(tokens, nil, nil, nil, nil, nil),
		authRateLimiter: middleware.NewRateLimiter(10, time.Minute),
		otpRateLimiter:  middleware.NewRateLimiter(5, time.Minute),
		aiRateLimiter:   middleware.NewRateLimiter(5, time.Minute),
	}
	apiGroup := router.Group("/api")
	authGroup := apiGroup.Group("", server.authGuard.RequireAuth())
	server.registerAuthRoutes(apiGroup, authGroup)
	server.registerSystemRoutes(apiGroup, authGroup)
	server.registerUserRoutes(apiGroup, authGroup)
	server.registerProjectRoutes(apiGroup, authGroup)
	server.registerOrganizationRoutes(authGroup)
	server.registerServerRoutes(apiGroup, authGroup)
	server.registerManagedRoutes(authGroup)
	server.registerMigrationRoutes(authGroup)
	server.registerAnalyticsRoutes(authGroup)
	server.registerDatabaseRoutes(authGroup)
	server.registerAppRoutes(apiGroup, authGroup)
	server.registerDeploymentRoutes(authGroup)
	server.registerBackupRoutes(authGroup)
	server.registerClusterRoutes(authGroup)
	server.registerSettingsRoutes(apiGroup, authGroup)
	server.registerMiscRoutes(apiGroup, authGroup)
	server.registerBillingRoutes(apiGroup, authGroup)
	return router
}

func TestEveryRouteHasSchema(t *testing.T) {
	router := schemaTestRouter(t)
	missing := []string{}
	for _, route := range router.Routes() {
		if route.Path == "/*" {
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
