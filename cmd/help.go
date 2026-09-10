package cmd

import (
	"fmt"

	"github.com/forgectl/forgectl/pkg/ui"
	"github.com/spf13/cobra"
)

var helpCmd = &cobra.Command{
	Use:   "help [command]",
	Short: "Help about any command",
	Long: `Show help for forgectl or for any of its commands.

Examples:
  forgectl help                 # Show overview
  forgectl help new             # Help for a specific command
  forgectl help config set      # Help for subcommands`,
	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// If a command is specified, delegate to Cobra's standard help
		if len(args) > 0 {
			target, _, err := rootCmd.Find(args)
			if err != nil || target == nil {
				ui.Errorf("Unknown help topic %q", args[0])
				fmt.Println()
				helpOverview()
				return
			}
			target.HelpFunc()(target, args)
			return
		}
		helpOverview()
	},
}

func helpOverview() {
	ui.Header("forgectl — Universal Project Management CLI")

	ui.Info("What is forgectl?")
	ui.Dim("  Create, initialize and manage complete development projects")
	ui.Dim("  locally and on remote platforms (GitHub, GitLab, Gitea).")
	fmt.Println()

	ui.Box("Quick Start", []string{
		"",
		"  1. forgectl remote add          Configure your provider",
		"  2. forgectl new                 Create a project (wizard)",
		"  3. cd my-project                Enter the project",
		"  4. forgectl sync                Sync local <-> remote",
		"  5. forgectl release new v1.0.0  Tag & release",
		"",
	})
	fmt.Println()

	ui.Newline()

	groups := [][]*cobra.Command{
		{
			newCmd,
			templateCmd,
			licenseCmd,
			readmeCmd,
		},
		{
			initCmd,
			syncCmd,
			pushCmd,
			pullCmd,
			statusCmd,
			tagCmd,
		},
		{
			releaseCmd,
			issueCmd,
			changelogCmd,
		},
		{
			configCmd,
			remoteCmd,
			projectCmd,
		},
		{
			doctorCmd,
			openCmd,
			searchCmd,
		},
	}

	groupTitles := []string{
		"  Project creation & scaffolding",
		"  Git operations",
		"  Releases, issues & changelog",
		"  Configuration & management",
		"  Utilities",
	}

	for i, group := range groups {
		fmt.Println(colorizeok(groupTitles[i]))
		for _, c := range group {
			name := colorizeCmd("    " + c.Name())
			desc := ui.White + c.Short + ui.Reset
			fmt.Printf("%-24s %s\n", name, desc)
		}
		fmt.Println()
	}

	ui.Box("Supported Templates", []string{
		"",
		"  go         Standard Go project layout",
		"  python     Standard Python project layout",
		"  node       Standard Node.js project layout",
		"  typescript Standard TypeScript project layout",
		"  rust       Standard Rust project layout",
		"  java       Standard Java project layout",
		"  kotlin     Standard Kotlin project layout",
		"  csharp     Standard .NET project layout",
		"  php        Standard PHP project layout",
		"  ruby       Standard Ruby project layout",
		"  swift      Standard Swift project layout",
		"  dart       Standard Dart project layout",
		"  c-cpp      Standard C/C++ project layout",
		"  minimal    Bare minimum structure",
		"",
	})

	ui.Box("Supported Licenses", []string{
		"",
		"  MIT        Short and permissive",
		"  Apache     Apache 2.0 with patent grant",
		"  GPL        GPL v3 copyleft",
		"  BSD        BSD 2-Clause permissive",
		"  ISC        Functionally identical to MIT",
		"  MPL        MPL 2.0 weak copyleft",
		"  Unlicense  Public domain dedication",
		"",
	})

	ui.Box("Examples", []string{
		"",
		"  # Interactive wizard                       ",
		"  forgectl new                               ",
		"  forgectl config set                        ",
		"  forgectl release new                       ",
		"  forgectl issue new \"Bug report\"            ",
		"",
		"  # Direct, scriptable                        ",
		"  forgectl new app -t go -l MIT -d \"My app\"  ",
		"  forgectl sync --tags                        ",
		"  forgectl doctor                             ",
		"  forgectl remote list                        ",
		"  forgectl release new v1.0.0 --remote        ",
		"  forgectl issue new \"Fix\" --labels bug       ",
		"",
	})
	fmt.Println()

	ui.Dim("Run 'forgectl help <command>' for details on a specific command.")
	ui.Dim("Run 'forgectl <command> --help' for the same.")
	fmt.Println()
}

func colorizeCmd(s string) string {
	return ui.Cyan + s + ui.Reset
}

func colorizeok(s string) string {
	return ui.Bold + ui.Purple + s + ui.Reset
}

func init() {
	rootCmd.SetHelpCommand(helpCmd)
}
