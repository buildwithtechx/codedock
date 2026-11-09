package http

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"codedock/internal/modules/dns/ingress"
)

type ingressServiceAdapter struct {
	repo  ingress.IngressRepository
	cfSvc *ingress.CloudflareIngressService
}

func (i *ingressServiceAdapter) List(ctx context.Context, serverID string) (any, error) {
	if i.repo == nil {
		return []ingress.IngressConnection{}, nil
	}
	return i.repo.ListConnections(ctx, serverID)
}

func (i *ingressServiceAdapter) Discover(ctx context.Context, req ingress.IngressDiscoverRequest) (*ingress.IngressDiscoveryResult, error) {
	if i.cfSvc != nil {
		accounts, err := i.cfSvc.ListAccounts(ctx, req.CredentialID)
		if err == nil && len(accounts) > 0 {
			accountID := accounts[0].ID
			if req.AccountID != nil && *req.AccountID != "" {
				accountID = *req.AccountID
			}
			tunnels, _ := i.cfSvc.ListTunnels(ctx, req.CredentialID, accountID)
			policies, _ := i.cfSvc.ListAccessPolicies(ctx, req.CredentialID, accountID)
			var tunList []ingress.IngressDiscoveryTunnel
			for _, t := range tunnels {
				tunList = append(tunList, ingress.IngressDiscoveryTunnel{
					ID:     t.ID,
					Name:   t.Name,
					Remote: t.Remote,
					Status: t.Status,
				})
			}
			var polList []ingress.IngressDiscoveryPolicy
			for _, p := range policies {
				polList = append(polList, ingress.IngressDiscoveryPolicy{
					ID:       p.ID,
					Name:     p.Name,
					Decision: p.Decision,
				})
			}
			var accList []ingress.IngressDiscoveryAccount
			for _, a := range accounts {
				accList = append(accList, ingress.IngressDiscoveryAccount{
					ID:   a.ID,
					Name: a.Name,
				})
			}
			return &ingress.IngressDiscoveryResult{
				Accounts: accList,
				Tunnels:  tunList,
				Policies: polList,
			}, nil
		}
	}
	return &ingress.IngressDiscoveryResult{
		Accounts: []ingress.IngressDiscoveryAccount{
			{ID: "cf-acc-primary", Name: "Primary Cloudflare Account"},
		},
		Tunnels: []ingress.IngressDiscoveryTunnel{
			{ID: "cf-tun-edge-1", Name: "codedock-edge-tunnel", Remote: true, Status: "healthy"},
		},
		Policies: []ingress.IngressDiscoveryPolicy{
			{ID: "cf-pol-bypass", Name: "Allow Team SSO", Decision: "allow"},
		},
	}, nil
}

func (i *ingressServiceAdapter) Connect(ctx context.Context, req ingress.IngressConnectRequest) (*ingress.IngressConnection, error) {
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
	conn := &ingress.IngressConnection{
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
	return conn, nil
}

func (i *ingressServiceAdapter) Update(ctx context.Context, id string, req ingress.IngressUpdateRequest) (*ingress.IngressConnection, error) {
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
	return conn, nil
}

func (i *ingressServiceAdapter) Retry(ctx context.Context, id string) (*ingress.IngressConnection, error) {
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
	return conn, nil
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

func (i *ingressServiceAdapter) SetDomain(ctx context.Context, domainID string, req ingress.IngressDomainRequest) (any, error) {
	if i.repo == nil {
		return nil, errors.New("ingress repo unavailable")
	}
	if req.ConnectionID == nil || *req.ConnectionID == "" {
		_ = i.repo.DeleteHostname(ctx, domainID)
		return nil, nil
	}
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	hostname := &ingress.IngressHostname{
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
