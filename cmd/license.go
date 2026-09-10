package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/forgectl/forgectl/internal/templates"
	"github.com/forgectl/forgectl/pkg/ui"
	"github.com/spf13/cobra"
)

var licenseCmd = &cobra.Command{
	Use:   "license",
	Short: "License management commands",
	Long:  "Add and manage licenses for your project.",
}

var licenseAddCmd = &cobra.Command{
	Use:   "add [license-type]",
	Short: "Add a LICENSE file to the project",
	Long: `Add a LICENSE file with the specified license type.
If run without arguments, enters interactive mode.

Examples:
  forgectl license add              # Interactive
  forgectl license add MIT          # Direct
  forgectl license add Apache --author "John Doe"`,
	Args: cobra.MaximumNArgs(1),
	RunE: runLicenseAdd,
}

func runLicenseAdd(cmd *cobra.Command, args []string) error {
	licenseType := ""
	author, _ := cmd.Flags().GetString("author")

	if len(args) == 0 && ui.IsTerminal() {
		ui.Header("Add License")

		choices := []ui.Choice{
			{Label: "MIT", Description: "Short and permissive", Value: "MIT"},
			{Label: "Apache 2.0", Description: "Permissive with patent grant", Value: "Apache"},
			{Label: "GPL v3", Description: "Copyleft, derivative works must be GPL", Value: "GPL"},
			{Label: "BSD 2-Clause", Description: "Permissive, similar to MIT", Value: "BSD"},
			{Label: "ISC", Description: "Functionally identical to MIT", Value: "ISC"},
			{Label: "MPL 2.0", Description: "Weak copyleft, file-level", Value: "MPL"},
			{Label: "Unlicense", Description: "Public domain dedication", Value: "Unlicense"},
		}
		lic := ui.PromptSelect("Select a license:", choices)
		licenseType = lic.Value

		author = ui.PromptText("Author name (optional)", author)
	} else if len(args) > 0 {
		licenseType = args[0]
	} else {
		return fmt.Errorf("license type is required")
	}

	content, err := templates.GetLicense(licenseType, time.Now().Year())
	if err != nil {
		return fmt.Errorf("cannot generate license: %w", err)
	}

	if author != "" {
		content = fmt.Sprintf("Copyright (c) %d %s\n\n%s", time.Now().Year(), author, content)
	}

	ui.Stepf("Adding %s license...", licenseType)
	if err := os.WriteFile("LICENSE", []byte(content), 0o644); err != nil {
		return fmt.Errorf("cannot write LICENSE: %w", err)
	}

	ui.Successf("LICENSE (%s) added successfully", licenseType)
	return nil
}

var licenseListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available license types",
	Run: func(cmd *cobra.Command, args []string) {
		ui.Header("Available Licenses")
		for _, l := range templates.ListLicenses() {
			ui.Step(l)
		}
	},
}

func init() {
	licenseCmd.AddCommand(licenseAddCmd)
	licenseCmd.AddCommand(licenseListCmd)
	licenseAddCmd.Flags().String("author", "", "author name for the license")
}
