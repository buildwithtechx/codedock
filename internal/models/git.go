package models

import "time"

type GitProviderConfig struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"userId" db:"user_id"`
	Provider    string    `json:"provider" db:"provider"`
	AccessToken string    `json:"accessToken,omitempty" db:"encrypted_access_token"`
	AccountName string    `json:"accountName" db:"account_name"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type GitRepository struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	FullName      string `json:"fullName"`
	Private       bool   `json:"private"`
	CloneURL      string `json:"cloneUrl"`
	HTMLURL       string `json:"htmlUrl"`
	DefaultBranch string `json:"defaultBranch"`
}

type GitConnectRequest struct {
	Provider    string `json:"provider"`
	AccessToken string `json:"accessToken"`
	AccountName string `json:"accountName"`
}

type GithubApp struct {
	ID             string    `json:"id" db:"id"`
	Name           string    `json:"name" db:"name"`
	AppID          string    `json:"appId" db:"app_id"`
	InstallationID string    `json:"installationId" db:"installation_id"`
	ClientID       string    `json:"clientId" db:"client_id"`
	ClientSecret   string    `json:"clientSecret" db:"client_secret"`
	WebhookSecret  string    `json:"webhookSecret" db:"webhook_secret"`
	PrivateKey     string    `json:"privateKey" db:"private_key"`
	IsPublic       bool      `json:"isPublic" db:"is_public"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}

type GitHubRepoDTO struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	Private       bool   `json:"private"`
	CloneURL      string `json:"clone_url"`
	HTMLURL       string `json:"html_url"`
	DefaultBranch string `json:"default_branch"`
}

type GitLabRepoDTO struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	PathWithNamespace string `json:"path_with_namespace"`
	Visibility        string `json:"visibility"`
	HTTPURL           string `json:"http_url_to_repo"`
	WebURL            string `json:"web_url"`
	DefaultBranch     string `json:"default_branch"`
}

type BitbucketCloneLinkDTO struct {
	Name string `json:"name"`
	Href string `json:"href"`
}

type BitbucketHTMLLinkDTO struct {
	Href string `json:"href"`
}

type BitbucketLinksDTO struct {
	HTML  BitbucketHTMLLinkDTO    `json:"html"`
	Clone []BitbucketCloneLinkDTO `json:"clone"`
}

type BitbucketBranchDTO struct {
	Name string `json:"name"`
}

type BitbucketRepoDTO struct {
	UUID       string             `json:"uuid"`
	Name       string             `json:"name"`
	FullName   string             `json:"full_name"`
	IsPriv     bool               `json:"is_private"`
	Links      BitbucketLinksDTO  `json:"links"`
	Mainbranch BitbucketBranchDTO `json:"mainbranch"`
}

type BitbucketRepoListDTO struct {
	Values []BitbucketRepoDTO `json:"values"`
}

type GiteaRepoDTO struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	Private       bool   `json:"private"`
	CloneURL      string `json:"clone_url"`
	HTMLURL       string `json:"html_url"`
	DefaultBranch string `json:"default_branch"`
}
