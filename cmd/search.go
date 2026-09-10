package cmd

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/forgectl/forgectl/internal/config"
	"github.com/forgectl/forgectl/internal/project"
	"github.com/forgectl/forgectl/pkg/ui"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for repositories on the configured provider",
	Long: `Search for repositories on the configured remote provider
(GitHub, GitLab, or Gitea).

Examples:
  forgectl search
  forgectl search my-awesome-project
  forgectl search --provider github --query terraform`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.New()
		if err != nil {
			return err
		}

		providerName, _ := cmd.Flags().GetString("provider")
		query := ""
		if len(args) > 0 {
			query = args[0]
		}

		if query == "" && ui.IsTerminal() {
			query = ui.PromptTextRequired("Search query")
		}

		if query == "" {
			return fmt.Errorf("search query is required")
		}

		// Get the remote to use
		var remoteCfg *config.RemoteConfig
		if providerName != "" {
			cfgApp, _ := cfg.GetAppConfig()
			for name, r := range cfgApp.Remotes {
				if r.Provider == providerName {
					rc := r
					remoteCfg = &rc
					_ = name
					break
				}
			}
			if remoteCfg == nil {
				return fmt.Errorf("no remote found for provider %q", providerName)
			}
		} else {
			rc, err := cfg.GetDefaultRemote()
			if err != nil {
				return fmt.Errorf("no remote configured: %v", err)
			}
			remoteCfg = rc
		}

		// Create provider and list repos
		provider, err := project.NewProvider(project.RemoteConfig{
			Provider: remoteCfg.Provider,
			Token:    remoteCfg.Token,
			Username: remoteCfg.Username,
			BaseURL:  remoteCfg.BaseURL,
		})
		if err != nil {
			return err
		}

		ui.Stepf("Searching %s for %q...", remoteCfg.Provider, query)
		repos, err := provider.ListRepositories()
		if err != nil {
			return fmt.Errorf("cannot list repositories: %w", err)
		}

		queryLower := strings.ToLower(query)
		var results []*project.Repository
		for _, r := range repos {
			name := strings.ToLower(r.Name)
			desc := strings.ToLower(r.Description)
			if strings.Contains(name, queryLower) || strings.Contains(desc, queryLower) {
				results = append(results, r)
			}
		}

		if len(results) == 0 {
			ui.Info("No repositories found matching your query")
			return nil
		}

		ui.Successf("Found %d repositories:", len(results))
		fmt.Println()
		for _, r := range results {
			ui.Box(r.Name, []string{
				"Description: " + r.Description,
				"URL:         " + r.URL,
				"Visibility:  " + visibilityString(r.Private),
			})
			fmt.Println()
		}

		return nil
	},
}

func visibilityString(private bool) string {
	if private {
		return "private"
	}
	return "public"
}

func normalizeURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)
	return u.String()
}

func init() {
	searchCmd.Flags().String("provider", "", "provider to search on")
}