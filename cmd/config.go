package cmd

import (
	"fmt"

	"github.com/forgectl/forgectl/internal/config"
	"github.com/forgectl/forgectl/pkg/ui"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management",
	Long:  "Manage forgectl configuration including remotes and defaults.",
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Long: `Set a configuration value. Supports dot notation for nested keys.
If run without arguments, enters interactive mode.

Examples:
  forgectl config set                              # Interactive
  forgectl config set default_provider github
  forgectl config set remotes.github.token ghp_xxx`,
	Args: cobra.MaximumNArgs(2),
	RunE: runConfigSet,
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	if len(args) == 0 && ui.IsTerminal() {
		return runConfigSetInteractive()
	}
	if len(args) < 2 {
		return fmt.Errorf("usage: forgectl config set <key> <value>")
	}

	cfg, err := config.New()
	if err != nil {
		return err
	}

	if err := cfg.Set(args[0], args[1]); err != nil {
		return fmt.Errorf("cannot set config: %w", err)
	}

	ui.Successf("Set %s = %s", args[0], args[1])
	return nil
}

func runConfigSetInteractive() error {
	ui.Header("Configuration Setup")

	cfg, err := config.New()
	if err != nil {
		return err
	}

	action := ui.PromptSelect("What do you want to configure?", []ui.Choice{
		{Label: "Remote provider (GitHub, GitLab, Gitea)", Value: "remote"},
		{Label: "Default settings", Value: "defaults"},
		{Label: "Set a custom key", Value: "custom"},
	})

	switch action.Value {
	case "remote":
		return setupRemoteInteractive(cfg, "")
	case "defaults":
		return setupDefaultsInteractive(cfg)
	case "custom":
		return setupCustomInteractive(cfg)
	}
	return nil
}

func setupDefaultsInteractive(cfg *config.Config) error {
	ui.HeaderSmall("Default Settings")

	providerChoices := []ui.Choice{
		{Label: "GitHub", Value: "github"},
		{Label: "GitLab", Value: "gitlab"},
		{Label: "Gitea", Value: "gitea"},
	}
	provider := ui.PromptSelect("Default provider:", providerChoices)
	cfg.Set("default_provider", provider.Value)

	licenseChoices := []ui.Choice{
		{Label: "MIT", Value: "MIT"},
		{Label: "Apache 2.0", Value: "Apache"},
		{Label: "GPL v3", Value: "GPL"},
		{Label: "BSD 2-Clause", Value: "BSD"},
		{Label: "ISC", Value: "ISC"},
		{Label: "None", Value: "none"},
	}
	lic := ui.PromptSelect("Default license:", licenseChoices)
	cfg.Set("default_license", lic.Value)

	templateChoices := []ui.Choice{
		{Label: "Go", Value: "go"},
		{Label: "Python", Value: "python"},
		{Label: "Node.js", Value: "node"},
		{Label: "Rust", Value: "rust"},
		{Label: "Minimal", Value: "minimal"},
	}
	tmpl := ui.PromptSelect("Default template:", templateChoices)
	cfg.Set("default_template", tmpl.Value)

	ui.Success("Default settings updated")
	return nil
}

func setupCustomInteractive(cfg *config.Config) error {
	ui.HeaderSmall("Custom Configuration")

	key := ui.PromptTextRequired("Key (e.g., remotes.github.token)")
	value := ui.PromptTextRequired("Value")

	cfg.Set(key, value)
	ui.Successf("Set %s = %s", key, value)
	return nil
}

var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.New()
		if err != nil {
			return err
		}

		value := cfg.Get(args[0])
		if value == nil {
			return fmt.Errorf("key %q not found", args[0])
		}

		fmt.Printf("%v\n", value)
		return nil
	},
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset configuration to defaults",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !ui.IsTerminal() || ui.PromptConfirm("Reset all configuration to defaults?", false) {
			cfg, err := config.New()
			if err != nil {
				return err
			}
			if err := cfg.Reset(); err != nil {
				return fmt.Errorf("cannot reset config: %w", err)
			}
			ui.Success("Configuration reset to defaults")
		}
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.New()
		if err != nil {
			return err
		}

		appCfg, err := cfg.GetAppConfig()
		if err != nil {
			return err
		}

		ui.Header("Current Configuration")

		ui.Box("Defaults", []string{
			fmt.Sprintf("Provider:  %s", appCfg.DefaultProvider),
			fmt.Sprintf("License:   %s", appCfg.DefaultLicense),
			fmt.Sprintf("Template:  %s", appCfg.DefaultTemplate),
		})

		if len(appCfg.Remotes) > 0 {
			fmt.Println()
			for name, r := range appCfg.Remotes {
				tokenDisplay := "not set"
				if r.Token != "" {
					tokenDisplay = r.Token[:min(8, len(r.Token))] + "..."
				}
				ui.Box("Remote: "+name, []string{
					fmt.Sprintf("Provider:  %s", r.Provider),
					fmt.Sprintf("Username:  %s", r.Username),
					fmt.Sprintf("Token:     %s", tokenDisplay),
					fmt.Sprintf("Base URL:  %s", r.BaseURL),
					fmt.Sprintf("Default:   %t", r.Default),
				})
			}
		} else {
			fmt.Println()
			ui.Warn("No remotes configured. Run: forgectl config set")
		}

		if len(appCfg.Projects) > 0 {
			fmt.Println()
			ui.HeaderSmall("Tracked Projects")
			for name, p := range appCfg.Projects {
				ui.Dim(fmt.Sprintf("  %s: %s (%s)", name, p.Description, p.Template))
			}
		}

		return nil
	},
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show configuration file path",
	Run: func(cmd *cobra.Command, args []string) {
		home, _ := config.ConfigDir()
		fmt.Println(home)
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configResetCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configPathCmd)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
