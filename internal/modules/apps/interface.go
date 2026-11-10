package apps

import "context"

// AppServiceRepository defines persistence operations for app services and related resources.
type AppServiceRepository interface {
	Create(ctx context.Context, svc *AppService) error
	GetByID(ctx context.Context, id string) (*AppService, error)
	ListByEnvironment(ctx context.Context, environmentID string) ([]*AppService, error)
	ListByProject(ctx context.Context, projectID string) ([]*AppService, error)
	ListByOrganization(ctx context.Context, organizationID string, search string) ([]*AppService, error)
	ListAll(ctx context.Context, search string) ([]*AppService, error)
	SearchIDsByName(ctx context.Context, term string, limit int) ([]string, error)
	Update(ctx context.Context, svc *AppService) error
	Delete(ctx context.Context, id string) error
	CreateWebhook(ctx context.Context, w *Webhook) error
	ListWebhooksByService(ctx context.Context, serviceID string) ([]*Webhook, error)
	DeleteWebhook(ctx context.Context, id, serviceID string) error
	CreateLogDrain(ctx context.Context, d *LogDrain) error
	ListLogDrainsByService(ctx context.Context, serviceID string) ([]*LogDrain, error)
	DeleteLogDrain(ctx context.Context, id, serviceID string) error
}

// AccessChecker defines permission checks for app services within projects.
type AccessChecker interface {
	HasPermission(ctx context.Context, projectID, userID, role, permission string) bool
	AccessibleServiceIDs(ctx context.Context, userID, role string) (map[string]bool, error)
}

// EnvironmentLookup defines resolution of project IDs from environment IDs.
type EnvironmentLookup interface {
	GetProjectID(ctx context.Context, environmentID string) (string, error)
}

// LifecycleDeployer defines lifecycle deployment and restart actions for app services.
type LifecycleDeployer interface {
	StopAppService(ctx context.Context, appID string) error
	RestartAppService(ctx context.Context, appID string) error
	RedeployAppService(ctx context.Context, appID string) (any, error)
	ApplyProjectUpdates(ctx context.Context, projectID string) (any, error)
}

// RuntimeEngine defines runtime observability and execution operations.
type RuntimeEngine interface {
	Get(ctx context.Context, id string) (any, error)
	Operation(ctx context.Context, userID, id, opID string) (any, error)
	Review(ctx context.Context, userID, id string, req any) (any, error)
	ApplyForService(ctx context.Context, userID, id, opID, confirmation string) error
	Observe(ctx context.Context, id string) (any, error)
	Logs(ctx context.Context, id string) (any, error)
	Graph(ctx context.Context, id string) (any, error)
	Resources(ctx context.Context, id string) (any, error)
}

// Vault defines encryption and decryption for sensitive runtime credentials.
type Vault interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

// DatabaseRepository defines database querying for linked runtime services.
type DatabaseRepository interface {
	ListByProject(ctx context.Context, projectID string, search string) ([]*Database, error)
}

// AppServiceOperations defines high-level business logic for managing app services.
type AppServiceOperations interface {
	CreateAppService(ctx context.Context, svc *AppService) (*AppService, error)
	GetAppService(ctx context.Context, id string) (*AppService, error)
	ListAppServicesByEnvironment(ctx context.Context, environmentID string) ([]*AppService, error)
	ListAppServicesByProject(ctx context.Context, projectID string) ([]*AppService, error)
	ListAppServicesByOrganization(ctx context.Context, organizationID string, search string) ([]*AppService, error)
	ListAllAppServices(ctx context.Context, search string) ([]*AppService, error)
	UpdateAppService(ctx context.Context, svc *AppService) (*AppService, error)
	DeleteAppService(ctx context.Context, id string) error
	CreateWebhook(ctx context.Context, w *Webhook) error
	ListWebhooks(ctx context.Context, serviceID string) ([]*Webhook, error)
	DeleteWebhook(ctx context.Context, id, serviceID string) error
	CreateLogDrain(ctx context.Context, d *LogDrain) error
	ListLogDrains(ctx context.Context, serviceID string) ([]*LogDrain, error)
	DeleteLogDrain(ctx context.Context, id, serviceID string) error
}
