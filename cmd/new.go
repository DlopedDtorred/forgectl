package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/forgectl/forgectl/internal/config"
	"github.com/forgectl/forgectl/internal/project"
	"github.com/forgectl/forgectl/pkg/ui"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a new project with full scaffolding",
	Long: `Create a new project with Git initialization, remote repository,
README, license, project structure, and initial commit.

If run without arguments, enters interactive wizard mode.

Examples:
  forgectl new                          # Interactive wizard
  forgectl new my-app                   # Quick create
  forgectl new my-app -t go -l MIT      # With flags
  forgectl new my-app --no-remote       # Local only`,
	Args: cobra.MaximumNArgs(1),
	RunE: runNew,
}

func runNew(cmd *cobra.Command, args []string) error {
	checkGitInstalled()

	interactive := len(args) == 0 && ui.IsTerminal()

	if interactive {
		return runNewInteractive()
	}

	name := args[0]
	description, _ := cmd.Flags().GetString("description")
	license, _ := cmd.Flags().GetString("license")
	template, _ := cmd.Flags().GetString("template")
	provider, _ := cmd.Flags().GetString("provider")
	remote, _ := cmd.Flags().GetString("remote")
	private, _ := cmd.Flags().GetBool("private")
	noRemote, _ := cmd.Flags().GetBool("no-remote")
	localPath, _ := cmd.Flags().GetString("path")
	noCI, _ := cmd.Flags().GetBool("no-ci")
	ci, _ := cmd.Flags().GetString("ci")
	hooks, _ := cmd.Flags().GetBool("hooks")

	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("cannot load config: %w", err)
	}

	mgr := project.NewManager(cfg)
	opts := project.NewProjectOptions{
		Name:        name,
		Description: description,
		License:     license,
		Template:    template,
		Provider:    provider,
		Remote:      remote,
		LocalPath:   localPath,
		Private:     private,
		NoRemote:    noRemote,
		NoCI:        noCI,
		CI:          ci,
		Hooks:       hooks,
	}

	return mgr.Create(opts)
}

func runNewInteractive() error {
	ui.Header("Create New Project")

	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("cannot load config: %w", err)
	}

	var opts project.NewProjectOptions

	// Step 1: Project name
	fmt.Println(ui.Bold + "  Let's set up your new project!" + ui.Reset)
	fmt.Println()
	opts.Name = ui.PromptTextRequired("Project name")

	// Validate name
	if strings.Contains(opts.Name, " ") || strings.ContainsAny(opts.Name, "!@#$%^&*()+=[]{}|\\:;\"'<>,?/") {
		ui.Warn("Project name contains special characters. Use kebab-case (e.g., my-project)")
		if !ui.PromptConfirm("Continue anyway?", false) {
			return nil
		}
	}

	// Check if directory exists
	projectPath := filepath.Join(func() string {
		dir, _ := os.Getwd()
		return dir
	}(), opts.Name)
	if _, err := os.Stat(projectPath); err == nil {
		ui.Errorf("Directory %q already exists", opts.Name)
		if !ui.PromptConfirm("Overwrite?", false) {
			return nil
		}
	}

	// Step 2: Description
	opts.Description = ui.PromptText("Description", "")

	// Step 3: Template
	templates := []ui.Choice{
		{Label: "Go", Description: "Standard Go project layout", Value: "go"},
		{Label: "Python", Description: "Standard Python project layout", Value: "python"},
		{Label: "Node.js", Description: "Standard Node.js project layout", Value: "node"},
		{Label: "Rust", Description: "Standard Rust project layout", Value: "rust"},
		{Label: "Minimal", Description: "Bare minimum structure", Value: "minimal"},
	}

	if ui.PromptConfirm("Do you want to select a project template?", true) {
		tmpl := ui.PromptSelect("Select a template:", templates)
		opts.Template = tmpl.Value
	} else {
		opts.Template = "minimal"
	}

	// Step 4: License
	licenses := []ui.Choice{
		{Label: "MIT", Description: "Short and permissive", Value: "MIT"},
		{Label: "Apache 2.0", Description: "Permissive with patent grant", Value: "Apache"},
		{Label: "GPL v3", Description: "Copyleft, derivative works must be GPL", Value: "GPL"},
		{Label: "BSD 2-Clause", Description: "Permissive, similar to MIT", Value: "BSD"},
		{Label: "ISC", Description: "Functionally identical to MIT", Value: "ISC"},
		{Label: "MPL 2.0", Description: "Weak copyleft, file-level", Value: "MPL"},
		{Label: "Unlicense", Description: "Public domain dedication", Value: "Unlicense"},
		{Label: "None", Description: "No license", Value: "none"},
	}

	if ui.PromptConfirm("Do you want to add a license?", true) {
		lic := ui.PromptSelect("Select a license:", licenses)
		opts.License = lic.Value
	}

	// Step 5: Remote provider
	cfgApp, _ := cfg.GetAppConfig()
	providerChoices := []ui.Choice{
		{Label: "GitHub", Description: "github.com", Value: "github"},
		{Label: "GitLab", Description: "gitlab.com", Value: "gitlab"},
		{Label: "Gitea", Description: "Self-hosted Gitea", Value: "gitea"},
		{Label: "None", Description: "No remote repository", Value: "none"},
	}

	if hasRemotes(cfgApp) {
		fmt.Println()
		ui.Info("Configured remotes found:")
		for name, r := range cfgApp.Remotes {
			ui.Dim(fmt.Sprintf("  %s (%s) - %s", name, r.Provider, r.Username))
		}
		fmt.Println()
	}

	if ui.PromptConfirm("Do you want to create a remote repository?", true) {
		provider := ui.PromptSelect("Select remote provider:", providerChoices)
		if provider.Value != "none" {
			opts.Provider = provider.Value
			opts.Private = ui.PromptConfirm("Make repository private?", false)
		} else {
			opts.NoRemote = true
		}
	} else {
		opts.NoRemote = true
	}

	// Step 6: CI/CD
	if ui.PromptConfirm("Do you want to generate CI/CD configuration?", true) {
		ciChoices := []ui.Choice{
			{Label: "GitHub Actions", Description: "For GitHub repositories", Value: "github-actions"},
			{Label: "GitLab CI", Description: "For GitLab repositories", Value: "gitlab-ci"},
			{Label: "None", Description: "No CI/CD configuration", Value: "none"},
		}
		ci := ui.PromptSelect("Select CI/CD provider:", ciChoices)
		if ci.Value != "none" {
			opts.CI = ci.Value
		} else {
			opts.NoCI = true
		}
	}

	// Step 7: Git hooks
	if ui.PromptConfirm("Do you want to generate Git hooks?", true) {
		opts.Hooks = true
	}

	// Summary
	ui.HeaderSmall("Project Summary")
	ui.Box(opts.Name, []string{
		"Template:    " + opts.Template,
		"License:     " + opts.License,
		"Description: " + opts.Description,
		"Provider:    " + opts.Provider,
		"Private:     " + fmt.Sprintf("%t", opts.Private),
		"CI/CD:       " + opts.CI,
		"Hooks:       " + fmt.Sprintf("%t", opts.Hooks),
	})

	fmt.Println()
	if !ui.PromptConfirm("Create this project?", true) {
		ui.Warn("Aborted.")
		return nil
	}

	mgr := project.NewManager(cfg)
	return mgr.Create(opts)
}

func hasRemotes(cfg *config.AppConfig) bool {
	if cfg == nil {
		return false
	}
	return len(cfg.Remotes) > 0
}

func init() {
	newCmd.Flags().StringP("description", "d", "", "project description")
	newCmd.Flags().StringP("license", "l", "", "license type")
	newCmd.Flags().StringP("template", "t", "", "project template")
	newCmd.Flags().StringP("provider", "p", "", "remote provider")
	newCmd.Flags().String("remote", "", "remote name from config")
	newCmd.Flags().BoolP("private", "P", false, "create private repository")
	newCmd.Flags().Bool("no-remote", false, "skip remote repository creation")
	newCmd.Flags().String("path", "", "custom local path")
	newCmd.Flags().String("ci", "", "CI/CD provider (github-actions, gitlab-ci)")
	newCmd.Flags().Bool("no-ci", false, "skip CI/CD generation")
	newCmd.Flags().Bool("hooks", false, "generate Git hooks")
}
