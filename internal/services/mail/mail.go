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
	"time"

	"codedock/internal/models"
	"codedock/internal/repositories"
)

var ErrMailNotFound = errors.New("not found")
var ErrMailExists = errors.New("already exists")

type MailService struct {
	repo   repositories.MailRepository
	client *http.Client
}

func NewMailService(repo repositories.MailRepository) *MailService {
	return &MailService{repo: repo, client: &http.Client{Timeout: 15 * time.Second}}
}

func (s *MailService) CreateServer(ctx context.Context, req *models.CreateMailServerRequest) (*models.MailServer, error) {
	if req.OrganizationID == "" || req.Domain == "" || req.MaildURL == "" || req.APIToken == "" {
		return nil, fmt.Errorf("organization, domain, maild url, and api token are required")
	}
	server := &models.MailServer{
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

func (s *MailService) ListServers(ctx context.Context, orgID string) ([]*models.MailServer, error) {
	return s.repo.ListByOrganization(ctx, orgID)
}

func (s *MailService) GetServer(ctx context.Context, id string) (*models.MailServer, error) {
	server, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, mapRepoError(err)
	}
	return server, nil
}

func (s *MailService) DeleteServer(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return mapRepoError(err)
	}
	return nil
}

func (s *MailService) ListDomains(ctx context.Context, id string) ([]models.MailDomain, error) {
	server, err := s.resolve(ctx, id)
	if err != nil {
		return nil, err
	}
	var domains []models.MailDomain
	if err := s.do(ctx, server, http.MethodGet, "/api/domains", nil, &domains); err != nil {
		return nil, err
	}
	if domains == nil {
		domains = []models.MailDomain{}
	}
	return domains, nil
}

func (s *MailService) CreateDomain(ctx context.Context, id, name string) (*models.MailDomain, error) {
	server, err := s.resolve(ctx, id)
	if err != nil {
		return nil, err
	}
	var domain models.MailDomain
	if err := s.do(ctx, server, http.MethodPost, "/api/domains", models.CreateMailDomainRequest{Name: name}, &domain); err != nil {
		return nil, err
	}
	return &domain, nil
}

func (s *MailService) DeleteDomain(ctx context.Context, id, domainID string) error {
	server, err := s.resolve(ctx, id)
	if err != nil {
		return err
	}
	return s.do(ctx, server, http.MethodDelete, "/api/domains/"+domainID, nil, nil)
}

func (s *MailService) ListMailboxes(ctx context.Context, id, domainID string) ([]models.Mailbox, error) {
	server, err := s.resolve(ctx, id)
	if err != nil {
		return nil, err
	}
	var mailboxes []models.Mailbox
	if err := s.do(ctx, server, http.MethodGet, "/api/domains/"+domainID+"/mailboxes", nil, &mailboxes); err != nil {
		return nil, err
	}
	if mailboxes == nil {
		mailboxes = []models.Mailbox{}
	}
	return mailboxes, nil
}

func (s *MailService) CreateMailbox(ctx context.Context, id, domainID, localPart string) (*models.MailboxCredentials, error) {
	server, err := s.resolve(ctx, id)
	if err != nil {
		return nil, err
	}
	var credentials models.MailboxCredentials
	if err := s.do(ctx, server, http.MethodPost, "/api/domains/"+domainID+"/mailboxes", models.CreateMailboxRequest{LocalPart: localPart}, &credentials); err != nil {
		return nil, err
	}
	return &credentials, nil
}

func (s *MailService) DeleteMailbox(ctx context.Context, id, domainID, mailboxID string) error {
	server, err := s.resolve(ctx, id)
	if err != nil {
		return err
	}
	return s.do(ctx, server, http.MethodDelete, "/api/domains/"+domainID+"/mailboxes/"+mailboxID, nil, nil)
}

func (s *MailService) CreateDKIM(ctx context.Context, id, domainID, selector, algorithm string) (*models.DKIMKey, error) {
	server, err := s.resolve(ctx, id)
	if err != nil {
		return nil, err
	}
	var key models.DKIMKey
	if err := s.do(ctx, server, http.MethodPost, "/api/domains/"+domainID+"/dkim", models.CreateDKIMRequest{Selector: selector, Algorithm: algorithm}, &key); err != nil {
		return nil, err
	}
	return &key, nil
}

func (s *MailService) GetDNS(ctx context.Context, id, domainID string) (*models.DNSResponse, error) {
	server, err := s.resolve(ctx, id)
	if err != nil {
		return nil, err
	}
	var response models.DNSResponse
	if err := s.do(ctx, server, http.MethodGet, "/api/domains/"+domainID+"/dns", nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *MailService) GetChecks(ctx context.Context, id, domainID, ip string) (*models.MailCheckReport, error) {
	server, err := s.resolve(ctx, id)
	if err != nil {
		return nil, err
	}
	path := "/api/domains/" + domainID + "/checks"
	if ip != "" {
		path += "?ip=" + ip
	}
	var report models.MailCheckReport
	if err := s.do(ctx, server, http.MethodGet, path, nil, &report); err != nil {
		return nil, err
	}
	return &report, nil
}

func (s *MailService) ListCerts(ctx context.Context, id string) ([]models.MailCert, error) {
	server, err := s.resolve(ctx, id)
	if err != nil {
		return nil, err
	}
	var certs []models.MailCert
	if err := s.do(ctx, server, http.MethodGet, "/api/certs", nil, &certs); err != nil {
		return nil, err
	}
	if certs == nil {
		certs = []models.MailCert{}
	}
	return certs, nil
}

func (s *MailService) ProbeCert(ctx context.Context, id, host string, port int) (*models.MailCert, error) {
	server, err := s.resolve(ctx, id)
	if err != nil {
		return nil, err
	}
	var cert models.MailCert
	if err := s.do(ctx, server, http.MethodPost, "/api/certs/probe", models.ProbeMailCertRequest{Host: host, Port: port}, &cert); err != nil {
		return nil, err
	}
	return &cert, nil
}

func (s *MailService) resolve(ctx context.Context, id string) (*models.MailServer, error) {
	server, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, mapRepoError(err)
	}
	return server, nil
}

func (s *MailService) do(ctx context.Context, server *models.MailServer, method, path string, body any, out any) error {
	var reader io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("encode mail request: %w", err)
		}
		reader = bytes.NewReader(encoded)
	}
	req, err := http.NewRequestWithContext(ctx, method, server.MaildURL+path, reader)
	if err != nil {
		return fmt.Errorf("build mail request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+server.APIToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("mail server unreachable: %w", err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return fmt.Errorf("read mail response: %w", err)
	}
	if resp.StatusCode == http.StatusNoContent {
		return nil
	}
	if resp.StatusCode >= 400 {
		return mapStatusError(resp.StatusCode, raw)
	}
	if out != nil && len(raw) > 0 {
		if err := json.Unmarshal(raw, out); err != nil {
			return fmt.Errorf("decode mail response: %w", err)
		}
	}
	return nil
}

func mapStatusError(status int, raw []byte) error {
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

func mapRepoError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrMailNotFound
	}
	return err
}
