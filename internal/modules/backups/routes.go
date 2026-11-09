package backups

import (
	"github.com/labstack/echo/v4"
)

type RouteMiddleware struct {
	Admin          echo.MiddlewareFunc
	BackupWrite    echo.MiddlewareFunc
	DatabaseManage echo.MiddlewareFunc
}

func RegisterRoutes(authGroup *echo.Group, h *Handler, mw RouteMiddleware) {
	if h == nil {
		h = &Handler{}
	}
	if authGroup == nil {
		return
	}

	withMw := func(m echo.MiddlewareFunc) []echo.MiddlewareFunc {
		if m != nil {
			return []echo.MiddlewareFunc{m}
		}
		return nil
	}

	// Backup policies
	authGroup.GET("/backups", h.List)
	authGroup.POST("/backups", h.Create)
	authGroup.GET("/backups/:id", h.Get)
	authGroup.PUT("/backups/:id", h.Update, withMw(mw.BackupWrite)...)
	authGroup.DELETE("/backups/:id", h.Delete)
	authGroup.POST("/backups/:id/trigger", h.Trigger)
	authGroup.POST("/backups/:id/restore", h.Restore)
	authGroup.POST("/backups/:id/runs", h.StartRun, withMw(mw.BackupWrite)...)

	// Backup records
	authGroup.GET("/backups/:id/records", h.ListRecords)
	authGroup.GET("/backups/:id/records/:recordId/download", h.DownloadRecord)
	authGroup.DELETE("/backups/:id/records/:recordId", h.DeleteRecord)
	authGroup.GET("/backup-records", h.ListAllRecords)
	authGroup.PUT("/backup-records/:recordId/protection", h.ProtectRecord, withMw(mw.BackupWrite)...)
	authGroup.POST("/backup-records/:recordId/restore/review", h.PrepareRestore, withMw(mw.BackupWrite)...)
	authGroup.POST("/backup-operations/:operationId/apply", h.ApplyRestore, withMw(mw.BackupWrite)...)

	// Volume restore
	authGroup.GET("/backup-records/:recordId/volume-restore", h.VolumeRestoreTarget)
	authGroup.DELETE("/backup-records/:recordId/volume-restore", h.CancelVolumeRestore)

	// S3 destinations
	authGroup.GET("/s3-destinations", h.ListS3Destinations, withMw(mw.Admin)...)
	authGroup.POST("/s3-destinations", h.CreateS3Destination, withMw(mw.Admin)...)
	authGroup.POST("/s3-destinations/verify", h.VerifyS3Draft, withMw(mw.Admin)...)
	authGroup.PUT("/s3-destinations/:id", h.UpdateS3Destination, withMw(mw.Admin)...)
	authGroup.POST("/s3-destinations/:id/verify", h.VerifyS3Destination, withMw(mw.Admin)...)
	authGroup.POST("/s3-destinations/:id/default", h.SetDefaultS3Destination, withMw(mw.Admin)...)
	authGroup.DELETE("/s3-destinations/:id", h.DeleteS3Destination, withMw(mw.Admin)...)

	// SFTP destinations
	authGroup.GET("/sftp-destinations", h.ListSFTPDestinations, withMw(mw.Admin)...)
	authGroup.POST("/sftp-destinations", h.CreateSFTPDestination, withMw(mw.Admin)...)
	authGroup.POST("/sftp-destinations/:id/verify", h.VerifySFTPDestination, withMw(mw.Admin)...)
	authGroup.DELETE("/sftp-destinations/:id", h.DeleteSFTPDestination, withMw(mw.Admin)...)

	// Policy batches
	authGroup.GET("/projects/:id/policy-batches", h.ListBatches, withMw(mw.Admin)...)
	authGroup.POST("/projects/:id/policy-batches", h.CreateBatch, withMw(mw.Admin)...)
	authGroup.POST("/policy-batches/:batchId/trigger", h.TriggerBatch, withMw(mw.Admin)...)

	// Database scoped backup records
	authGroup.GET("/databases/:id/backups", h.ListRecordsByDatabase, withMw(mw.DatabaseManage)...)
	authGroup.POST("/databases/:id/backups", h.TriggerDatabaseBackup, withMw(mw.DatabaseManage)...)
}
