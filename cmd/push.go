package cmd

import (
	"fmt"
	"os"

	"github.com/forgectl/forgectl/internal/git"
	"github.com/forgectl/forgectl/pkg/ui"
	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push [remote] [branch]",
	Short: "Push local changes to remote repository",
	Long: `Push local commits to the remote repository.
Defaults to origin and current branch.

Examples:
  forgectl push
  forgectl push origin main
  forgectl push --tags`,
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

		pushTags, _ := cmd.Flags().GetBool("tags")

		ui.Stepf("Pushing to %s/%s...", remote, branch)
		if err := client.Push(remote, branch); err != nil {
			return fmt.Errorf("push failed: %w", err)
		}
		ui.Successf("Pushed to %s/%s", remote, branch)

		if pushTags {
			ui.Step("Pushing tags...")
			if err := client.PushTags(remote); err != nil {
				return fmt.Errorf("push tags failed: %w", err)
			}
			ui.Success("Tags pushed successfully")
		}

		return nil
	},
}

func init() {
	pushCmd.Flags().Bool("tags", false, "push tags along with commits")
}
