package config

import (
	"os"
	"path/filepath"
	"testing"
)

func newTestConfig(t *testing.T) *Config {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	cfg, err := New()
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}
	return cfg
}

func TestNewCreatesConfigDirAndFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	cfg, err := New()
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}
	if cfg == nil {
		t.Fatal("New() returned nil config")
	}

	configDir := filepath.Join(home, ConfigDirName)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Error("Config directory was not created")
	}

	configFile := filepath.Join(configDir, ConfigFileName)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}
}

func TestNewCreatesDefaultConfig(t *testing.T) {
	cfg := newTestConfig(t)

	appCfg, err := cfg.GetAppConfig()
	if err != nil {
		t.Fatalf("GetAppConfig() returned error: %v", err)
	}

	if appCfg.DefaultProvider != "github" {
		t.Errorf("DefaultProvider = %q, want %q", appCfg.DefaultProvider, "github")
	}
	if appCfg.DefaultLicense != "MIT" {
		t.Errorf("DefaultLicense = %q, want %q", appCfg.DefaultLicense, "MIT")
	}
	if appCfg.DefaultTemplate != "go" {
		t.Errorf("DefaultTemplate = %q, want %q", appCfg.DefaultTemplate, "go")
	}
}

func TestGetAppConfig(t *testing.T) {
	cfg := newTestConfig(t)

	appCfg, err := cfg.GetAppConfig()
	if err != nil {
		t.Fatalf("GetAppConfig() returned error: %v", err)
	}
	if appCfg == nil {
		t.Fatal("GetAppConfig() returned nil")
	}
	if appCfg.Remotes == nil {
		t.Error("Remotes map is nil")
	}
	if appCfg.Projects == nil {
		t.Error("Projects map is nil")
	}
}

func TestSetAndGet(t *testing.T) {
	cfg := newTestConfig(t)

	if err := cfg.Set("default_license", "Apache"); err != nil {
		t.Fatalf("Set() returned error: %v", err)
	}

	val := cfg.Get("default_license")
	if val != "Apache" {
		t.Errorf("Get(default_license) = %v, want Apache", val)
	}

	appCfg, err := cfg.GetAppConfig()
	if err != nil {
		t.Fatalf("GetAppConfig() returned error: %v", err)
	}
	if appCfg.DefaultLicense != "Apache" {
		t.Errorf("DefaultLicense = %q, want Apache", appCfg.DefaultLicense)
	}
}

func TestSetAndGetNestedKey(t *testing.T) {
	cfg := newTestConfig(t)

	if err := cfg.Set("default_template", "rust"); err != nil {
		t.Fatalf("Set() returned error: %v", err)
	}

	val := cfg.Get("default_template")
	if val != "rust" {
		t.Errorf("Get(default_template) = %v, want rust", val)
	}
}

func TestAddAndGetRemote(t *testing.T) {
	cfg := newTestConfig(t)

	remote := RemoteConfig{
		Provider: "github",
		Token:    "test-token-123",
		Username: "testuser",
		BaseURL:  "https://api.github.com",
		Default:  true,
	}

	if err := cfg.AddRemote("origin", remote); err != nil {
		t.Fatalf("AddRemote() returned error: %v", err)
	}

	got, err := cfg.GetRemote("origin")
	if err != nil {
		t.Fatalf("GetRemote() returned error: %v", err)
	}

	if got.Provider != "github" {
		t.Errorf("Provider = %q, want github", got.Provider)
	}
	if got.Token != "test-token-123" {
		t.Errorf("Token = %q, want test-token-123", got.Token)
	}
	if got.Username != "testuser" {
		t.Errorf("Username = %q, want testuser", got.Username)
	}
	if got.BaseURL != "https://api.github.com" {
		t.Errorf("BaseURL = %q, want https://api.github.com", got.BaseURL)
	}
	if !got.Default {
		t.Error("Default = false, want true")
	}
}

func TestGetRemoteNotFound(t *testing.T) {
	cfg := newTestConfig(t)

	_, err := cfg.GetRemote("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent remote")
	}
}

func TestGetDefaultRemote(t *testing.T) {
	cfg := newTestConfig(t)

	remote := RemoteConfig{
		Provider: "github",
		Token:    "token-a",
		Username: "user-a",
	}
	if err := cfg.AddRemote("origin", remote); err != nil {
		t.Fatalf("AddRemote() returned error: %v", err)
	}

	got, err := cfg.GetDefaultRemote()
	if err != nil {
		t.Fatalf("GetDefaultRemote() returned error: %v", err)
	}
	if got.Provider != "github" {
		t.Errorf("Provider = %q, want github", got.Provider)
	}
	if got.Token != "token-a" {
		t.Errorf("Token = %q, want token-a", got.Token)
	}
}

func TestGetDefaultRemoteNoneConfigured(t *testing.T) {
	cfg := newTestConfig(t)

	_, err := cfg.GetDefaultRemote()
	if err == nil {
		t.Error("Expected error when no remotes configured")
	}
}

func TestAddAndGetProject(t *testing.T) {
	cfg := newTestConfig(t)

	project := ProjectConfig{
		Name:        "my-project",
		Description: "A test project",
		License:     "MIT",
		Template:    "go",
		Provider:    "github",
		LocalPath:   "/tmp/my-project",
	}

	if err := cfg.AddProject("my-project", project); err != nil {
		t.Fatalf("AddProject() returned error: %v", err)
	}

	got, err := cfg.GetProject("my-project")
	if err != nil {
		t.Fatalf("GetProject() returned error: %v", err)
	}

	if got.Name != "my-project" {
		t.Errorf("Name = %q, want my-project", got.Name)
	}
	if got.Description != "A test project" {
		t.Errorf("Description = %q, want 'A test project'", got.Description)
	}
	if got.License != "MIT" {
		t.Errorf("License = %q, want MIT", got.License)
	}
	if got.Template != "go" {
		t.Errorf("Template = %q, want go", got.Template)
	}
	if got.LocalPath != "/tmp/my-project" {
		t.Errorf("LocalPath = %q, want /tmp/my-project", got.LocalPath)
	}
}

func TestGetProjectNotFound(t *testing.T) {
	cfg := newTestConfig(t)

	_, err := cfg.GetProject("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent project")
	}
}

func TestListProjects(t *testing.T) {
	cfg := newTestConfig(t)

	projects := cfg.ListProjects()
	if len(projects) != 0 {
		t.Errorf("Expected empty projects, got %d", len(projects))
	}

	cfg.AddProject("proj1", ProjectConfig{Name: "proj1"})
	cfg.AddProject("proj2", ProjectConfig{Name: "proj2"})

	projects = cfg.ListProjects()
	if len(projects) != 2 {
		t.Errorf("Expected 2 projects, got %d", len(projects))
	}
}

func TestRemoveProject(t *testing.T) {
	cfg := newTestConfig(t)

	cfg.AddProject("proj1", ProjectConfig{Name: "proj1"})
	cfg.AddProject("proj2", ProjectConfig{Name: "proj2"})

	if err := cfg.RemoveProject("proj1"); err != nil {
		t.Fatalf("RemoveProject() returned error: %v", err)
	}

	projects := cfg.ListProjects()
	if len(projects) != 1 {
		t.Errorf("Expected 1 project after removal, got %d", len(projects))
	}
	if _, ok := projects["proj2"]; !ok {
		t.Error("proj2 should still exist")
	}
	if _, ok := projects["proj1"]; ok {
		t.Error("proj1 should have been removed")
	}
}

func TestReset(t *testing.T) {
	cfg := newTestConfig(t)

	cfg.Set("default_license", "GPL")
	cfg.Set("default_template", "rust")
	cfg.AddRemote("origin", RemoteConfig{Provider: "github"})
	cfg.AddProject("test", ProjectConfig{Name: "test"})

	if err := cfg.Reset(); err != nil {
		t.Fatalf("Reset() returned error: %v", err)
	}

	appCfg, err := cfg.GetAppConfig()
	if err != nil {
		t.Fatalf("GetAppConfig() returned error: %v", err)
	}

	if appCfg.DefaultLicense != "MIT" {
		t.Errorf("DefaultLicense = %q, want MIT after reset", appCfg.DefaultLicense)
	}
	if appCfg.DefaultTemplate != "go" {
		t.Errorf("DefaultTemplate = %q, go after reset", appCfg.DefaultTemplate)
	}
	if appCfg.DefaultProvider != "github" {
		t.Errorf("DefaultProvider = %q, want github after reset", appCfg.DefaultProvider)
	}
}

func TestConfigDir(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	dir, err := ConfigDir()
	if err != nil {
		t.Fatalf("ConfigDir() returned error: %v", err)
	}

	expected := filepath.Join(home, ConfigDirName)
	if dir != expected {
		t.Errorf("ConfigDir() = %q, want %q", dir, expected)
	}
}
