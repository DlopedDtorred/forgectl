package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/forgectl/forgectl/internal/templates"
	"github.com/forgectl/forgectl/pkg/ui"
	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Project template management",
	Long:  "Apply and manage project templates.",
}

var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available project templates",
	Run: func(cmd *cobra.Command, args []string) {
		ui.Header("Available Templates")
		for name, tmpl := range templates.ScaffoldTemplates {
			ui.Box(name, []string{
				tmpl.Description,
				fmt.Sprintf("Directories: %d", len(tmpl.Directories)),
			})
			fmt.Println()
		}
	},
}

var templateApplyCmd = &cobra.Command{
	Use:   "apply [template-name]",
	Short: "Apply a project template to the current directory",
	Long: `Apply a project template by creating the standard directory structure.
If run without arguments, enters interactive mode.

Examples:
  forgectl template apply            # Interactive
  forgectl template apply go         # Direct`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tmplName := ""

		if len(args) == 0 && ui.IsTerminal() {
			ui.Header("Apply Template")
			choices := []ui.Choice{}
			for name, tmpl := range templates.ScaffoldTemplates {
				choices = append(choices, ui.Choice{
					Label:       name,
					Description: tmpl.Description,
					Value:       name,
				})
			}
			tmpl := ui.PromptSelect("Select a template:", choices)
			tmplName = tmpl.Value
		} else if len(args) > 0 {
			tmplName = args[0]
		} else {
			return fmt.Errorf("template name is required")
		}

		tmpl, ok := templates.GetScaffoldTemplate(tmplName)
		if !ok {
			return fmt.Errorf("template %q not found. Run: forgectl template list", tmplName)
		}

		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		ui.Stepf("Applying %s template...", tmplName)

		created := 0
		for _, d := range tmpl.Directories {
			dirPath := filepath.Join(dir, d)
			if err := os.MkdirAll(dirPath, 0o755); err != nil {
				return fmt.Errorf("cannot create directory %s: %w", d, err)
			}
			keepFile := filepath.Join(dirPath, ".gitkeep")
			os.WriteFile(keepFile, []byte(""), 0o644)
			created++
			ui.Dim(d)
		}

		ui.Step("Generating .gitignore...")
		gitignore := templates.GetGitignore(tmplName)
		gitignorePath := filepath.Join(dir, ".gitignore")
		os.WriteFile(gitignorePath, []byte(gitignore), 0o644)

		ui.Successf("Template %q applied: %d directories created", tmplName, created)
		return nil
	},
}

var templateInfoCmd = &cobra.Command{
	Use:   "info [template-name]",
	Short: "Show detailed template information",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tmpl, ok := templates.GetScaffoldTemplate(args[0])
		if !ok {
			return fmt.Errorf("template %q not found", args[0])
		}

		ui.Header("Template: " + tmpl.Name)
		ui.Box(tmpl.Name, []string{
			tmpl.Description,
			fmt.Sprintf("Directories: %d", len(tmpl.Directories)),
		})

		fmt.Println()
		ui.Info("Directory structure:")
		for _, d := range tmpl.Directories {
			ui.Dim("  " + d)
		}

		return nil
	},
}

func init() {
	templateCmd.AddCommand(templateListCmd)
	templateCmd.AddCommand(templateApplyCmd)
	templateCmd.AddCommand(templateInfoCmd)
}
