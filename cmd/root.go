package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/forgectl/forgectl/pkg/ui"
	"github.com/spf13/cobra"
)

var version = "0.2.0"

var rootCmd = &cobra.Command{
	Use:   "forgectl",
	Short: "Professional project management CLI",
	Long: `forgectl is a universal project management tool that automates
the creation, initialization, and management of development projects.

It supports multiple remote providers (GitHub, GitLab, Gitea) and
automates Git workflows, project scaffolding, and more.

Quick start:
  forgectl config set              # Configure provider & tokens
  forgectl new                     # Create a new project (interactive)
  forgectl sync                    # Sync local with remote
  forgectl doctor                  # Diagnose your setup`,
	Version: version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(pushCmd)
	rootCmd.AddCommand(pullCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(remoteCmd)
	rootCmd.AddCommand(readmeCmd)
	rootCmd.AddCommand(changelogCmd)
	rootCmd.AddCommand(licenseCmd)
	rootCmd.AddCommand(releaseCmd)
	rootCmd.AddCommand(issueCmd)
	rootCmd.AddCommand(tagCmd)
	rootCmd.AddCommand(templateCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(projectCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(openCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(versionCmd)

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().Bool("json", false, "output in JSON format")
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print forgectl version",
	Run: func(cmd *cobra.Command, args []string) {
		ui.Header("forgectl")
		ui.Box("Version", []string{
			version,
			"Universal Project Management CLI",
			"https://forgectl.dev",
		})
	},
}

func checkGitInstalled() {
	if err := ensureGit(); err != nil {
		ui.Error(err.Error())
		os.Exit(1)
	}
}

func ensureGit() error {
	_, err := exec.LookPath("git")
	if err != nil {
		return fmt.Errorf("git is not installed or not found in PATH")
	}
	return nil
}

func isJSONOutput(cmd *cobra.Command) bool {
	val, _ := cmd.Flags().GetBool("json")
	return val
}
