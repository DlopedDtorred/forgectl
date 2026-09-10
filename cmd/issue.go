package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/forgectl/forgectl/pkg/ui"
	"github.com/spf13/cobra"
)

var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Issue management commands",
	Long:  "Create and manage issues on remote providers.",
}

var issueNewCmd = &cobra.Command{
	Use:   "new [title]",
	Short: "Create a new issue on remote provider",
	Long: `Create a new issue on the configured remote provider.

Examples:
  forgectl issue new                         # Interactive
  forgectl issue new "Bug report"            # Direct
  forgectl issue new "Feature" --body "..."  # With body
  forgectl issue new "Task" --labels "bug,help wanted"`,
	Args: cobra.MaximumNArgs(1),
	RunE: runIssueNew,
}

var issueListCmd = &cobra.Command{
	Use:   "list",
	Short: "List issues from remote provider",
	Long:  "List issues from the configured remote provider.",
	RunE:  runIssueList,
}

func runIssueNew(cmd *cobra.Command, args []string) error {
	title := ""
	body, _ := cmd.Flags().GetString("body")
	labelsStr, _ := cmd.Flags().GetString("labels")
	remoteName, _ := cmd.Flags().GetString("remote")
	jsonOutput := isJSONOutput(cmd)

	if len(args) > 0 {
		title = args[0]
	} else if ui.IsTerminal() {
		ui.Header("Create Issue")
		title = ui.PromptTextRequired("Issue title")
		body = ui.PromptText("Issue body", "")
		labelsStr = ui.PromptText("Labels (comma-separated)", "")
	} else {
		return fmt.Errorf("issue title is required as argument")
	}

	if title == "" {
		return fmt.Errorf("issue title is required")
	}

	var labels []string
	if labelsStr != "" {
		for _, l := range strings.Split(labelsStr, ",") {
			l = strings.TrimSpace(l)
			if l != "" {
				labels = append(labels, l)
			}
		}
	}

	provider, err := getRemoteProvider(remoteName)
	if err != nil {
		return fmt.Errorf("cannot get remote provider: %w", err)
	}

	repoName := detectRepoName()
	if repoName == "" {
		return fmt.Errorf("cannot detect repository name. Ensure you are in a git repository with a remote")
	}

	ui.Stepf("Creating issue on %s...", provider.Name())
	issue, err := provider.CreateIssue(repoName, title, body, labels)
	if err != nil {
		return fmt.Errorf("cannot create issue: %w", err)
	}

	if jsonOutput {
		data, _ := json.MarshalIndent(issue, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	ui.Successf("Issue #%d created: %s", issue.Number, issue.URL)
	return nil
}

func runIssueList(cmd *cobra.Command, args []string) error {
	remoteName, _ := cmd.Flags().GetString("remote")
	jsonOutput := isJSONOutput(cmd)

	provider, err := getRemoteProvider(remoteName)
	if err != nil {
		return fmt.Errorf("cannot get remote provider: %w", err)
	}

	repoName := detectRepoName()
	if repoName == "" {
		return fmt.Errorf("cannot detect repository name")
	}

	ui.Stepf("Listing issues from %s...", provider.Name())

	issues, err := provider.ListIssues(repoName)
	if err != nil {
		return fmt.Errorf("cannot list issues: %w", err)
	}

	if len(issues) == 0 {
		if jsonOutput {
			fmt.Println("[]")
		} else {
			ui.Info("No open issues found.")
		}
		return nil
	}

	if jsonOutput {
		data, _ := json.MarshalIndent(issues, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	ui.Header("Issues")
	for _, issue := range issues {
		ui.Stepf("#%d %s - %s", issue.Number, issue.Title, issue.URL)
	}
	return nil
}

func init() {
	issueCmd.AddCommand(issueNewCmd)
	issueCmd.AddCommand(issueListCmd)

	issueNewCmd.Flags().String("body", "", "issue body/description")
	issueNewCmd.Flags().String("labels", "", "comma-separated list of labels")
	issueNewCmd.Flags().String("remote", "", "remote name from config")

	issueListCmd.Flags().String("remote", "", "remote name from config")
}
