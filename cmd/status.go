package cmd

import (
	"fmt"
	"os"

	"github.com/DinethDilhara/glo/internal/formatters"
	"github.com/DinethDilhara/glo/internal/gitexec"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the working tree status with enhanced formatting",
	Long: `Display the current status of the Git repository with beautiful formatting.

Shows staged files, modified files, untracked files, and remote tracking information
in a clean, easy-to-read format with colors and helpful next-step suggestions.

Output formats:
- color (default): Colorized output with emojis and helpful sections
- table: Clean table format for easy scanning
- summary: Brief one-line summary
- json: Machine-readable JSON format

Examples:
  glo status                    # Colorized status with emojis
  glo status --format=table     # Clean table format
  glo status --format=summary   # Brief summary
  glo status --format=json      # JSON output for scripts`,
	Run: runStatus,
}

func runStatus(cmd *cobra.Command, args []string) {
	format, _ := cmd.Flags().GetString("format")
	
	gitExec := gitexec.NewGitExecutor()
	
	if !gitExec.IsGitRepository() {
		fmt.Fprintf(os.Stderr, "Error: Not a git repository\n")
		fmt.Fprintf(os.Stderr, "Run 'git init' to initialize a new repository\n")
		os.Exit(1)
	}
	
	status, err := gitExec.GetRepositoryStatus()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting repository status: %v\n", err)
		os.Exit(1)
	}
	
	useColor := format == "color"
	statusFormatter := formatters.NewStatusFormatter(useColor)
	
	var output string
	switch format {
	case "json":
		jsonOutput, err := statusFormatter.FormatJSON(status)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error formatting JSON: %v\n", err)
			os.Exit(1)
		}
		output = jsonOutput
	case "table":
		output = statusFormatter.FormatTable(status)
	case "summary":
		output = statusFormatter.FormatSummary(status)
	case "color":
		fallthrough
	default:
		output = statusFormatter.FormatColor(status)
	}
	
	fmt.Print(output)
}

func init() {
	rootCmd.AddCommand(statusCmd)
	
	statusCmd.Flags().StringP("format", "f", "color", "Output format: color, table, summary, json")
}
