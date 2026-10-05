package middleware

import (
	"codedock.run/codedock/internal/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInstanceOwnerCanAccessOrganizationWithoutMembership(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/organizations/test", nil)
	c := echo.New().NewContext(request, httptest.NewRecorder())
	c.Set("user", &models.UserClaims{UserID: "owner", Role: models.UserRoleOwner})
	guard := &AuthGuard{}
	if err := guard.verifyOrgRole(c, models.MemberPermissionOwner); err != nil {
		t.Fatal(err)
	}
}
