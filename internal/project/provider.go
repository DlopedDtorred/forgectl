package project

import (
	"fmt"
	"time"
)

type Repository struct {
	Name        string
	Description string
	Private     bool
	URL         string
	CloneURL    string
	CreatedAt   time.Time
}

type Release struct {
	Name       string
	TagName    string
	Body       string
	Draft      bool
	Prerelease bool
	URL        string
	CreatedAt  time.Time
}

type Issue struct {
	Number int
	Title  string
	Body   string
	URL    string
	Labels []string
}

type Provider interface {
	Name() string
	CreateRepository(name, description string, private bool) (*Repository, error)
	DeleteRepository(name string) error
	GetRepository(name string) (*Repository, error)
	ListRepositories() ([]*Repository, error)
	GetCurrentUsername() (string, error)
	CreateRelease(repoName, tagName, name, body string, draft, prerelease bool) (*Release, error)
	ListReleases(repoName string) ([]*Release, error)
	CreateIssue(repoName, title, body string, labels []string) (*Issue, error)
	ListIssues(repoName string) ([]*Issue, error)
}

type RemoteConfig struct {
	Provider     string
	Token        string
	Username     string
	BaseURL      string
	Organization string
}

func NewProvider(cfg RemoteConfig) (Provider, error) {
	switch cfg.Provider {
	case "github", "":
		return newGitHubClient(cfg)
	case "gitlab":
		return newGitLabClient(cfg)
	case "gitea":
		return newGiteaClient(cfg)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", cfg.Provider)
	}
}
