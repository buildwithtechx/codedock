package backups

import (
	"context"
	"reflect"

	"github.com/labstack/echo/v4"
)

type UserContext struct {
	UserID string
	Role   string
}

func extractUser(c echo.Context) (userID string, role string, ok bool) {
	val := c.Get("user")
	if val == nil {
		return "", "", false
	}
	r := reflect.ValueOf(val)
	if r.Kind() == reflect.Ptr {
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

func auditActorID(c echo.Context) string {
	userID, _, ok := extractUser(c)
	if ok && userID != "" {
		return userID
	}
	return "system"
}

func extractStructField(v any, name string) string {
	if v == nil {
		return ""
	}
	r := reflect.ValueOf(v)
	if r.Kind() == reflect.Ptr {
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

type ResourceProjectResolver interface {
	GetDatabaseProjectID(ctx context.Context, databaseID string) (string, error)
	GetServiceProjectID(ctx context.Context, serviceID string) (string, error)
	HasPermission(ctx context.Context, projectID, userID, role, permission string) bool
}
