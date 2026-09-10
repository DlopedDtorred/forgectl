package cmd

import (
	"fmt"
	"os"

	"github.com/forgectl/forgectl/internal/git"
	"github.com/forgectl/forgectl/pkg/ui"
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull [remote] [branch]",
	Short: "Pull changes from remote repository",
	Long: `Pull the latest changes from the remote repository.
Defaults to origin and current branch.

Examples:
  forgectl pull
  forgectl pull origin main`,
	Args: cobra.MaximumNArgs(2),
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

		remote := "origin"
		branch, _ := client.CurrentBranch()

		if len(args) > 0 {
			remote = args[0]
		}
		if len(args) > 1 {
			branch = args[1]
		}

		ui.Stepf("Pulling from %s/%s...", remote, branch)
		if err := client.Pull(remote, branch); err != nil {
			return fmt.Errorf("pull failed: %w", err)
		}
		ui.Successf("Pulled from %s/%s", remote, branch)
		return nil
	},
}
