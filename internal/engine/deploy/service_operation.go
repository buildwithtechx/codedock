package deploy

import "context"

type serviceOperationKey struct{}

func (d *Deployer) BeginServiceOperation(ctx context.Context, id string) (context.Context, func(), error) {
	release, err := d.serviceOperation(ctx, id)
	if err != nil {
		return ctx, nil, err
	}
	return context.WithValue(ctx, serviceOperationKey{}, id), release, nil
}

func (d *Deployer) serviceOperation(ctx context.Context, id string) (func(), error) {
	if owned, ok := ctx.Value(serviceOperationKey{}).(string); ok && owned == id {
		return func() {}, nil
	}
	return d.containerManager.acquireService(id)
}
