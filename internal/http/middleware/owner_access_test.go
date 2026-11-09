package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"codedock/internal/modules/organizations"
	"codedock/internal/modules/users"
)

func TestOwnerBypassesOrgMemberLookup(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/organizations/org-1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("orgId")
	c.SetParamValues("org-1")
	c.Set("user", &users.UserClaims{UserID: "owner", Role: users.UserRoleOwner})

	guard := &AuthGuard{}
	if err := guard.verifyOrgRole(c, organizations.MemberPermissionOwner); err != nil {
		t.Fatalf("expected owner role to bypass org membership check, got error: %v", err)
	}
}
