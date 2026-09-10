package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/forgectl/forgectl/internal/config"
	"github.com/forgectl/forgectl/internal/git"
	"github.com/forgectl/forgectl/pkg/ui"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose system and configuration issues",
	Long: `Check system prerequisites, configuration, and project health.
Use this command to troubleshoot issues with forgectl.

Examples:
  forgectl doctor
  forgectl doctor --fix`,
	RunE: runDoctor,
}

func runDoctor(cmd *cobra.Command, args []string) error {
	ui.Header("System Diagnostics")
	fix, _ := cmd.Flags().GetBool("fix")
	checks := 0
	passed := 0
	warnings := 0
	errors := 0

	// Check 1: Go
	ui.Step("Checking Go installation...")
	if err := checkCommand("go"); err != nil {
		ui.Errorf("Go not found: %v", err)
		errors++
	} else {
		out, _ := exec.Command("go", "version").Output()
		ui.Successf("Go: %s", strings.TrimSpace(string(out)))
		passed++
	}
	checks++

	// Check 2: Git
	ui.Step("Checking Git installation...")
	if err := checkCommand("git"); err != nil {
		ui.Errorf("Git not found: %v", err)
		errors++
	} else {
		out, _ := exec.Command("git", "version").Output()
		ui.Successf("Git: %s", strings.TrimSpace(string(out)))
		passed++
	}
	checks++

	// Check 3: Node.js
	ui.Step("Checking Node.js installation...")
	if err := checkCommand("node"); err != nil {
		ui.Warn("Node.js not found (optional)")
		warnings++
	} else {
		out, _ := exec.Command("node", "--version").Output()
		ui.Successf("Node.js: %s", strings.TrimSpace(string(out)))
		passed++
	}
	checks++

	// Check 4: Python
	ui.Step("Checking Python installation...")
	if err := checkCommand("python3"); err != nil {
		if err := checkCommand("python"); err != nil {
			ui.Warn("Python not found (optional)")
			warnings++
		} else {
			out, _ := exec.Command("python", "--version").Output()
			ui.Successf("Python: %s", strings.TrimSpace(string(out)))
			passed++
		}
	} else {
		out, _ := exec.Command("python3", "--version").Output()
		ui.Successf("Python: %s", strings.TrimSpace(string(out)))
		passed++
	}
	checks++

	// Check 5: Cargo (Rust)
	ui.Step("Checking Rust/Cargo installation...")
	if err := checkCommand("cargo"); err != nil {
		ui.Warn("Cargo not found (optional)")
		warnings++
	} else {
		out, _ := exec.Command("cargo", "--version").Output()
		ui.Successf("Cargo: %s", strings.TrimSpace(string(out)))
		passed++
	}
	checks++

	// Check 6: Configuration
	ui.Step("Checking configuration...")
	cfg, err := config.New()
	if err != nil {
		ui.Errorf("Configuration error: %v", err)
		errors++
	} else {
		appCfg, _ := cfg.GetAppConfig()
		if len(appCfg.Remotes) > 0 {
			ui.Successf("Remotes configured: %d", len(appCfg.Remotes))
			for name, r := range appCfg.Remotes {
				ui.Dim(fmt.Sprintf("  %s (%s) - %s", name, r.Provider, r.Username))
			}
		} else {
			ui.Warn("No remotes configured. Run: forgectl config set")
			warnings++
		}
		passed++
	}
	checks++

	// Check 7: Git config
	ui.Step("Checking Git configuration...")
	if name := getGitConfig("user.name"); name != "" {
		ui.Successf("Git user.name: %s", name)
		passed++
	} else {
		ui.Warn("Git user.name not set")
		warnings++
	}
	checks++

	if email := getGitConfig("user.email"); email != "" {
		ui.Successf("Git user.email: %s", email)
		passed++
	} else {
		ui.Warn("Git user.email not set")
		warnings++
	}
	checks++

	// Check 8: Platform
	ui.Step("Platform information...")
	ui.Successf("OS: %s / Arch: %s", runtime.GOOS, runtime.GOARCH)
	passed++
	checks++

	// Summary
	fmt.Println()
	ui.Header("Diagnostics Summary")
	ui.Box("Results", []string{
		fmt.Sprintf("Total checks:  %d", checks),
		fmt.Sprintf("Passed:        %d", passed),
		fmt.Sprintf("Warnings:      %d", warnings),
		fmt.Sprintf("Errors:        %d", errors),
	})

	if errors > 0 {
		fmt.Println()
		ui.Error("Some checks failed. Please fix the issues above.")
		if fix {
			ui.Info("Attempting to fix issues...")
			fixIssues()
		}
	} else if warnings > 0 {
		fmt.Println()
		ui.Warn("Some warnings found. forgectl may work with reduced functionality.")
	} else {
		fmt.Println()
		ui.Success("All checks passed! forgectl is ready to use.")
	}

	return nil
}

func checkCommand(name string) error {
	_, err := exec.LookPath(name)
	return err
}

func getGitConfig(key string) string {
	out, err := exec.Command("git", "config", "--global", key).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func fixIssues() {
	// Check if git user is configured
	if getGitConfig("user.name") == "" {
		ui.Info("Setting git user.name...")
		fmt.Printf("  Enter your name: ")
		var name string
		fmt.Scanln(&name)
		if name != "" {
			exec.Command("git", "config", "--global", "user.name", name).Run()
			ui.Success("Git user.name set")
		}
	}
	if getGitConfig("user.email") == "" {
		ui.Info("Setting git user.email...")
		fmt.Printf("  Enter your email: ")
		var email string
		fmt.Scanln(&email)
		if email != "" {
			exec.Command("git", "config", "--global", "user.email", email).Run()
			ui.Success("Git user.email set")
		}
	}
}

func init() {
	doctorCmd.Flags().Bool("fix", false, "attempt to fix issues automatically")
}

// openCmd opens the project in browser or file manager
var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open project in browser or file manager",
	Long: `Open the current project's remote URL in the default browser,
or open the file manager at the project directory.

Examples:
  forgectl open              # Open remote URL in browser
  forgectl open --dir        # Open file manager
  forgectl open --repo       # Open repository page`,
	RunE: func(cmd *cobra.Command, args []string) error {
		openDir, _ := cmd.Flags().GetBool("dir")

		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		if openDir {
			return openFileManager(dir)
		}

		// Try to get remote URL
		cfg, err := config.New()
		if err != nil {
			return openFileManager(dir)
		}

		// Find project in config
		projects := cfg.ListProjects()
		for _, p := range projects {
			if p.LocalPath == dir || dir == filepath.Join(p.LocalPath, p.Name) {
				if p.Remote != "" {
					ui.Stepf("Opening %s in browser...", p.Name)
					return openBrowser(p.Remote)
				}
			}
		}

		// Try git remote
		gitDir, err := git.FindRepoRoot(dir)
		if err == nil {
			client, err := git.Open(gitDir)
			if err == nil {
				url, err := client.GetRemoteURL("origin")
				if err == nil {
					ui.Stepf("Opening %s in browser...", url)
					return openBrowser(url)
				}
			}
		}

		ui.Warn("No remote URL found. Opening file manager...")
		return openFileManager(dir)
	},
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	return cmd.Start()
}

func openFileManager(dir string) error {
	ui.Stepf("Opening file manager at %s...", dir)
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", dir)
	case "darwin":
		cmd = exec.Command("open", dir)
	case "windows":
		cmd = exec.Command("explorer", dir)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	return cmd.Start()
}

func init() {
	openCmd.Flags().Bool("dir", false, "open file manager instead of browser")
	openCmd.Flags().Bool("repo", false, "open repository page in browser")
}
