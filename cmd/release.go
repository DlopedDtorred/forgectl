package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/forgectl/forgectl/internal/config"
	"github.com/forgectl/forgectl/internal/git"
	"github.com/forgectl/forgectl/internal/project"
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
  forgectl release new v1.0.0 --push
  forgectl release new v1.0.0 --remote           # Also create on remote provider
  forgectl release new v1.0.0 --remote-only      # Only create on remote provider
  forgectl release new v1.0.0 --remote --draft --prerelease`,
	Args: cobra.MaximumNArgs(1),
	RunE: runReleaseNew,
}

func runReleaseNew(cmd *cobra.Command, args []string) error {
	checkGitInstalled()

	version := ""
	message := ""
	pushRelease := false

	remoteName, _ := cmd.Flags().GetString("remote")
	createRemote := remoteName != ""
	remoteOnly, _ := cmd.Flags().GetBool("remote-only")
	draft, _ := cmd.Flags().GetBool("draft")
	prerelease, _ := cmd.Flags().GetBool("prerelease")

	if len(args) == 0 && ui.IsTerminal() {
		ui.Header("Create Release")

		version = ui.PromptTextRequired("Version (e.g., v1.0.0)")
		message = ui.PromptText("Release message", "Release "+version)
		pushRelease = ui.PromptConfirm("Push to remote?", true)
		if !remoteOnly {
			useRemote := ui.PromptConfirm("Create release on remote provider?", false)
			if useRemote {
				createRemote = true
				draft = ui.PromptConfirm("Is this a draft release?", false)
				if !draft {
					prerelease = ui.PromptConfirm("Is this a prerelease?", false)
				}
			}
		}
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

	if remoteOnly {
		createRemote = true
	}

	if !remoteOnly {
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
	}

	// Create release on remote provider
	if createRemote {
		provider, err := getRemoteProvider(remoteName)
		if err != nil {
			return fmt.Errorf("cannot get remote provider: %w", err)
		}

		repoName := detectRepoName()
		if repoName == "" {
			return fmt.Errorf("cannot detect repository name. Ensure you are in a git repository with a remote")
		}

		ui.Stepf("Creating release %s on %s...", version, provider.Name())
		release, err := provider.CreateRelease(repoName, version, version, message, draft, prerelease)
		if err != nil {
			return fmt.Errorf("cannot create remote release: %w", err)
		}
		ui.Successf("Release created: %s", release.URL)
	}

	ui.Successf("Release %s created successfully", version)
	return nil
}

var releaseListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all releases (tags or remote)",
	RunE: func(cmd *cobra.Command, args []string) error {
		checkGitInstalled()

		remoteName, _ := cmd.Flags().GetString("remote")
		useRemote := remoteName != ""
		jsonOutput := isJSONOutput(cmd)

		if useRemote {
			provider, err := getRemoteProvider(remoteName)
			if err != nil {
				return fmt.Errorf("cannot get remote provider: %w", err)
			}

			repoName := detectRepoName()
			if repoName == "" {
				return fmt.Errorf("cannot detect repository name")
			}

			releases, err := provider.ListReleases(repoName)
			if err != nil {
				return fmt.Errorf("cannot list remote releases: %w", err)
			}

			if len(releases) == 0 {
				if jsonOutput {
					fmt.Println("[]")
				} else {
					ui.Info("No releases found on remote.")
				}
				return nil
			}

			if jsonOutput {
				data, _ := json.MarshalIndent(releases, "", "  ")
				fmt.Println(string(data))
				return nil
			}

			ui.Header("Remote Releases")
			for _, r := range releases {
				status := ""
				if r.Draft {
					status = " [draft]"
				}
				if r.Prerelease {
					status = " [prerelease]"
				}
				ui.Stepf("%s%s - %s", r.TagName, status, r.URL)
			}
			return nil
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

		tags, err := client.Tags()
		if err != nil {
			return fmt.Errorf("cannot list tags: %w", err)
		}

		if len(tags) == 0 {
			if jsonOutput {
				fmt.Println("[]")
			} else {
				ui.Info("No releases found. Create one with: forgectl release new <version>")
			}
			return nil
		}

		if jsonOutput {
			type tagInfo struct {
				Tag string `json:"tag"`
			}
			var tagList []tagInfo
			for _, t := range tags {
				tagList = append(tagList, tagInfo{Tag: t})
			}
			data, _ := json.MarshalIndent(tagList, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		ui.Header("Releases")
		for _, tag := range tags {
			ui.Step(tag)
		}
		return nil
	},
}

func getRemoteProvider(remoteName string) (project.Provider, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	var remote *config.RemoteConfig
	if remoteName != "" {
		remote, err = cfg.GetRemote(remoteName)
	} else {
		remote, err = cfg.GetDefaultRemote()
	}
	if err != nil {
		return nil, err
	}

	return project.NewProvider(project.RemoteConfig{
		Provider:     remote.Provider,
		Token:        remote.Token,
		Username:     remote.Username,
		BaseURL:      remote.BaseURL,
		Organization: remote.Organization,
	})
}

func detectRepoName() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	repoRoot, err := git.FindRepoRoot(dir)
	if err != nil {
		return ""
	}
	return filepath.Base(repoRoot)
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
	releaseNewCmd.Flags().String("remote", "", "create release on remote provider (remote name from config)")
	releaseNewCmd.Flags().Bool("remote-only", false, "create release ONLY on remote provider (no local tag)")
	releaseNewCmd.Flags().Bool("draft", false, "mark release as draft")
	releaseNewCmd.Flags().Bool("prerelease", false, "mark release as prerelease")

	releaseListCmd.Flags().String("remote", "", "list releases from remote provider (remote name from config)")
}
