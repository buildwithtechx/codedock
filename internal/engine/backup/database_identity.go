package backup

import (
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
)

func databaseContainerIdentity(database *models.Database) string {
	if database.ContainerID != "" {
		return database.ContainerID
	}
	if database.Name != "" {
		return utils.NormalizeContainerName("codedock-db-" + database.Name)
	}
	return utils.NormalizeContainerName(database.ID)
}
