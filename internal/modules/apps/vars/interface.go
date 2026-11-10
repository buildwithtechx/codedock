package vars

import "context"

// ServiceVarRepository defines database persistence for application environment variables.
type ServiceVarRepository interface {
	Create(ctx context.Context, v *Variable) error
	Update(ctx context.Context, v *Variable) error
	GetByID(ctx context.Context, id string) (*Variable, error)
	ListByService(ctx context.Context, serviceID string) ([]*Variable, error)
	Delete(ctx context.Context, id string) error
}

// VariableService provides business logic for injecting and managing environment variables.
type VariableService interface {
	CreateVariable(ctx context.Context, v *Variable) (*Variable, error)
	UpdateVariable(ctx context.Context, v *Variable) error
	GetVariable(ctx context.Context, id string) (*Variable, error)
	ListVariablesByService(ctx context.Context, serviceID string) ([]*Variable, error)
	DeleteVariable(ctx context.Context, id string) error
}
