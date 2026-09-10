package cmd

import (
	"fmt"
	"os"

	"github.com/forgectl/forgectl/internal/git"
	"github.com/forgectl/forgectl/pkg/ui"
	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:   "tag [tag-name]",
	Short: "Create a Git tag",
	Long: `Create a lightweight or annotated Git tag.
If run without arguments, enters interactive mode.

Examples:
  forgectl tag              # Interactive
  forgectl tag v1.0.0       # Direct
  forgectl tag --list       # List all tags`,
	Args: cobra.MaximumNArgs(1),
	RunE: runTag,
}

func runTag(cmd *cobra.Command, args []string) error {
	checkGitInstalled()

	listTags, _ := cmd.Flags().GetBool("list")

	if listTags {
		return runTagList()
	}

	if len(args) == 0 && ui.IsTerminal() {
		ui.Header("Create Tag")

		tagName := ui.PromptTextRequired("Tag name (e.g., v1.0.0)")
		message := ui.PromptText("Tag message", tagName)

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

		ui.Stepf("Creating tag %s...", tagName)
		if err := client.Tag(tagName, message); err != nil {
			return fmt.Errorf("cannot create tag: %w", err)
		}
		ui.Successf("Tag %s created", tagName)
		return nil
	}

	if len(args) == 0 {
		return fmt.Errorf("tag name is required (or use --list)")
	}

	tagName := args[0]
	message, _ := cmd.Flags().GetString("message")

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

	if message != "" {
		if err := client.Tag(tagName, message); err != nil {
			return fmt.Errorf("cannot create tag: %w", err)
		}
	} else {
		if err := client.Tag(tagName, tagName); err != nil {
			return fmt.Errorf("cannot create tag: %w", err)
		}
	}

	ui.Successf("Tag %s created", tagName)
	return nil
}

func runTagList() error {
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
		return err
	}

	if len(tags) == 0 {
		ui.Info("No tags found. Create one with: forgectl tag <name>")
		return nil
	}

	ui.Header("Tags")
	for _, tag := range tags {
		ui.Step(tag)
	}
	return nil
}

func init() {
	tagCmd.Flags().BoolP("list", "l", false, "list all tags")
	tagCmd.Flags().StringP("message", "m", "", "tag message")
}
