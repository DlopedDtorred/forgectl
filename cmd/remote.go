package cmd

import (
	"fmt"

	"github.com/forgectl/forgectl/internal/config"
	"github.com/forgectl/forgectl/pkg/ui"
	"github.com/spf13/cobra"
)

var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Manage remote repositories and providers",
	Long:  "List, add, and manage remote providers (GitHub, GitLab, Gitea).",
}

var remoteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured remotes",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.New()
		if err != nil {
			return err
		}

		appCfg, err := cfg.GetAppConfig()
		if err != nil {
			return err
		}

		if len(appCfg.Remotes) == 0 {
			ui.Warn("No remotes configured. Run: forgectl remote add")
			return nil
		}

		ui.Header("Configured Remotes")
		for name, r := range appCfg.Remotes {
			lines := []string{
				fmt.Sprintf("Provider:  %s", r.Provider),
				fmt.Sprintf("Username:  %s", r.Username),
				fmt.Sprintf("Base URL:  %s", r.BaseURL),
			}
			if r.Token != "" {
				lines = append(lines, fmt.Sprintf("Token:     %s", r.Token[:min(len(r.Token), 6)]+"..."))
			} else {
				lines = append(lines, fmt.Sprintf("Token:     %s", ui.Yellow+"not set"+ui.Reset))
			}
			if r.Default {
				lines = append(lines, ui.Green+"DEFAULT"+ui.Reset)
			}
			ui.Box(name, lines)
			fmt.Println()
		}

		ui.Infof("Default provider: %s", appCfg.DefaultProvider)
		return nil
	},
}

var remoteAddCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add a new remote",
	Long: `Add a new remote provider interactively.

Examples:
  forgectl remote add                # Interactive
  forgectl remote add github
  forgectl remote add gitlab

  # Or set values directly:
  forgectl config set remotes.github.token ghp_xxx
  forgectl config set remotes.github.username myuser`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.New()
		if err != nil {
			return err
		}

		var name string
		if len(args) > 0 {
			name = args[0]
		} else if ui.IsTerminal() {
			name = ui.PromptTextRequired("Remote name (e.g., github, gitlab, gitea)")
		} else {
			return fmt.Errorf("remote name is required")
		}

		ui.Stepf("Configuring remote %q...", name)

		if ui.IsTerminal() {
			return setupRemoteInteractive(cfg, name)
		}

		// Non-interactive: check flags
		provider, _ := cmd.Flags().GetString("provider")
		username, _ := cmd.Flags().GetString("username")
		token, _ := cmd.Flags().GetString("token")
		baseURL, _ := cmd.Flags().GetString("base-url")

		if provider == "" {
			provider = name
		}
		if username == "" || token == "" {
			return fmt.Errorf("username and token are required (or run interactively)")
		}

		cfg.Set(fmt.Sprintf("remotes.%s.provider", name), provider)
		cfg.Set(fmt.Sprintf("remotes.%s.username", name), username)
		cfg.Set(fmt.Sprintf("remotes.%s.token", name), token)
		if baseURL != "" {
			cfg.Set(fmt.Sprintf("remotes.%s.base_url", name), baseURL)
		}
		cfg.Set(fmt.Sprintf("remotes.%s.default", name), true)

		ui.Successf("Remote %q configured", name)
		return nil
	},
}

var remoteRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove a configured remote",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("remote name is required")
		}

		cfg, err := config.New()
		if err != nil {
			return err
		}

		if _, err := cfg.GetRemote(args[0]); err != nil {
			return err
		}

		if !ui.PromptConfirm(fmt.Sprintf("Remove remote %q?", args[0]), false) {
			ui.Info("Aborted")
			return nil
		}

		cfg.Set(fmt.Sprintf("remotes.%s", args[0]), nil)
		ui.Successf("Remote %q removed", args[0])
		return nil
	},
}

func setupRemoteInteractive(cfg *config.Config, name string) error {
	if name == "" {
		name = ui.PromptTextRequired("Remote name (e.g., github, gitlab)")
	}

	ui.HeaderSmall("Remote Provider Setup")

	provider := ui.PromptSelect("Select provider:", []ui.Choice{
		{Label: "GitHub", Description: "github.com", Value: "github"},
		{Label: "GitLab", Description: "gitlab.com", Value: "gitlab"},
		{Label: "Gitea", Description: "Self-hosted Gitea instance", Value: "gitea"},
	})

	ui.Info("Enter your " + provider.Label + " credentials:")
	username := ui.PromptTextRequired("Username")
	token := ui.PromptPassword("API Token")

	baseURL := ""
	if provider.Value == "gitea" {
		baseURL = ui.PromptTextRequired("Gitea base URL (e.g., https://gitea.example.com)")
	}

	isDefault := ui.PromptConfirm("Set as default remote?", true)

	cfg.Set(fmt.Sprintf("remotes.%s.provider", name), provider.Value)
	cfg.Set(fmt.Sprintf("remotes.%s.username", name), username)
	cfg.Set(fmt.Sprintf("remotes.%s.token", name), token)
	cfg.Set(fmt.Sprintf("remotes.%s.default", name), isDefault)
	if baseURL != "" {
		cfg.Set(fmt.Sprintf("remotes.%s.base_url", name), baseURL)
	}

	ui.Successf("Remote %q configured successfully", name)
	return nil
}

func init() {
	remoteCmd.AddCommand(remoteListCmd)
	remoteCmd.AddCommand(remoteAddCmd)
	remoteCmd.AddCommand(remoteRemoveCmd)

	remoteAddCmd.Flags().String("provider", "", "provider type (github, gitlab, gitea)")
	remoteAddCmd.Flags().String("username", "", "username")
	remoteAddCmd.Flags().String("token", "", "API token")
	remoteAddCmd.Flags().String("base-url", "", "base URL (required for gitea)")
}