package cmd

import (
	"fmt"
	"os"

	"github.com/forgectl/forgectl/internal/git"
	"github.com/forgectl/forgectl/pkg/ui"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronize local repository with remote",
	Long: `Synchronize the local repository with the remote origin.
This performs a git pull followed by a git push.

Examples:
  forgectl sync
  forgectl sync --force`,
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

		branch, err := client.CurrentBranch()
		if err != nil {
			return fmt.Errorf("cannot determine current branch: %w", err)
		}

		remote := "origin"

		ui.Stepf("Syncing with %s/%s...", remote, branch)

		ui.Step("Pulling changes...")
		if err := client.Pull(remote, branch); err != nil {
			ui.Warnf("Pull failed: %v", err)
		} else {
			ui.Success("Pulled from remote")
		}

		ui.Step("Pushing changes...")
		if err := client.Push(remote, branch); err != nil {
			ui.Warnf("Push failed: %v", err)
		} else {
			ui.Success("Pushed to remote")
		}

		ui.Success("Sync completed")
		return nil
	},
}
