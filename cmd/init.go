package cmd

import (
	"fmt"
	"os"

	"github.com/forgectl/forgectl/internal/git"
	"github.com/forgectl/forgectl/pkg/ui"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a Git repository in the current directory",
	Long: `Initialize a new Git repository in the current directory.
If a repository already exists, this command will do nothing.

Examples:
  forgectl init
  forgectl init --branch main`,
	RunE: func(cmd *cobra.Command, args []string) error {
		checkGitInstalled()

		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		branch, _ := cmd.Flags().GetString("branch")

		ui.Step("Initializing Git repository...")
		client, err := git.OpenOrInit(dir)
		if err != nil {
			return err
		}

		if branch != "" {
			if err := client.CreateOrphanBranch(branch); err != nil {
				return fmt.Errorf("cannot create branch: %w", err)
			}
			ui.Successf("Created and switched to branch: %s", branch)
		} else {
			ui.Success("Git repository initialized")
		}

		ui.Dim("Repository initialized in: " + dir)
		return nil
	},
}

func init() {
	initCmd.Flags().StringP("branch", "b", "main", "initial branch name")
}
