package http

import (
	"codedock/internal/adapters"
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
		&adapters.WSAccessCheckerAdapter{AppService: d.appService, ProjectService: d.projectService},
		&adapters.WSUserValidatorAdapter{UserRepo: d.userRepo},
	)
	wsHandler.Streams = d.deploymentService

	ingressAdapter := &adapters.IngressServiceAdapter{
		Repo:  ingress.NewIngressRepository(d.db),
		CfSvc: ingress.NewCloudflareIngressService(nil),
	}

	return monitoringHandlersGroup{
		operationHandler: operations.NewOperationHandler(&adapters.OperationsServiceAdapter{Svc: d.operationService}),
		scheduledTaskHandler: tasks.NewScheduledTaskHandler(
			&adapters.ScheduledTaskServiceAdapter{Svc: d.scheduledTaskService},
			&adapters.ScheduledTaskAccessAdapter{AppService: d.appService, ProjectService: d.projectService},
		),
		notifSettingsHandler: notifsmod.NewNotificationSettingsHandler(&adapters.NotifSettingsServiceAdapter{
			Svc:     d.notifSettingsService,
			OrgRepo: d.orgRepo,
		}),
		updaterHandler:      system.NewUpdaterHandler(&adapters.UpdaterServiceAdapter{Svc: updaterService}),
		notificationHandler: notifsmod.NewNotificationHandler(&adapters.NotificationServiceAdapter{}),
		systemHandler: system.NewSystemHandler(
			&adapters.SystemServiceAdapter{Svc: system.NewSystemService()},
			&adapters.SystemAdminFinderAdapter{UserRepo: d.userRepo},
			&adapters.SystemTokenGenAdapter{TokenService: d.tokenService},
		),
		metricsHandler:       monitoring.NewMetricsHandler(&adapters.MetricsServiceAdapter{Svc: metricsService}),
		logHandler:           monitoring.NewLogHandler(&adapters.LogServiceAdapter{Svc: logService}),
		exampleHandler:       monitoring.NewExampleHandler(&adapters.ExampleServiceAdapter{Svc: monitoring.NewExampleService()}),
		serviceLogsWSHandler: wsHandler,
		ingressHandler:       ingress.NewIngressHandler(ingressAdapter),
	}
}
