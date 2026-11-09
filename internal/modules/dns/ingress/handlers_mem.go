package ingress

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

type inMemoryIngressService struct {
	mu          sync.RWMutex
	connections map[string]*IngressConnection
	domains     map[string]*IngressHostname
}

func newInMemoryIngressService() *inMemoryIngressService {
	return &inMemoryIngressService{
		connections: make(map[string]*IngressConnection),
		domains:     make(map[string]*IngressHostname),
	}
}

func (s *inMemoryIngressService) List(ctx context.Context, serverID string) (any, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var list []*IngressConnection
	for _, c := range s.connections {
		if serverID == "" || c.ServerID == serverID {
			list = append(list, c)
		}
	}
	if list == nil {
		list = make([]*IngressConnection, 0)
	}
	return list, nil
}

func (s *inMemoryIngressService) Discover(ctx context.Context, req IngressDiscoverRequest) (*IngressDiscoveryResult, error) {
	return &IngressDiscoveryResult{
		Accounts: []IngressDiscoveryAccount{
			{ID: "cf-acc-primary", Name: "Primary Cloudflare Account"},
		},
		Tunnels: []IngressDiscoveryTunnel{
			{ID: "cf-tun-edge-1", Name: "codedock-edge-tunnel", Remote: true, Status: "healthy"},
		},
		Policies: []IngressDiscoveryPolicy{
			{ID: "cf-pol-bypass", Name: "Allow Team SSO", Decision: "allow"},
		},
	}, nil
}

func (s *inMemoryIngressService) Connect(ctx context.Context, req IngressConnectRequest) (*IngressConnection, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	b := make([]byte, 8)
	_, _ = rand.Read(b)
	id := "ing-" + hex.EncodeToString(b)
	now := time.Now().UTC().Format(time.RFC3339)

	port := 80
	if req.EdgePort != nil && *req.EdgePort > 0 {
		port = *req.EdgePort
	}

	conn := &IngressConnection{
		ID:                    id,
		ServerID:              req.ServerID,
		CredentialID:          req.CredentialID,
		AccountID:             req.AccountID,
		TunnelID:              req.TunnelID,
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
	s.connections[id] = conn
	return conn, nil
}

func (s *inMemoryIngressService) Update(ctx context.Context, id string, req IngressUpdateRequest) (*IngressConnection, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	conn, ok := s.connections[id]
	if !ok {
		return nil, errors.New("connection not found")
	}
	conn.DefaultForDomains = req.DefaultForDomains
	conn.DefaultAccessPolicyID = req.DefaultAccessPolicyID
	conn.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return conn, nil
}

func (s *inMemoryIngressService) Retry(ctx context.Context, id string) (*IngressConnection, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	conn, ok := s.connections[id]
	if !ok {
		return nil, errors.New("connection not found")
	}
	conn.Status = "ready"
	conn.Error = nil
	conn.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return conn, nil
}

func (s *inMemoryIngressService) Disconnect(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.connections[id]; !ok {
		return errors.New("connection not found")
	}
	delete(s.connections, id)
	return nil
}

func (s *inMemoryIngressService) GetDomain(ctx context.Context, domainID string) (any, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	host, ok := s.domains[domainID]
	if !ok {
		return nil, nil
	}
	return host, nil
}

func (s *inMemoryIngressService) SetDomain(ctx context.Context, domainID string, req IngressDomainRequest) (any, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if req.ConnectionID == nil || *req.ConnectionID == "" {
		delete(s.domains, domainID)
		return nil, nil
	}
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	host := &IngressHostname{
		ID:             "igh-" + hex.EncodeToString(b),
		DomainID:       &domainID,
		ConnectionID:   *req.ConnectionID,
		Hostname:       domainID,
		AccessPolicyID: req.AccessPolicyID,
		Status:         "ready",
		Phase:          "dns",
		UpdatedAt:      time.Now().UTC().Format(time.RFC3339),
	}
	s.domains[domainID] = host
	return host, nil
}
