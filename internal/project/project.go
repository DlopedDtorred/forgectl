package project

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/forgectl/forgectl/internal/config"
	"github.com/forgectl/forgectl/internal/git"
	"github.com/forgectl/forgectl/internal/templates"
	"github.com/forgectl/forgectl/pkg/ui"
)

type Manager struct {
	config *config.Config
}

func NewManager(cfg *config.Config) *Manager {
	return &Manager{config: cfg}
}

type NewProjectOptions struct {
	Name        string
	Description string
	License     string
	Template    string
	Provider    string
	Remote      string
	LocalPath   string
	Private     bool
	NoRemote    bool
	NoCI        bool
	CI          string
	Hooks       bool
	Directories []string
}

func (m *Manager) Create(opts NewProjectOptions) error {
	cfg, err := m.config.GetAppConfig()
	if err != nil {
		return err
	}

	if opts.License == "" {
		opts.License = cfg.DefaultLicense
	}
	if opts.Template == "" {
		opts.Template = cfg.DefaultTemplate
	}
	if opts.Provider == "" {
		opts.Provider = cfg.DefaultProvider
	}

	projectPath := opts.LocalPath
	if projectPath == "" {
		projectDir, _ := os.Getwd()
		projectPath = filepath.Join(projectDir, opts.Name)
	}

	// Validate
	if opts.Name == "" {
		return fmt.Errorf("project name is required")
	}

	ui.Header("Creating Project: " + opts.Name)

	// Step 1: Create directory
	ui.Step("Creating project directory...")
	if err := os.MkdirAll(projectPath, 0o755); err != nil {
		return fmt.Errorf("cannot create project directory: %w", err)
	}
	ui.Successf("Directory: %s", projectPath)

	// Step 2: Init git
	ui.Step("Initializing Git repository...")
	repo, err := git.Init(projectPath)
	if err != nil {
		return err
	}
	ui.Success("Git repository initialized")

	// Step 3: Create project structure
	ui.Step("Creating project structure...")
	structure := templates.DefaultStructure
	if opts.Template != "" {
		if tmpl, ok := templates.ScaffoldTemplates[opts.Template]; ok {
			structure = tmpl.Directories
		}
	}
	for _, dir := range opts.Directories {
		structure = append(structure, dir)
	}

	createdDirs := 0
	for _, dir := range structure {
		dirPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(dirPath, 0o755); err != nil {
			return fmt.Errorf("cannot create directory %s: %w", dir, err)
		}
		keepFile := filepath.Join(dirPath, ".gitkeep")
		os.WriteFile(keepFile, []byte(""), 0o644)
		createdDirs++
	}
	ui.Successf("Created %d directories (%s template)", createdDirs, opts.Template)

	// Step 4: License
	if opts.License != "" && opts.License != "none" {
		ui.Stepf("Adding %s license...", opts.License)
		licenseContent, err := templates.GetLicense(opts.License, time.Now().Year())
		if err != nil {
			ui.Warnf("Could not generate license: %v", err)
		} else {
			licensePath := filepath.Join(projectPath, "LICENSE")
			if err := os.WriteFile(licensePath, []byte(licenseContent), 0o644); err != nil {
				return fmt.Errorf("cannot write license: %w", err)
			}
			ui.Successf("License: %s", opts.License)
		}
	}

	// Step 5: README
	ui.Step("Generating README.md...")
	tmplData := templates.ReadmeData{
		Name:        opts.Name,
		Description: opts.Description,
		License:     opts.License,
		Year:        time.Now().Year(),
	}
	readmeContent, err := templates.GenerateReadme(tmplData)
	if err != nil {
		ui.Warnf("Could not generate README: %v", err)
	} else {
		readmePath := filepath.Join(projectPath, "README.md")
		if err := os.WriteFile(readmePath, []byte(readmeContent), 0o644); err != nil {
			return fmt.Errorf("cannot write README: %w", err)
		}
		ui.Success("README.md generated")
	}

	// Step 6: .gitignore
	ui.Step("Generating .gitignore...")
	gitignoreContent := templates.GetGitignore(opts.Template)
	gitignorePath := filepath.Join(projectPath, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0o644); err != nil {
		return fmt.Errorf("cannot write .gitignore: %w", err)
	}
	ui.Success(".gitignore generated")

	// Step 7: CI/CD
	if opts.CI != "" && !opts.NoCI {
		ui.Stepf("Generating %s CI/CD configuration...", opts.CI)
		ciContent := templates.GetCITemplate(opts.CI, opts.Template)
		if ciContent != "" {
			ciDir := ""
			ciFile := ""
			switch opts.CI {
			case "github-actions":
				ciDir = filepath.Join(projectPath, ".github", "workflows")
				ciFile = "ci.yml"
			case "gitlab-ci":
				ciDir = projectPath
				ciFile = ".gitlab-ci.yml"
			}
			if ciDir != "" && ciFile != "" {
				os.MkdirAll(ciDir, 0o755)
				ciPath := filepath.Join(ciDir, ciFile)
				if err := os.WriteFile(ciPath, []byte(ciContent), 0o644); err != nil {
					ui.Warnf("Could not write CI config: %v", err)
				} else {
					ui.Successf("CI/CD: %s", opts.CI)
				}
			}
		}
	}

	// Step 8: Git hooks
	if opts.Hooks {
		ui.Step("Generating Git hooks...")
		hooksDir := filepath.Join(projectPath, ".git", "hooks")
		if err := os.MkdirAll(hooksDir, 0o755); err == nil {
			hookContent := templates.GetHook("pre-commit")
			hookPath := filepath.Join(hooksDir, "pre-commit")
			if err := os.WriteFile(hookPath, []byte(hookContent), 0o755); err == nil {
				ui.Success("Git hooks generated")
			}
		}
	}

	// Step 9: Remote repository
	var remoteURL string
	if !opts.NoRemote {
		ui.Stepf("Creating remote repository on %s...", opts.Provider)
		provider, err := m.getProvider(opts.Provider, opts.Remote)
		if err != nil {
			ui.Warnf("Cannot create remote repository: %v", err)
			ui.Dim("Run 'forgectl config set' to configure a remote")
		} else {
			repoObj, err := provider.CreateRepository(opts.Name, opts.Description, opts.Private)
			if err != nil {
				ui.Warnf("Cannot create remote repository: %v", err)
			} else {
				remoteURL = repoObj.CloneURL
				remoteName := "origin"
				if opts.Remote != "" {
					remoteName = opts.Remote
				}
				repo.RemoteAddOrSet(remoteName, remoteURL)
				ui.Successf("Remote: %s", repoObj.URL)
			}
		}
	}

	// Step 10: Initial commit
	ui.Step("Creating initial commit...")
	if err := repo.AddAndCommit("Initial commit: project scaffolded by forgectl"); err != nil {
		return fmt.Errorf("initial commit failed: %w", err)
	}
	ui.Success("Initial commit created")

	// Save project config
	projectCfg := config.ProjectConfig{
		Name:        opts.Name,
		Description: opts.Description,
		License:     opts.License,
		Template:    opts.Template,
		Provider:    opts.Provider,
		LocalPath:   projectPath,
	}
	m.config.AddProject(opts.Name, projectCfg)

	// Step 11: Push
	if remoteURL != "" {
		ui.Step("Pushing to remote...")
		branch, _ := repo.CurrentBranch()
		if branch == "" {
			branch = "main"
		}
		remoteName := "origin"
		if opts.Remote != "" {
			remoteName = opts.Remote
		}
		if err := repo.Push(remoteName, branch); err != nil {
			ui.Warnf("Push failed: %v", err)
		} else {
			ui.Success("Pushed to remote")
		}
	}

	// Final summary
	ui.HeaderSmall("Project Created Successfully!")
	ui.Box(opts.Name, []string{
		"Path:     " + projectPath,
		"Template: " + opts.Template,
		"License:  " + opts.License,
		"Remote:   " + remoteURL,
	})

	fmt.Println()
	ui.Stepf("Get started:")
	ui.Dim(fmt.Sprintf("  cd %s", projectPath))
	if opts.Template == "go" {
		ui.Dim("  go mod init " + opts.Name)
		ui.Dim("  go run .")
	} else if opts.Template == "python" {
		ui.Dim("  python -m venv venv")
		ui.Dim("  source venv/bin/activate")
		ui.Dim("  pip install -r requirements.txt")
	} else if opts.Template == "node" {
		ui.Dim("  npm init -y")
		ui.Dim("  npm install")
	} else if opts.Template == "rust" {
		ui.Dim("  cargo build")
		ui.Dim("  cargo run")
	}
	fmt.Println()
	return nil
}

func (m *Manager) getProvider(providerName, remoteName string) (Provider, error) {
	if remoteName != "" {
		remote, err := m.config.GetRemote(remoteName)
		if err != nil {
			return nil, err
		}
		return NewProvider(RemoteConfig{
			Provider: remote.Provider,
			Token:    remote.Token,
			Username: remote.Username,
			BaseURL:  remote.BaseURL,
		})
	}

	remote, err := m.config.GetDefaultRemote()
	if err != nil {
		return nil, err
	}
	return NewProvider(RemoteConfig{
		Provider: remote.Provider,
		Token:    remote.Token,
		Username: remote.Username,
		BaseURL:  remote.BaseURL,
	})
}
