package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/forgectl/forgectl/internal/git"
	"github.com/forgectl/forgectl/pkg/ui"
	"github.com/spf13/cobra"
)

var changelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "Changelog management commands",
	Long:  "Generate and manage changelogs for your project.",
}

var changelogGenCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate a CHANGELOG.md from git history",
	Long: `Generate a CHANGELOG.md file based on git commit history and tags.

Examples:
  forgectl changelog gen
  forgectl changelog gen --output CHANGELOG.md`,
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

		output, _ := cmd.Flags().GetString("output")
		if output == "" {
			output = "CHANGELOG.md"
		}

		ui.Step("Analyzing git history...")

		tags, err := client.Tags()
		if err != nil {
			return fmt.Errorf("cannot list tags: %w", err)
		}

		sort.Sort(sort.Reverse(sort.StringSlice(tags)))

		var changelog strings.Builder
		changelog.WriteString("# Changelog\n\n")
		changelog.WriteString("All notable changes to this project will be documented in this file.\n\n")
		changelog.WriteString("The format is based on [Keep a Changelog](https://keepachangelog.com/).\n\n")

		if len(tags) == 0 {
			logs, err := client.Log(50)
			if err != nil {
				return fmt.Errorf("cannot get git log: %w", err)
			}
			if len(logs) > 0 {
				changelog.WriteString("## [Unreleased]\n\n")
				for _, line := range logs {
					changelog.WriteString(line + "\n")
				}
				changelog.WriteString("\n")
			}
		} else {
			for i, tag := range tags {
				changelog.WriteString(fmt.Sprintf("## [%s]\n\n", tag))

				var fromTag string
				if i < len(tags)-1 {
					fromTag = tags[i+1]
				}

				logs, err := client.GenerateChangelog(fromTag, tag)
				if err != nil {
					continue
				}
				if logs != "" {
					changelog.WriteString(logs + "\n\n")
				}
			}

			logs, err := client.GenerateChangelog(tags[0], "HEAD")
			if err == nil && logs != "" {
				changelog.WriteString("## [Unreleased]\n\n")
				changelog.WriteString(logs + "\n\n")
			}
		}

		outputPath := filepath.Join(repoRoot, output)
		ui.Stepf("Writing %s...", output)
		if err := os.WriteFile(outputPath, []byte(changelog.String()), 0o644); err != nil {
			return fmt.Errorf("cannot write changelog: %w", err)
		}

		ui.Successf("Changelog generated: %s", outputPath)
		return nil
	},
}

func init() {
	changelogCmd.AddCommand(changelogGenCmd)
	changelogGenCmd.Flags().StringP("output", "o", "", "output file path (default: CHANGELOG.md)")
}
