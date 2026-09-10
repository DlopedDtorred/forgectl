package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/forgectl/forgectl/internal/git"
	"github.com/forgectl/forgectl/pkg/ui"
	"github.com/spf13/cobra"
)

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Release management commands",
	Long:  "Create and manage releases for your project.",
}

var releaseNewCmd = &cobra.Command{
	Use:   "new [version]",
	Short: "Create a new release with tag",
	Long: `Create a new release by tagging the current commit and optionally pushing.
If run without arguments, enters interactive mode.

Examples:
  forgectl release new              # Interactive
  forgectl release new v1.0.0       # Direct
  forgectl release new v1.0.0 --push`,
	Args: cobra.MaximumNArgs(1),
	RunE: runReleaseNew,
}

func runReleaseNew(cmd *cobra.Command, args []string) error {
	checkGitInstalled()

	version := ""
	message := ""
	pushRelease := false

	if len(args) == 0 && ui.IsTerminal() {
		ui.Header("Create Release")

		version = ui.PromptTextRequired("Version (e.g., v1.0.0)")
		message = ui.PromptText("Release message", "Release "+version)
		pushRelease = ui.PromptConfirm("Push to remote?", true)
	} else if len(args) > 0 {
		version = args[0]
		message, _ = cmd.Flags().GetString("message")
		pushRelease, _ = cmd.Flags().GetBool("push")
		if message == "" {
			message = "Release " + version
		}
	} else {
		return fmt.Errorf("version is required")
	}

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

	ui.Step("Checking working directory...")
	clean, err := client.IsClean()
	if err != nil {
		return fmt.Errorf("cannot check status: %w", err)
	}
	if !clean {
		ui.Warn("Working directory is not clean")
		if !ui.PromptConfirm("Continue anyway?", false) {
			return nil
		}
	}

	ui.Stepf("Creating tag %s...", version)
	if err := client.Tag(version, message); err != nil {
		return fmt.Errorf("cannot create tag: %w", err)
	}
	ui.Successf("Tag %s created", version)

	// Offer to generate changelog
	genChangelog, _ := cmd.Flags().GetBool("changelog")
	if !genChangelog && ui.IsTerminal() && len(args) == 0 {
		genChangelog = ui.PromptConfirm("Generate CHANGELOG.md?", true)
	}
	if genChangelog {
		ui.Step("Generating CHANGELOG.md...")
		if err := generateChangelogFile(repoRoot, version); err != nil {
			ui.Warnf("Could not generate changelog: %v", err)
		} else {
			ui.Success("CHANGELOG.md updated")
			ui.Dim("Remember to commit the changelog before pushing")
		}
	}

	if pushRelease {
		ui.Step("Pushing tags to remote...")
		if err := client.PushTags("origin"); err != nil {
			return fmt.Errorf("cannot push tags: %w", err)
		}
		ui.Success("Tags pushed to origin")
	}

	ui.Successf("Release %s created successfully", version)
	return nil
}

var releaseListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all releases (tags)",
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

		tags, err := client.Tags()
		if err != nil {
			return fmt.Errorf("cannot list tags: %w", err)
		}

		if len(tags) == 0 {
			ui.Info("No releases found. Create one with: forgectl release new <version>")
			return nil
		}

		ui.Header("Releases")
		for _, tag := range tags {
			ui.Step(tag)
		}
		return nil
	},
}

func generateChangelogFile(repoRoot, releaseVersion string) error {
	client, err := git.Open(repoRoot)
	if err != nil {
		return err
	}

	tags, err := client.Tags()
	if err != nil {
		return err
	}

	sort.Sort(sort.Reverse(sort.StringSlice(tags)))

	var changelog strings.Builder
	changelog.WriteString("# Changelog\n\n")
	changelog.WriteString("All notable changes to this project will be documented in this file.\n\n")
	changelog.WriteString("The format is based on [Keep a Changelog](https://keepachangelog.com/).\n\n")

	changelog.WriteString(fmt.Sprintf("## [%s] - %s\n\n", releaseVersion, time.Now().Format("2006-01-02")))

	var fromTag string
	for i, t := range tags {
		if t == releaseVersion && i+1 < len(tags) {
			fromTag = tags[i+1]
			break
		}
	}

	logs, err := client.GenerateChangelog(fromTag, releaseVersion)
	if err == nil && logs != "" {
		changelog.WriteString(logs + "\n\n")
	} else {
		changelog.WriteString("- Initial release\n\n")
	}

	outputPath := filepath.Join(repoRoot, "CHANGELOG.md")
	return os.WriteFile(outputPath, []byte(changelog.String()), 0o644)
}

func init() {
	releaseCmd.AddCommand(releaseNewCmd)
	releaseCmd.AddCommand(releaseListCmd)
	releaseNewCmd.Flags().StringP("message", "m", "", "release message")
	releaseNewCmd.Flags().Bool("push", false, "push tags to remote after creating")
	releaseNewCmd.Flags().Bool("changelog", false, "generate CHANGELOG.md")
}
