package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/forgectl/forgectl/internal/templates"
	"github.com/forgectl/forgectl/pkg/ui"
	"github.com/spf13/cobra"
)

var readmeCmd = &cobra.Command{
	Use:   "readme",
	Short: "README management commands",
	Long:  "Generate, update and manage README files for your project.",
}

var readmeGenCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate a professional README.md",
	Long: `Generate a professional README.md file in the current directory.
If run without arguments, enters interactive mode.

Examples:
  forgectl readme gen              # Interactive
  forgectl readme gen -n my-project`,
	Args: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		license, _ := cmd.Flags().GetString("license")
		author, _ := cmd.Flags().GetString("author")

		if ui.IsTerminal() && name == "" {
			ui.Header("Generate README")
			name = ui.PromptTextRequired("Project name")
			description = ui.PromptText("Description", "")
			license = ui.PromptText("License", "MIT")
			author = ui.PromptText("Author name", "")
		}

		if name == "" {
			dir, _ := os.Getwd()
			name = filepath.Base(dir)
		}

		data := templates.ReadmeData{
			Name:        name,
			Description: description,
			License:     license,
			Year:        time.Now().Year(),
			Author:      author,
		}

		ui.Step("Generating README.md...")
		content, err := templates.GenerateReadme(data)
		if err != nil {
			return fmt.Errorf("cannot generate README: %w", err)
		}

		if err := os.WriteFile("README.md", []byte(content), 0o644); err != nil {
			return fmt.Errorf("cannot write README.md: %w", err)
		}

		ui.Success("README.md generated successfully")
		return nil
	},
}

func init() {
	readmeCmd.AddCommand(readmeGenCmd)
	readmeGenCmd.Flags().StringP("name", "n", "", "project name")
	readmeGenCmd.Flags().StringP("description", "d", "", "project description")
	readmeGenCmd.Flags().StringP("license", "l", "", "license type")
	readmeGenCmd.Flags().String("author", "", "author name")
}
