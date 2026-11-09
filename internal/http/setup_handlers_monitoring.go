package http

import (
	"codedock/internal/modules/dns/ingress"
	"codedock/internal/modules/monitoring"
	notifsmod "codedock/internal/modules/notifications"
	"codedock/internal/modules/operations"
	"codedock/internal/modules/system"
	"codedock/internal/modules/tasks"
)

type monitoringHandlersGroup struct {
	operationHandler     *operations.OperationHandler
	scheduledTaskHandler *tasks.ScheduledTaskHandler
	notifSettingsHandler *notifsmod.NotificationSettingsHandler
	updaterHandler       *system.UpdaterHandler
	notificationHandler  *notifsmod.NotificationHandler
	systemHandler        *system.SystemHandler
	metricsHandler       *monitoring.MetricsHandler
	logHandler           *monitoring.LogHandler
	exampleHandler       *monitoring.ExampleHandler
	serviceLogsWSHandler *monitoring.ServiceLogsWSHandler
	ingressHandler       *ingress.IngressHandler
}

func createMonitoringHandlers(
	d *serverSetupDeps,
	updaterService *system.UpdaterService,
	metricsService *monitoring.MetricsService,
	logService *monitoring.LogService,
) monitoringHandlersGroup {
	wsHandler := monitoring.NewServiceLogsWSHandler(
		d.tokenService,
		&wsAccessCheckerAdapter{appService: d.appService, projectService: d.projectService},
		&wsUserValidatorAdapter{userRepo: d.userRepo},
	)
	wsHandler.Streams = d.deploymentService

	ingressAdapter := &ingressServiceAdapter{
		repo:  ingress.NewIngressRepository(d.db),
		cfSvc: ingress.NewCloudflareIngressService(nil),
	}

	return monitoringHandlersGroup{
		operationHandler: operations.NewOperationHandler(&operationsServiceAdapter{svc: d.operationService}),
		scheduledTaskHandler: tasks.NewScheduledTaskHandler(
			&scheduledTaskServiceAdapter{svc: d.scheduledTaskService},
			&scheduledTaskAccessAdapter{appService: d.appService, projectService: d.projectService},
		),
		notifSettingsHandler: notifsmod.NewNotificationSettingsHandler(&notifSettingsServiceAdapter{
			svc:     d.notifSettingsService,
			orgRepo: d.orgRepo,
		}),
		updaterHandler:      system.NewUpdaterHandler(&updaterServiceAdapter{svc: updaterService}),
		notificationHandler: notifsmod.NewNotificationHandler(&notificationServiceAdapter{}),
		systemHandler: system.NewSystemHandler(
			&systemServiceAdapter{svc: system.NewSystemService()},
			&systemAdminFinderAdapter{userRepo: d.userRepo},
			&systemTokenGenAdapter{tokenService: d.tokenService},
		),
		metricsHandler:       monitoring.NewMetricsHandler(&metricsServiceAdapter{svc: metricsService}),
		logHandler:           monitoring.NewLogHandler(&logServiceAdapter{svc: logService}),
		exampleHandler:       monitoring.NewExampleHandler(&exampleServiceAdapter{svc: monitoring.NewExampleService()}),
		serviceLogsWSHandler: wsHandler,
		ingressHandler:       ingress.NewIngressHandler(ingressAdapter),
	}
}
