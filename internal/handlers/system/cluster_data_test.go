package system

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClusterDataReviewRequiresAuthenticatedOwner(t *testing.T) {
	context := echo.New().NewContext(httptest.NewRequest(http.MethodPost, "/review", nil), httptest.NewRecorder())
	if err := (&ClusterDataHandler{}).Review(context); err == nil {
		t.Fatal("anonymous database review accepted")
	} else if response, ok := err.(*echo.HTTPError); !ok || response.Code != http.StatusUnauthorized {
		t.Fatal(err)
	}
}
