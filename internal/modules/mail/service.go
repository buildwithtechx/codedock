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
