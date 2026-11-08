package users

import (
	"reflect"

	"github.com/labstack/echo/v4"
)

func extractUserID(c echo.Context) string {
	if u := c.Get("user"); u != nil {
		v := reflect.ValueOf(u)
		if v.Kind() == reflect.Ptr && !v.IsNil() {
			v = v.Elem()
		}
		if v.Kind() == reflect.Struct {
			if f := v.FieldByName("UserID"); f.IsValid() && f.Kind() == reflect.String {
				return f.String()
			}
		}
	}
	return ""
}

func auditActorID(c echo.Context) string {
	if id := extractUserID(c); id != "" {
		return id
	}
	return "system"
}
