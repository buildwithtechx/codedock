package http

import (
	"codedock/internal/models"
	"codedock/internal/services/backups"
	"codedock/internal/services/projects"
	"context"
	"fmt"
)

type backupUsers interface {
	GetUserByID(context.Context, string) (*models.User, error)
}

func configureScheduledBackupAuthorization(service *backups.BackupService, users backupUsers, projects *projects.ProjectService) {
	service.RunAuthorization = func(ctx context.Context, userID, project string) error {
		user, err := users.GetUserByID(ctx, userID)
		if err != nil || user == nil || !user.IsActive {
			return fmt.Errorf("scheduled backup owner is unavailable")
		}
		if project == "" {
			if user.Role != "owner" && user.Role != "admin" {
				return fmt.Errorf("control-plane backups require an instance administrator")
			}
			return nil
		}
		if !projects.HasPermission(ctx, project, userID, models.UserRole(user.Role), models.MemberPermissionAdmin) {
			return fmt.Errorf("scheduled backup owner no longer has project admin access")
		}
		return nil
	}
}
