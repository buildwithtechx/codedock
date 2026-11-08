package http

import (
	"codedock/internal/config"
	"codedock/internal/modules/monitoring"
	"codedock/internal/repositories"
	systemservices "codedock/internal/services/system"
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
	billingHandler       *monitoring.BillingHandler
	serviceLogsWSHandler *monitoring.ServiceLogsWSHandler
	ingressHandler       *monitoring.IngressHandler
	actionsHandler       *monitoring.ActionsHandler
	supportHandler       *monitoring.SupportHandler
}

func createMonitoringHandlers(
	d *serverSetupDeps,
	updaterService *systemservices.UpdaterService,
	metricsService *systemservices.MetricsService,
	logService *systemservices.LogService,
) monitoringHandlersGroup {
	wsHandler := monitoring.NewServiceLogsWSHandler(
		d.tokenService,
		&wsAccessCheckerAdapter{appService: d.appService, projectService: d.projectService},
		&wsUserValidatorAdapter{userRepo: d.userRepo},
	)
	wsHandler.Streams = d.deploymentService

	var billingHandler *monitoring.BillingHandler
	if config.Get().Cloud.Enabled {
		billingHandler = monitoring.NewBillingHandler(&billingServiceAdapter{svc: systemservices.NewBillingService(d.userRepo)})
	}

	ingressAdapter := &ingressServiceAdapter{
		repo:  repositories.NewIngressRepository(d.db),
		cfSvc: systemservices.NewCloudflareIngressService(nil),
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
		updaterHandler: monitoring.NewUpdaterHandler(&updaterServiceAdapter{svc: updaterService}),
		notificationHandler: monitoring.NewNotificationHandler(&notificationServiceAdapter{
			svc: systemservices.NewNotificationService(d.dispatcherService),
		}),
		systemHandler: monitoring.NewSystemHandler(
			&systemServiceAdapter{svc: systemservices.NewSystemService()},
			&systemAdminFinderAdapter{userRepo: d.userRepo},
			&systemTokenGenAdapter{tokenService: d.tokenService},
		),
		metricsHandler:       monitoring.NewMetricsHandler(&metricsServiceAdapter{svc: metricsService}),
		logHandler:           monitoring.NewLogHandler(&logServiceAdapter{svc: logService}),
		exampleHandler:       monitoring.NewExampleHandler(&exampleServiceAdapter{svc: systemservices.NewExampleService()}),
		billingHandler:       billingHandler,
		serviceLogsWSHandler: wsHandler,
		ingressHandler:       monitoring.NewIngressHandler(ingressAdapter),
		actionsHandler:       monitoring.NewActionsHandler(),
		supportHandler:       monitoring.NewSupportHandler(),
	}
}
