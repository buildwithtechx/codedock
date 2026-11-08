package http

import (
	"context"

	"codedock/internal/models"
	"codedock/internal/modules/servers"
	systemservices "codedock/internal/services/system"
)

type managedServiceAdapter struct {
	svc systemservices.ManagedService
}

func (a *managedServiceAdapter) Catalog() servers.ManagedCatalog {
	if a.svc == nil {
		return servers.ManagedCatalog{}
	}
	c := a.svc.Catalog()
	tiers := make([]servers.ManagedTier, len(c.Tiers))
	for i, t := range c.Tiers {
		tiers[i] = servers.ManagedTier{Name: t.Name, Cores: t.Cores, MemoryGB: t.MemoryGB, DiskGB: t.DiskGB, MonthlyPrice: t.MonthlyPrice}
	}
	return servers.ManagedCatalog{
		Provider: servers.ManagedProvider(c.Provider),
		Tiers:    tiers,
		Regions:  c.Regions,
		Images:   c.Images,
	}
}

func (a *managedServiceAdapter) Offers(ctx context.Context, userID, orgID, credentialID string) (*servers.ManagedOffers, error) {
	if a.svc == nil {
		return &servers.ManagedOffers{}, nil
	}
	res, err := a.svc.Offers(ctx, userID, orgID, credentialID)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	locs := make([]servers.ManagedLocation, len(res.Locations))
	for i, l := range res.Locations {
		locs[i] = servers.ManagedLocation{Name: l.Name, City: l.City, Country: l.Country}
	}
	offers := make([]servers.ManagedOffer, len(res.Offers))
	for i, o := range res.Offers {
		offers[i] = servers.ManagedOffer{
			Tier:   servers.ManagedTier{Name: o.Tier.Name, Cores: o.Tier.Cores, MemoryGB: o.Tier.MemoryGB, DiskGB: o.Tier.DiskGB, MonthlyPrice: o.Tier.MonthlyPrice},
			Prices: o.Prices,
		}
	}
	return &servers.ManagedOffers{
		Provider:  servers.ManagedProvider(res.Provider),
		Locations: locs,
		Offers:    offers,
		Images:    res.Images,
	}, nil
}

func (a *managedServiceAdapter) GetOperation(ctx context.Context, userID, operationID string) (any, error) {
	if a.svc == nil {
		return nil, nil
	}
	return a.svc.GetOperation(ctx, userID, operationID)
}

func (a *managedServiceAdapter) StreamOperation(ctx context.Context, userID, operationID string, emit func(any) error) error {
	if a.svc == nil {
		return nil
	}
	return a.svc.StreamOperation(ctx, userID, operationID, func(op *models.Operation) error {
		return emit(op)
	})
}

func (a *managedServiceAdapter) CreateCredential(ctx context.Context, userID, orgID string, req servers.CreateManagedCredentialRequest) (*servers.ManagedCredential, error) {
	if a.svc == nil {
		return nil, nil
	}
	res, err := a.svc.CreateCredential(ctx, userID, orgID, models.CreateManagedCredentialRequest{
		Provider: models.ManagedProvider(req.Provider),
		Label:    req.Label,
		Token:    req.Token,
	})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return &servers.ManagedCredential{
		ID:             res.ID,
		OrganizationID: res.OrganizationID,
		Provider:       servers.ManagedProvider(res.Provider),
		Label:          res.Label,
		CreatedAt:      res.CreatedAt,
		UpdatedAt:      res.UpdatedAt,
	}, nil
}

func (a *managedServiceAdapter) ListCredentials(ctx context.Context, userID, orgID string) ([]*servers.ManagedCredential, error) {
	if a.svc == nil {
		return nil, nil
	}
	list, err := a.svc.ListCredentials(ctx, userID, orgID)
	if err != nil {
		return nil, err
	}
	result := make([]*servers.ManagedCredential, len(list))
	for i, c := range list {
		result[i] = &servers.ManagedCredential{
			ID:             c.ID,
			OrganizationID: c.OrganizationID,
			Provider:       servers.ManagedProvider(c.Provider),
			Label:          c.Label,
			CreatedAt:      c.CreatedAt,
			UpdatedAt:      c.UpdatedAt,
		}
	}
	return result, nil
}

func (a *managedServiceAdapter) DeleteCredential(ctx context.Context, userID, orgID, credentialID string) error {
	if a.svc == nil {
		return nil
	}
	return a.svc.DeleteCredential(ctx, userID, orgID, credentialID)
}

func (a *managedServiceAdapter) GetQuota(ctx context.Context, userID, orgID string) (*servers.ManagedQuota, error) {
	if a.svc == nil {
		return nil, nil
	}
	q, err := a.svc.GetQuota(ctx, userID, orgID)
	if err != nil {
		return nil, err
	}
	if q == nil {
		return nil, nil
	}
	return &servers.ManagedQuota{
		OrganizationID: q.OrganizationID,
		MaxServers:     q.MaxServers,
		MaxMemoryGB:    q.MaxMemoryGB,
		UsedServers:    q.UsedServers,
		UsedMemoryGB:   q.UsedMemoryGB,
		UpdatedAt:      q.UpdatedAt,
	}, nil
}

func (a *managedServiceAdapter) SetQuota(ctx context.Context, userID, orgID string, req servers.ManagedQuotaRequest) (*servers.ManagedQuota, error) {
	if a.svc == nil {
		return nil, nil
	}
	q, err := a.svc.SetQuota(ctx, userID, orgID, models.ManagedQuotaRequest{
		MaxServers:  req.MaxServers,
		MaxMemoryGB: req.MaxMemoryGB,
	})
	if err != nil {
		return nil, err
	}
	if q == nil {
		return nil, nil
	}
	return &servers.ManagedQuota{
		OrganizationID: q.OrganizationID,
		MaxServers:     q.MaxServers,
		MaxMemoryGB:    q.MaxMemoryGB,
		UsedServers:    q.UsedServers,
		UsedMemoryGB:   q.UsedMemoryGB,
		UpdatedAt:      q.UpdatedAt,
	}, nil
}

func (a *managedServiceAdapter) ReviewProvision(ctx context.Context, userID, orgID string, req servers.ReviewManagedProvisionRequest) (*servers.ManagedReview, error) {
	if a.svc == nil {
		return nil, nil
	}
	r, err := a.svc.ReviewProvision(ctx, userID, orgID, models.ReviewManagedProvisionRequest{
		CredentialID: req.CredentialID,
		ServerID:     req.ServerID,
		Name:         req.Name,
		Region:       req.Region,
		ServerType:   req.ServerType,
		Image:        req.Image,
	})
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, nil
	}
	return &servers.ManagedReview{
		Review: r.Review,
		Plan:   r.Plan,
	}, nil
}

func (a *managedServiceAdapter) ReviewResize(ctx context.Context, userID, orgID, serverID string, req servers.ReviewManagedResizeRequest) (*servers.ManagedReview, error) {
	if a.svc == nil {
		return nil, nil
	}
	r, err := a.svc.ReviewResize(ctx, userID, orgID, serverID, models.ReviewManagedResizeRequest{ServerType: req.ServerType})
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, nil
	}
	return &servers.ManagedReview{
		Review: r.Review,
	}, nil
}

func (a *managedServiceAdapter) ReviewDelete(ctx context.Context, userID, orgID, serverID string) (*servers.ManagedReview, error) {
	if a.svc == nil {
		return nil, nil
	}
	r, err := a.svc.ReviewDelete(ctx, userID, orgID, serverID)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, nil
	}
	return &servers.ManagedReview{
		Review: r.Review,
	}, nil
}

func (a *managedServiceAdapter) ApplyOperation(ctx context.Context, userID, operationID, confirmation string) error {
	if a.svc == nil {
		return nil
	}
	return a.svc.ApplyOperation(ctx, userID, operationID, confirmation)
}

func (a *managedServiceAdapter) RefreshServer(ctx context.Context, userID, orgID, serverID string) (*servers.Server, error) {
	if a.svc == nil {
		return nil, nil
	}
	s, err := a.svc.RefreshServer(ctx, userID, orgID, serverID)
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, nil
	}
	return &servers.Server{
		ID:        s.ID,
		UserID:    s.UserID,
		Name:      s.Name,
		IPAddress: s.IPAddress,
		Status:    servers.ServerStatus(s.Status),
	}, nil
}
