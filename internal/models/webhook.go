package models

import "time"

type GithubWebhookPayload struct {
	Action      string                `json:"action"`
	Number      int                   `json:"number"`
	PullRequest GithubPullRequest     `json:"pull_request"`
	Ref         string                `json:"ref"`
	After       string                `json:"after"`
	Deleted     bool                  `json:"deleted"`
	Repository  GithubRepoInfo        `json:"repository"`
	HeadCommit  *GithubCommitInfo     `json:"head_commit"`
	Commits     []GithubWebhookCommit `json:"commits"`
}

type GithubRepoInfo struct {
	Name          string          `json:"name"`
	FullName      string          `json:"full_name"`
	DefaultBranch string          `json:"default_branch"`
	Owner         GithubRepoOwner `json:"owner"`
}

type GithubRepoOwner struct {
	Login string `json:"login"`
}

type GithubCommitInfo struct {
	ID       string   `json:"id"`
	Message  string   `json:"message"`
	Added    []string `json:"added"`
	Modified []string `json:"modified"`
	Removed  []string `json:"removed"`
}

type GithubWebhookCommit struct {
	ID       string   `json:"id"`
	Message  string   `json:"message"`
	Added    []string `json:"added"`
	Modified []string `json:"modified"`
	Removed  []string `json:"removed"`
}

type GithubPullRequest struct {
	State string                `json:"state"`
	Draft bool                  `json:"draft"`
	Head  GithubPullRequestHead `json:"head"`
	Base  GithubPullRequestBase `json:"base"`
}

type GithubPullRequestHead struct {
	Ref  string                `json:"ref"`
	Sha  string                `json:"sha"`
	Repo GithubPullRequestRepo `json:"repo"`
}

type GithubPullRequestBase struct {
	Ref  string                `json:"ref"`
	Repo GithubPullRequestRepo `json:"repo"`
}

type GithubPullRequestRepo struct {
	FullName string `json:"full_name"`
}

type Webhook struct {
	ID                    string    `json:"id" db:"id"`
	ServiceID             string    `json:"serviceId" db:"service_id"`
	URL                   string    `json:"url" db:"url"`
	EventTypes            []string  `json:"eventTypes" db:"-"`
	IncludePREnvironments bool      `json:"includePrEnvironments" db:"include_pr_environments"`
	CreatedAt             time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt             time.Time `json:"updatedAt" db:"updated_at"`
}

type CreateWebhookRequest struct {
	URL                   string   `json:"url"`
	EventTypes            []string `json:"eventTypes"`
	IncludePREnvironments bool     `json:"includePrEnvironments"`
}
