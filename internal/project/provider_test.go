package project

import (
	"testing"
	"time"
)

func TestRepositoryStruct(t *testing.T) {
	now := time.Now()
	repo := Repository{
		Name:        "test-repo",
		Description: "A test repository",
		Private:     true,
		URL:         "https://github.com/user/test-repo",
		CloneURL:    "https://github.com/user/test-repo.git",
		CreatedAt:   now,
	}

	if repo.Name != "test-repo" {
		t.Errorf("Name = %q, want test-repo", repo.Name)
	}
	if repo.Description != "A test repository" {
		t.Errorf("Description = %q, want 'A test repository'", repo.Description)
	}
	if !repo.Private {
		t.Error("Private = false, want true")
	}
	if repo.URL != "https://github.com/user/test-repo" {
		t.Errorf("URL = %q", repo.URL)
	}
	if repo.CloneURL != "https://github.com/user/test-repo.git" {
		t.Errorf("CloneURL = %q", repo.CloneURL)
	}
	if !repo.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v, want %v", repo.CreatedAt, now)
	}
}

func TestReleaseStruct(t *testing.T) {
	now := time.Now()
	release := Release{
		Name:       "v1.0.0",
		TagName:    "v1.0.0",
		Body:       "First release",
		Draft:      false,
		Prerelease: true,
		URL:        "https://github.com/user/repo/releases/tag/v1.0.0",
		CreatedAt:  now,
	}

	if release.Name != "v1.0.0" {
		t.Errorf("Name = %q, want v1.0.0", release.Name)
	}
	if release.TagName != "v1.0.0" {
		t.Errorf("TagName = %q, want v1.0.0", release.TagName)
	}
	if release.Body != "First release" {
		t.Errorf("Body = %q, want 'First release'", release.Body)
	}
	if release.Draft {
		t.Error("Draft = true, want false")
	}
	if !release.Prerelease {
		t.Error("Prerelease = false, want true")
	}
}

func TestIssueStruct(t *testing.T) {
	issue := Issue{
		Number: 42,
		Title:  "Fix bug",
		Body:   "Description of the bug",
		URL:    "https://github.com/user/repo/issues/42",
		Labels: []string{"bug", "urgent"},
	}

	if issue.Number != 42 {
		t.Errorf("Number = %d, want 42", issue.Number)
	}
	if issue.Title != "Fix bug" {
		t.Errorf("Title = %q, want 'Fix bug'", issue.Title)
	}
	if issue.Body != "Description of the bug" {
		t.Errorf("Body = %q", issue.Body)
	}
	if len(issue.Labels) != 2 {
		t.Errorf("Labels has %d items, want 2", len(issue.Labels))
	}
	if issue.Labels[0] != "bug" || issue.Labels[1] != "urgent" {
		t.Errorf("Labels = %v, want [bug urgent]", issue.Labels)
	}
}

func TestRemoteConfigStruct(t *testing.T) {
	cfg := RemoteConfig{
		Provider:     "github",
		Token:        "ghp_test123",
		Username:     "testuser",
		BaseURL:      "https://api.github.com",
		Organization: "my-org",
	}

	if cfg.Provider != "github" {
		t.Errorf("Provider = %q, want github", cfg.Provider)
	}
	if cfg.Token != "ghp_test123" {
		t.Errorf("Token = %q, want ghp_test123", cfg.Token)
	}
	if cfg.Username != "testuser" {
		t.Errorf("Username = %q, want testuser", cfg.Username)
	}
	if cfg.BaseURL != "https://api.github.com" {
		t.Errorf("BaseURL = %q", cfg.BaseURL)
	}
	if cfg.Organization != "my-org" {
		t.Errorf("Organization = %q, want my-org", cfg.Organization)
	}
}

func TestNewProviderGitHub(t *testing.T) {
	cfg := RemoteConfig{
		Provider: "github",
		Token:    "test-token",
		Username: "testuser",
	}

	p, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("NewProvider(github) returned error: %v", err)
	}
	if p.Name() != "github" {
		t.Errorf("Name() = %q, want github", p.Name())
	}
}

func TestNewProviderGitHubDefault(t *testing.T) {
	cfg := RemoteConfig{
		Token:    "test-token",
		Username: "testuser",
	}

	p, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("NewProvider(default) returned error: %v", err)
	}
	if p.Name() != "github" {
		t.Errorf("Name() = %q, want github", p.Name())
	}
}

func TestNewProviderGitLab(t *testing.T) {
	cfg := RemoteConfig{
		Provider: "gitlab",
		Token:    "test-token",
		Username: "testuser",
	}

	p, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("NewProvider(gitlab) returned error: %v", err)
	}
	if p.Name() != "gitlab" {
		t.Errorf("Name() = %q, want gitlab", p.Name())
	}
}

func TestNewProviderGiteaRequiresBaseURL(t *testing.T) {
	cfg := RemoteConfig{
		Provider: "gitea",
		Token:    "test-token",
		Username: "testuser",
	}

	_, err := NewProvider(cfg)
	if err == nil {
		t.Error("Expected error for gitea without base_url")
	}
}

func TestNewProviderGitea(t *testing.T) {
	cfg := RemoteConfig{
		Provider: "gitea",
		Token:    "test-token",
		Username: "testuser",
		BaseURL:  "https://gitea.example.com",
	}

	p, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("NewProvider(gitea) returned error: %v", err)
	}
	if p.Name() != "gitea" {
		t.Errorf("Name() = %q, want gitea", p.Name())
	}
}

func TestNewProviderUnsupported(t *testing.T) {
	cfg := RemoteConfig{
		Provider: "bitbucket",
		Token:    "test-token",
	}

	_, err := NewProvider(cfg)
	if err == nil {
		t.Error("Expected error for unsupported provider")
	}
}
