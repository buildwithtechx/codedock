package http

import (
	"codedock/internal/modules/dns"
	"codedock/internal/modules/monitoring"
)

type monitoringHandlersGroup struct {
	operationHandler     *monitoring.OperationHandler
	scheduledTaskHandler *monitoring.ScheduledTaskHandler
	notifSettingsHandler *monitoring.NotificationSettingsHandler
	updaterHandler       *monitoring.UpdaterHandler
	notificationHandler  *monitoring.NotificationHandler
	systemHandler        *monitoring.SystemHandler
	metricsHandler       *monitoring.MetricsHandler
	logHandler           *monitoring.LogHandler
	exampleHandler       *monitoring.ExampleHandler
	serviceLogsWSHandler *monitoring.ServiceLogsWSHandler
	ingressHandler       *dns.IngressHandler
}

func createMonitoringHandlers(
	d *serverSetupDeps,
	updaterService *monitoring.UpdaterService,
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
		repo:  dns.NewIngressRepository(d.db),
		cfSvc: dns.NewCloudflareIngressService(nil),
	}

	return monitoringHandlersGroup{
		operationHandler: monitoring.NewOperationHandler(&operationsServiceAdapter{svc: d.operationService}),
		scheduledTaskHandler: monitoring.NewScheduledTaskHandler(
			&scheduledTaskServiceAdapter{svc: d.scheduledTaskService},
			&scheduledTaskAccessAdapter{appService: d.appService, projectService: d.projectService},
		),
		notifSettingsHandler: monitoring.NewNotificationSettingsHandler(&notifSettingsServiceAdapter{
			svc:     d.notifSettingsService,
			orgRepo: d.orgRepo,
		}),
		updaterHandler:      monitoring.NewUpdaterHandler(&updaterServiceAdapter{svc: updaterService}),
		notificationHandler: monitoring.NewNotificationHandler(&notificationServiceAdapter{}),
		systemHandler: monitoring.NewSystemHandler(
			&systemServiceAdapter{svc: monitoring.NewSystemService()},
			&systemAdminFinderAdapter{userRepo: d.userRepo},
			&systemTokenGenAdapter{tokenService: d.tokenService},
		),
		metricsHandler:       monitoring.NewMetricsHandler(&metricsServiceAdapter{svc: metricsService}),
		logHandler:           monitoring.NewLogHandler(&logServiceAdapter{svc: logService}),
		exampleHandler:       monitoring.NewExampleHandler(&exampleServiceAdapter{svc: monitoring.NewExampleService()}),
		serviceLogsWSHandler: wsHandler,
		ingressHandler:       dns.NewIngressHandler(ingressAdapter),
	}
}
