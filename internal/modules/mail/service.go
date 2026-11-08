package mail

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Service struct {
	repo   Repository
	client *http.Client
}

func NewService(repo Repository) *Service {
	return &Service{
		repo:   repo,
		client: &http.Client{Timeout: DefaultHTTPTimeout},
	}
}

func (s *Service) CreateServer(ctx context.Context, req *CreateMailServerRequest) (*MailServer, error) {
	if req.OrganizationID == "" || req.Domain == "" || req.MaildURL == "" || req.APIToken == "" {
		return nil, ErrServerRequired
	}
	server := &MailServer{
		OrganizationID: req.OrganizationID,
		ServerID:       req.ServerID,
		Domain:         strings.TrimSpace(req.Domain),
		MaildURL:       strings.TrimSuffix(strings.TrimSpace(req.MaildURL), "/"),
		WebmailURL:     strings.TrimSpace(req.WebmailURL),
		APIToken:       req.APIToken,
	}
	if err := s.repo.Create(ctx, server); err != nil {
		return nil, err
	}
	return server, nil
}

func (s *Service) ListServers(ctx context.Context, orgID string) ([]*MailServer, error) {
	return s.repo.ListByOrganization(ctx, orgID)
}

func (s *Service) GetServer(ctx context.Context, id string) (*MailServer, error) {
	server, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, s.mapRepoError(err)
	}
	return server, nil
}

func (s *Service) DeleteServer(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return s.mapRepoError(err)
	}
	return nil
}

func (s *Service) ListDomains(ctx context.Context, id string) ([]MailDomain, error) {
	server, err := s.resolve(ctx, id)
	if err != nil {
		return nil, err
	}
	var domains []MailDomain
	if err := s.do(ctx, server, http.MethodGet, "/api/domains", nil, &domains); err != nil {
		return nil, err
	}
	if domains == nil {
		domains = []MailDomain{}
	}
	return domains, nil
}

func (s *Service) CreateDomain(ctx context.Context, id, name string) (*MailDomain, error) {
	server, err := s.resolve(ctx, id)
	if err != nil {
		return nil, err
	}
	var domain MailDomain
	if err := s.do(ctx, server, http.MethodPost, "/api/domains", CreateMailDomainRequest{Name: name}, &domain); err != nil {
		return nil, err
	}
	return &domain, nil
}

func (s *Service) DeleteDomain(ctx context.Context, id, domainID string) error {
	server, err := s.resolve(ctx, id)
	if err != nil {
		return err
	}
	return s.do(ctx, server, http.MethodDelete, "/api/domains/"+domainID, nil, nil)
}

func (s *Service) ListMailboxes(ctx context.Context, id, domainID string) ([]Mailbox, error) {
	server, err := s.resolve(ctx, id)
	if err != nil {
		return nil, err
	}
	var mailboxes []Mailbox
	if err := s.do(ctx, server, http.MethodGet, "/api/domains/"+domainID+"/mailboxes", nil, &mailboxes); err != nil {
		return nil, err
	}
	if mailboxes == nil {
		mailboxes = []Mailbox{}
	}
	return mailboxes, nil
}

func (s *Service) CreateMailbox(ctx context.Context, id, domainID, localPart string) (*MailboxCredentials, error) {
	server, err := s.resolve(ctx, id)
	if err != nil {
		return nil, err
	}
	var credentials MailboxCredentials
	if err := s.do(ctx, server, http.MethodPost, "/api/domains/"+domainID+"/mailboxes", CreateMailboxRequest{LocalPart: localPart}, &credentials); err != nil {
		return nil, err
	}
	return &credentials, nil
}

func (s *Service) DeleteMailbox(ctx context.Context, id, domainID, mailboxID string) error {
	server, err := s.resolve(ctx, id)
	if err != nil {
		return err
	}
	return s.do(ctx, server, http.MethodDelete, "/api/domains/"+domainID+"/mailboxes/"+mailboxID, nil, nil)
}

func (s *Service) CreateDKIM(ctx context.Context, id, domainID, selector, algorithm string) (*DKIMKey, error) {
	server, err := s.resolve(ctx, id)
	if err != nil {
		return nil, err
	}
	var key DKIMKey
	if err := s.do(ctx, server, http.MethodPost, "/api/domains/"+domainID+"/dkim", CreateDKIMRequest{Selector: selector, Algorithm: algorithm}, &key); err != nil {
		return nil, err
	}
	return &key, nil
}

func (s *Service) GetDNS(ctx context.Context, id, domainID string) (*DNSResponse, error) {
	server, err := s.resolve(ctx, id)
	if err != nil {
		return nil, err
	}
	var response DNSResponse
	if err := s.do(ctx, server, http.MethodGet, "/api/domains/"+domainID+"/dns", nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) GetChecks(ctx context.Context, id, domainID, ip string) (*MailCheckReport, error) {
	server, err := s.resolve(ctx, id)
	if err != nil {
		return nil, err
	}
	path := "/api/domains/" + domainID + "/checks"
	if ip != "" {
		path += "?ip=" + ip
	}
	var report MailCheckReport
	if err := s.do(ctx, server, http.MethodGet, path, nil, &report); err != nil {
		return nil, err
	}
	return &report, nil
}

func (s *Service) ListCerts(ctx context.Context, id string) ([]MailCert, error) {
	server, err := s.resolve(ctx, id)
	if err != nil {
		return nil, err
	}
	var certs []MailCert
	if err := s.do(ctx, server, http.MethodGet, "/api/certs", nil, &certs); err != nil {
		return nil, err
	}
	if certs == nil {
		certs = []MailCert{}
	}
	return certs, nil
}

func (s *Service) ProbeCert(ctx context.Context, id, host string, port int) (*MailCert, error) {
	server, err := s.resolve(ctx, id)
	if err != nil {
		return nil, err
	}
	var cert MailCert
	if err := s.do(ctx, server, http.MethodPost, "/api/certs/probe", ProbeMailCertRequest{Host: host, Port: port}, &cert); err != nil {
		return nil, err
	}
	return &cert, nil
}

func (s *Service) GetFail2banStatus(ctx context.Context, id, jail string) (*Fail2banStatusResponse, error) {
	server, err := s.resolve(ctx, id)
	if err != nil {
		return nil, err
	}
	path := "/api/fail2ban/status"
	if jail != "" {
		if err := ValidateJailName(jail); err != nil {
			return nil, err
		}
		path += "?jail=" + jail
	}
	var resp Fail2banStatusResponse
	if err := s.do(ctx, server, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	resp.CanManage = true
	return &resp, nil
}

func (s *Service) BanIP(ctx context.Context, id, ip, jail string) (*Fail2banActionResponse, error) {
	if err := ValidateIP(ip); err != nil {
		return nil, err
	}
	if err := ValidateJailName(jail); err != nil {
		return nil, err
	}
	server, err := s.resolve(ctx, id)
	if err != nil {
		return nil, err
	}
	var resp Fail2banActionResponse
	req := Fail2banActionRequest{IP: ip, Jail: jail}
	if err := s.do(ctx, server, http.MethodPost, "/api/fail2ban/ban", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *Service) UnbanIP(ctx context.Context, id, ip, jail string) (*Fail2banActionResponse, error) {
	if err := ValidateIP(ip); err != nil {
		return nil, err
	}
	if jail != "" {
		if err := ValidateJailName(jail); err != nil {
			return nil, err
		}
	}
	server, err := s.resolve(ctx, id)
	if err != nil {
		return nil, err
	}
	var resp Fail2banActionResponse
	req := Fail2banActionRequest{IP: ip, Jail: jail}
	if err := s.do(ctx, server, http.MethodPost, "/api/fail2ban/unban", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *Service) ListAliases(ctx context.Context, id, domainID string) ([]MailAlias, error) {
	server, err := s.resolve(ctx, id)
	if err != nil {
		return nil, err
	}
	var aliases []MailAlias
	if err := s.do(ctx, server, http.MethodGet, "/api/domains/"+domainID+"/aliases", nil, &aliases); err != nil {
		return nil, err
	}
	if aliases == nil {
		aliases = []MailAlias{}
	}
	return aliases, nil
}

func (s *Service) CreateAlias(ctx context.Context, id, domainID string, req CreateMailAliasRequest) (*MailAlias, error) {
	server, err := s.resolve(ctx, id)
	if err != nil {
		return nil, err
	}
	req.Source = strings.TrimSpace(req.Source)
	req.Destination = strings.TrimSpace(req.Destination)
	if req.Source == "" || req.Destination == "" {
		return nil, fmt.Errorf("source and destination are required")
	}
	var alias MailAlias
	if err := s.do(ctx, server, http.MethodPost, "/api/domains/"+domainID+"/aliases", req, &alias); err != nil {
		return nil, err
	}
	return &alias, nil
}

func (s *Service) DeleteAlias(ctx context.Context, id, domainID, aliasID string) error {
	server, err := s.resolve(ctx, id)
	if err != nil {
		return err
	}
	return s.do(ctx, server, http.MethodDelete, "/api/domains/"+domainID+"/aliases/"+aliasID, nil, nil)
}

func (s *Service) resolve(ctx context.Context, id string) (*MailServer, error) {
	server, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, s.mapRepoError(err)
	}
	return server, nil
}

func v1Path(path string) string {
	if strings.HasPrefix(path, "/api/") && !strings.HasPrefix(path, "/api/v1/") {
		return "/api/v1" + strings.TrimPrefix(path, "/api")
	}
	return path
}

func (s *Service) do(ctx context.Context, server *MailServer, method, path string, body any, out any) error {
	var encoded []byte
	if body != nil {
		var err error
		encoded, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("encode mail request: %w", err)
		}
	}

	exec := func(p string) (*http.Response, []byte, error) {
		var reader io.Reader
		if encoded != nil {
			reader = bytes.NewReader(encoded)
		}
		req, err := http.NewRequestWithContext(ctx, method, server.MaildURL+p, reader)
		if err != nil {
			return nil, nil, fmt.Errorf("build mail request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+server.APIToken)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		resp, err := s.client.Do(req)
		if err != nil {
			return nil, nil, fmt.Errorf("mail server unreachable: %w", err)
		}
		defer resp.Body.Close()
		raw, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
		if err != nil {
			return nil, nil, fmt.Errorf("read mail response: %w", err)
		}
		return resp, raw, nil
	}

	targetPath := v1Path(path)
	resp, raw, err := exec(targetPath)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusNotFound && targetPath != path {
		if fbResp, fbRaw, fbErr := exec(path); fbErr == nil && fbResp.StatusCode != http.StatusNotFound {
			resp = fbResp
			raw = fbRaw
		}
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}
	if resp.StatusCode >= 400 {
		return s.mapStatusError(resp.StatusCode, raw)
	}
	if out != nil && len(raw) > 0 {
		if err := json.Unmarshal(raw, out); err != nil {
			return fmt.Errorf("decode mail response: %w", err)
		}
	}
	return nil
}

func (s *Service) mapStatusError(status int, raw []byte) error {
	var payload struct {
		Error string `json:"error"`
	}
	message := strings.TrimSpace(string(raw))
	if err := json.Unmarshal(raw, &payload); err == nil && payload.Error != "" {
		message = payload.Error
	}
	switch status {
	case http.StatusNotFound:
		return ErrMailNotFound
	case http.StatusConflict:
		return ErrMailExists
	default:
		return fmt.Errorf("maild: status %d: %s", status, message)
	}
}

func (s *Service) mapRepoError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrMailNotFound
	}
	return err
}
