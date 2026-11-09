package storage

import (
	"codedock/internal/utils"
)

func databaseContainerIdentity(database *Database) string {
	if database.ContainerID != "" {
		return database.ContainerID
	}
	if database.Name != "" {
		return utils.NormalizeContainerName("codedock-db-" + database.Name)
	}
	return utils.NormalizeContainerName(database.ID)
}
