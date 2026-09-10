package cmd

import (
	"fmt"
	"os"

	"github.com/forgectl/forgectl/internal/git"
	"github.com/forgectl/forgectl/pkg/ui"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show repository status",
	Long: `Display the current status of the repository including
branch, remote, uncommitted changes, and tags.

Examples:
  forgectl status`,
	RunE: func(cmd *cobra.Command, args []string) error {
		checkGitInstalled()

		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		repoRoot, err := git.FindRepoRoot(dir)
		if err != nil {
			return err
		}

		client, err := git.Open(repoRoot)
		if err != nil {
			return err
		}

		ui.Header("Repository Status")

		// Branch
		branch, err := client.CurrentBranch()
		if err != nil {
			ui.Warnf("Cannot determine branch: %v", err)
		} else {
			ui.Info("Branch: " + branch)
		}

		// Remote
		if url, err := client.GetRemoteURL("origin"); err == nil {
			ui.Info("Remote: " + url)
		} else {
			ui.Warn("No origin remote configured")
		}

		// Clean?
		clean, err := client.IsClean()
		if err != nil {
			ui.Warnf("Cannot check status: %v", err)
		} else if clean {
			ui.Success("Working tree: clean")
		} else {
			ui.Warn("Working tree: has uncommitted changes")
			ui.Dim("Run 'forgectl sync' to push changes, or 'git status' for details")
		}

		// Tags
		tags, _ := client.Tags()
		if len(tags) > 0 {
			ui.Infof("Tags: %d", len(tags))
			fmt.Println()
			for _, tag := range tags {
				ui.Dim("  " + tag)
			}
		}

		// Recent commits
		fmt.Println()
		ui.Info("Recent commits:")
		logs, err := client.Log(5)
		if err == nil && len(logs) > 0 {
			for _, line := range logs {
				ui.Dim("  " + line)
			}
		}

		return nil
	},
}