package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func SetVersionInfo(v, c, d string) {
	version = v
	commit = c
	date = d
}

var rootCmd = &cobra.Command{
	Use:   "glo",
	Short: "A CLI tool to explore Git history with style",
	Long: `glo is a powerful CLI tool for exploring Git history with beautiful formatting options.

Filter commits by author, date, or message and export results in different formats:
- Colorized terminal output for easy reading
- JSON format for programmatic use
- Markdown format for documentation

Examples:
  glo log                              # Show recent commits with colors
  glo log --author="John Doe"          # Filter by author
  glo log --since="2024-01-01"         # Show commits since date
  glo log --format=json                # Export as JSON
  glo log --format=markdown            # Export as Markdown`,
}


var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Print version, commit hash, and build date information for glo.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("glo version %s\n", version)
		fmt.Printf("Git commit: %s\n", commit)
		fmt.Printf("Built: %s\n", date)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	
	rootCmd.PersistentFlags().StringP("format", "f", "color", "Output format: color, json, markdown")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
}
