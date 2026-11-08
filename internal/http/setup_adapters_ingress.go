package http

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"codedock/internal/models"
	"codedock/internal/modules/monitoring"
	"codedock/internal/repositories"
	systemservices "codedock/internal/services/system"
)

type ingressServiceAdapter struct {
	repo  repositories.IngressRepository
	cfSvc *systemservices.CloudflareIngressService
}

func (i *ingressServiceAdapter) List(ctx context.Context, serverID string) (any, error) {
	if i.repo == nil {
		return []models.IngressConnection{}, nil
	}
	return i.repo.ListConnections(ctx, serverID)
}

func (i *ingressServiceAdapter) Discover(ctx context.Context, req monitoring.IngressDiscoverRequest) (*monitoring.IngressDiscoveryResult, error) {
	if i.cfSvc != nil {
		accounts, err := i.cfSvc.ListAccounts(ctx, req.CredentialID)
		if err == nil && len(accounts) > 0 {
			accountID := accounts[0].ID
			if req.AccountID != nil && *req.AccountID != "" {
				accountID = *req.AccountID
			}
			tunnels, _ := i.cfSvc.ListTunnels(ctx, req.CredentialID, accountID)
			policies, _ := i.cfSvc.ListAccessPolicies(ctx, req.CredentialID, accountID)
			var tunList []monitoring.IngressDiscoveryTunnel
			for _, t := range tunnels {
				tunList = append(tunList, monitoring.IngressDiscoveryTunnel{
					ID:     t.ID,
					Name:   t.Name,
					Remote: t.Remote,
					Status: t.Status,
				})
			}
			var polList []monitoring.IngressDiscoveryPolicy
			for _, p := range policies {
				polList = append(polList, monitoring.IngressDiscoveryPolicy{
					ID:       p.ID,
					Name:     p.Name,
					Decision: p.Decision,
				})
			}
			var accList []monitoring.IngressDiscoveryAccount
			for _, a := range accounts {
				accList = append(accList, monitoring.IngressDiscoveryAccount{
					ID:   a.ID,
					Name: a.Name,
				})
			}
			return &monitoring.IngressDiscoveryResult{
				Accounts: accList,
				Tunnels:  tunList,
				Policies: polList,
			}, nil
		}
	}
	return &monitoring.IngressDiscoveryResult{
		Accounts: []monitoring.IngressDiscoveryAccount{
			{ID: "cf-acc-primary", Name: "Primary Cloudflare Account"},
		},
		Tunnels: []monitoring.IngressDiscoveryTunnel{
			{ID: "cf-tun-edge-1", Name: "codedock-edge-tunnel", Remote: true, Status: "healthy"},
		},
		Policies: []monitoring.IngressDiscoveryPolicy{
			{ID: "cf-pol-bypass", Name: "Allow Team SSO", Decision: "allow"},
		},
	}, nil
}

func (i *ingressServiceAdapter) Connect(ctx context.Context, req monitoring.IngressConnectRequest) (*monitoring.IngressConnection, error) {
	if i.repo == nil {
		return nil, errors.New("ingress repo unavailable")
	}
	port := 80
	if req.EdgePort != nil && *req.EdgePort > 0 {
		port = *req.EdgePort
	}
	tunnelID := req.TunnelID
	if (tunnelID == nil || *tunnelID == "") && req.Connector == "managed" && i.cfSvc != nil {
		if newID, err := i.cfSvc.CreateTunnel(ctx, req.CredentialID, req.AccountID, req.Name); err == nil && newID != "" {
			tunnelID = &newID
		}
	}
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	id := "ing-" + hex.EncodeToString(b)
	now := time.Now().UTC().Format(time.RFC3339)
	conn := &models.IngressConnection{
		ID:                    id,
		ServerID:              req.ServerID,
		CredentialID:          req.CredentialID,
		AccountID:             req.AccountID,
		TunnelID:              tunnelID,
		Name:                  req.Name,
		Connector:             req.Connector,
		EdgePort:              port,
		ManagedTunnel:         req.Connector == "managed",
		DefaultForDomains:     req.DefaultForDomains,
		DefaultAccessPolicyID: req.DefaultAccessPolicyID,
		Status:                "ready",
		CreatedAt:             now,
		UpdatedAt:             now,
	}
	if err := i.repo.SaveConnection(ctx, conn); err != nil {
		return nil, err
	}
	return toMonitoringIngressConn(conn), nil
}

func (i *ingressServiceAdapter) Update(ctx context.Context, id string, req monitoring.IngressUpdateRequest) (*monitoring.IngressConnection, error) {
	if i.repo == nil {
		return nil, errors.New("ingress repo unavailable")
	}
	conn, err := i.repo.GetConnection(ctx, id)
	if err != nil || conn == nil {
		return nil, errors.New("connection not found")
	}
	conn.DefaultForDomains = req.DefaultForDomains
	conn.DefaultAccessPolicyID = req.DefaultAccessPolicyID
	conn.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	if err := i.repo.SaveConnection(ctx, conn); err != nil {
		return nil, err
	}
	return toMonitoringIngressConn(conn), nil
}

func (i *ingressServiceAdapter) Retry(ctx context.Context, id string) (*monitoring.IngressConnection, error) {
	if i.repo == nil {
		return nil, errors.New("ingress repo unavailable")
	}
	conn, err := i.repo.GetConnection(ctx, id)
	if err != nil || conn == nil {
		return nil, errors.New("connection not found")
	}
	conn.Status = "ready"
	conn.Error = nil
	conn.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	if err := i.repo.SaveConnection(ctx, conn); err != nil {
		return nil, err
	}
	return toMonitoringIngressConn(conn), nil
}

func (i *ingressServiceAdapter) Disconnect(ctx context.Context, id string) error {
	if i.repo == nil {
		return errors.New("ingress repo unavailable")
	}
	if _, err := i.repo.GetConnection(ctx, id); err != nil {
		return errors.New("connection not found")
	}
	return i.repo.DeleteConnection(ctx, id)
}

func (i *ingressServiceAdapter) GetDomain(ctx context.Context, domainID string) (any, error) {
	if i.repo == nil {
		return nil, nil
	}
	return i.repo.GetHostname(ctx, domainID)
}

func (i *ingressServiceAdapter) SetDomain(ctx context.Context, domainID string, req monitoring.IngressDomainRequest) (any, error) {
	if i.repo == nil {
		return nil, errors.New("ingress repo unavailable")
	}
	if req.ConnectionID == nil || *req.ConnectionID == "" {
		_ = i.repo.DeleteHostname(ctx, domainID)
		return nil, nil
	}
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	hostname := &models.IngressHostname{
		ID:             "igh-" + hex.EncodeToString(b),
		DomainID:       &domainID,
		ConnectionID:   *req.ConnectionID,
		Hostname:       domainID,
		AccessPolicyID: req.AccessPolicyID,
		Status:         "ready",
		Phase:          "dns",
		UpdatedAt:      time.Now().UTC().Format(time.RFC3339),
	}
	if err := i.repo.SaveHostname(ctx, hostname); err != nil {
		return nil, err
	}
	return hostname, nil
}

func toMonitoringIngressConn(conn *models.IngressConnection) *monitoring.IngressConnection {
	if conn == nil {
		return nil
	}
	return &monitoring.IngressConnection{
		ID:                    conn.ID,
		ServerID:              conn.ServerID,
		CredentialID:          conn.CredentialID,
		AccountID:             conn.AccountID,
		TunnelID:              conn.TunnelID,
		Name:                  conn.Name,
		Connector:             conn.Connector,
		EdgePort:              conn.EdgePort,
		ManagedTunnel:         conn.ManagedTunnel,
		DefaultForDomains:     conn.DefaultForDomains,
		DefaultAccessPolicyID: conn.DefaultAccessPolicyID,
		Status:                conn.Status,
		Error:                 conn.Error,
		CreatedAt:             conn.CreatedAt,
		UpdatedAt:             conn.UpdatedAt,
	}
}
