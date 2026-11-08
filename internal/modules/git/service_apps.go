package git

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type GitAppsService struct {
	repo GitAppRepository
}

func NewGitAppsService(repo GitAppRepository) *GitAppsService {
	return &GitAppsService{repo: repo}
}

type githubManifestConversionResponse struct {
	ID            int    `json:"id"`
	Slug          string `json:"slug"`
	ClientID      string `json:"client_id"`
	ClientSecret  string `json:"client_secret"`
	WebhookSecret string `json:"webhook_secret"`
	PEM           string `json:"pem"`
	HTMLURL       string `json:"html_url"`
	Name          string `json:"name"`
}

func (s *GitAppsService) ExchangeGithubManifestCode(ctx context.Context, code string) (*GithubApp, error) {
	if code == "" {
		return nil, errors.New("conversion code is required")
	}

	url := fmt.Sprintf("https://api.github.com/app-manifests/%s/conversions", code)
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github api error: %s", string(body))
	}

	var conversion githubManifestConversionResponse
	if err := json.NewDecoder(resp.Body).Decode(&conversion); err != nil {
		return nil, err
	}

	app := &GithubApp{
		ID:            uuid.NewString(),
		Name:          conversion.Name,
		AppID:         fmt.Sprintf("%d", conversion.ID),
		ClientID:      conversion.ClientID,
		ClientSecret:  conversion.ClientSecret,
		WebhookSecret: conversion.WebhookSecret,
		PrivateKey:    conversion.PEM,
		IsPublic:      false,
	}

	if err := s.repo.SaveGithubApp(ctx, app); err != nil {
		return nil, err
	}

	return app, nil
}

func (s *GitAppsService) ListGithubApps(ctx context.Context) ([]GithubApp, error) {
	return s.repo.ListGithubApps(ctx)
}

func (s *GitAppsService) GetGithubApp(ctx context.Context, id string) (*GithubApp, error) {
	if id == "" {
		return nil, errors.New("app ID is required")
	}
	return s.repo.GetGithubApp(ctx, id)
}

func (s *GitAppsService) SaveGithubApp(ctx context.Context, app *GithubApp) error {
	if app == nil {
		return errors.New("app config is required")
	}
	if app.ID == "" {
		app.ID = uuid.NewString()
	}
	return s.repo.SaveGithubApp(ctx, app)
}

func (s *GitAppsService) DeleteGithubApp(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("app ID is required")
	}
	return s.repo.DeleteGithubApp(ctx, id)
}
