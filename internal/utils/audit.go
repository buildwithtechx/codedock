package utils

import (
	"reflect"

	"github.com/labstack/echo/v4"
)

func ExtractStructField(v any, name string) string {
	if v == nil {
		return ""
	}
	r := reflect.ValueOf(v)
	if r.Kind() == reflect.Ptr {
		if r.IsNil() {
			return ""
		}
		r = r.Elem()
	}
	if r.Kind() != reflect.Struct {
		return ""
	}
	field := r.FieldByName(name)
	if field.IsValid() && field.Kind() == reflect.String {
		return field.String()
	}
	return ""
}

func ExtractUserID(c echo.Context) string {
	return ExtractStructField(c.Get("user"), "UserID")
}

func ExtractClaims(c echo.Context) (userID, email string) {
	u := c.Get("user")
	return ExtractStructField(u, "UserID"), ExtractStructField(u, "Email")
}

func ExtractUser(c echo.Context) (userID string, role string, ok bool) {
	u := c.Get("user")
	if u == nil {
		return "", "", false
	}
	r := reflect.ValueOf(u)
	if r.Kind() == reflect.Ptr {
		if r.IsNil() {
			return "", "", false
		}
		r = r.Elem()
	}
	if r.Kind() != reflect.Struct {
		return "", "", false
	}
	uField := r.FieldByName("UserID")
	rField := r.FieldByName("Role")
	if uField.IsValid() && rField.IsValid() {
		return uField.String(), rField.String(), true
	}
	return "", "", false
}

func AuditActorID(c echo.Context) string {
	if id := ExtractUserID(c); id != "" {
		return id
	}
	return "system"
}
