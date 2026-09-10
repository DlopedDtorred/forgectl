package project

import (
	"testing"

	"github.com/forgectl/forgectl/internal/config"
)

func TestNewManager(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	cfg, err := config.New()
	if err != nil {
		t.Fatalf("config.New() returned error: %v", err)
	}

	m := NewManager(cfg)
	if m == nil {
		t.Fatal("NewManager() returned nil")
	}
	if m.config != cfg {
		t.Error("Manager config does not match provided config")
	}
}

func TestNewProjectOptionsDefaults(t *testing.T) {
	opts := NewProjectOptions{}

	if opts.Name != "" {
		t.Errorf("Name = %q, want empty", opts.Name)
	}
	if opts.License != "" {
		t.Errorf("License = %q, want empty", opts.License)
	}
	if opts.Template != "" {
		t.Errorf("Template = %q, want empty", opts.Template)
	}
	if opts.Provider != "" {
		t.Errorf("Provider = %q, want empty", opts.Provider)
	}
	if opts.Private {
		t.Error("Private = true, want false")
	}
	if opts.NoRemote {
		t.Error("NoRemote = true, want false")
	}
	if opts.NoCI {
		t.Error("NoCI = true, want false")
	}
	if opts.Hooks {
		t.Error("Hooks = true, want false")
	}
	if opts.EditorConfig {
		t.Error("EditorConfig = true, want false")
	}
	if opts.Codeowners {
		t.Error("Codeowners = true, want false")
	}
	if opts.Security {
		t.Error("Security = true, want false")
	}
}

func TestNewProjectOptionsWithValues(t *testing.T) {
	opts := NewProjectOptions{
		Name:        "my-app",
		Description: "A test app",
		License:     "MIT",
		Template:    "go",
		Provider:    "github",
		Private:     true,
		Hooks:       true,
	}

	if opts.Name != "my-app" {
		t.Errorf("Name = %q, want my-app", opts.Name)
	}
	if opts.Description != "A test app" {
		t.Errorf("Description = %q, want 'A test app'", opts.Description)
	}
	if opts.License != "MIT" {
		t.Errorf("License = %q, want MIT", opts.License)
	}
	if opts.Template != "go" {
		t.Errorf("Template = %q, want go", opts.Template)
	}
	if opts.Provider != "github" {
		t.Errorf("Provider = %q, want github", opts.Provider)
	}
	if !opts.Private {
		t.Error("Private = false, want true")
	}
	if !opts.Hooks {
		t.Error("Hooks = false, want true")
	}
}
