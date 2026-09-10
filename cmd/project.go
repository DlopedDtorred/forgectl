package cmd

import (
	"fmt"
	"os"

	"github.com/forgectl/forgectl/internal/config"
	"github.com/forgectl/forgectl/pkg/ui"
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Project management commands",
	Long:  "List, inspect, and manage tracked projects.",
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tracked projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.New()
		if err != nil {
			return err
		}

		projects := cfg.ListProjects()
		if len(projects) == 0 {
			ui.Info("No projects tracked. Create one with: forgectl new <name>")
			return nil
		}

		ui.Header("Tracked Projects")
		for name, p := range projects {
			ui.Box(name, []string{
				"Description: " + p.Description,
				"Template:    " + p.Template,
				"License:     " + p.License,
				"Provider:    " + p.Provider,
				"Path:        " + p.LocalPath,
			})
			fmt.Println()
		}
		return nil
	},
}

var projectInfoCmd = &cobra.Command{
	Use:   "info [project-name]",
	Short: "Show detailed project information",
	Long: `Show detailed information about a tracked project.
If run without arguments, tries to detect the current directory.

Examples:
  forgectl project info my-project
  forgectl project info              # Detect from current directory`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.New()
		if err != nil {
			return err
		}

		var projectName string
		if len(args) > 0 {
			projectName = args[0]
		} else if ui.IsTerminal() {
			// Try to detect from current directory
			dir, _ := os.Getwd()
			projects := cfg.ListProjects()
			for name, p := range projects {
				if p.LocalPath == dir {
					projectName = name
					break
				}
			}
			if projectName == "" {
				ui.Header("Select Project")
				choices := []ui.Choice{}
				for name, p := range projects {
					choices = append(choices, ui.Choice{
						Label:       name,
						Description: p.Description,
						Value:       name,
					})
				}
				if len(choices) > 0 {
					sel := ui.PromptSelect("Select a project:", choices)
					projectName = sel.Value
				}
			}
		}

		if projectName == "" {
			return fmt.Errorf("project name is required")
		}

		project, err := cfg.GetProject(projectName)
		if err != nil {
			return err
		}

		ui.Header("Project: " + project.Name)
		ui.Box(project.Name, []string{
			"Description: " + project.Description,
			"Template:    " + project.Template,
			"License:     " + project.License,
			"Provider:    " + project.Provider,
			"Remote:      " + project.Remote,
			"Local Path:  " + project.LocalPath,
		})

		// Check if path exists
		if _, err := os.Stat(project.LocalPath); os.IsNotExist(err) {
			fmt.Println()
			ui.Warn("Project directory does not exist at: " + project.LocalPath)
		} else {
			fmt.Println()
			ui.Success("Project directory exists")
		}

		return nil
	},
}

var projectRemoveCmd = &cobra.Command{
	Use:   "remove [project-name]",
	Short: "Remove a project from the tracked list",
	Long: `Remove a project from the tracked list. This does NOT delete
the project files or the remote repository.

Examples:
  forgectl project remove my-project
  forgectl project remove --force`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.New()
		if err != nil {
			return err
		}

		force, _ := cmd.Flags().GetBool("force")
		projectName := ""

		if len(args) > 0 {
			projectName = args[0]
		} else if ui.IsTerminal() {
			projects := cfg.ListProjects()
			if len(projects) == 0 {
				ui.Info("No projects to remove")
				return nil
			}

			ui.Header("Remove Project")
			choices := []ui.Choice{}
			for name, p := range projects {
				choices = append(choices, ui.Choice{
					Label:       name,
					Description: p.Description,
					Value:       name,
				})
			}
			sel := ui.PromptSelect("Select project to remove:", choices)
			projectName = sel.Value
		}

		if projectName == "" {
			return fmt.Errorf("project name is required")
		}

		if !force {
			if !ui.PromptConfirm(fmt.Sprintf("Remove %q from tracked list?", projectName), false) {
				ui.Info("Aborted")
				return nil
			}
		}

		if err := cfg.RemoveProject(projectName); err != nil {
			return err
		}

		ui.Successf("Project %q removed from tracked list", projectName)
		return nil
	},
}

func init() {
	projectCmd.AddCommand(projectListCmd)
	projectCmd.AddCommand(projectInfoCmd)
	projectCmd.AddCommand(projectRemoveCmd)
	projectRemoveCmd.Flags().BoolP("force", "f", false, "skip confirmation")
}
