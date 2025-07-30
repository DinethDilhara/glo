package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/DinethDilhara/glo/internal/formatter"
	"github.com/DinethDilhara/glo/internal/gitexec"
	"github.com/DinethDilhara/glo/internal/models"
	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show git commit history with various formatting options",
	Long: `Display git commit history with filtering and formatting capabilities.

The log command allows you to:
- Filter commits by author, date range, or message content
- Output in different formats (color, JSON, markdown)
- Limit the number of commits shown
- Search within commit messages

Examples:
  glo log                                    # Show recent commits
  glo log --author="John Doe"                # Filter by author
  glo log --since="2024-01-01"               # Show commits since date
  glo log --until="2024-12-31"               # Show commits until date
  glo log --message="fix"                    # Search in commit messages
  glo log --limit=10                         # Limit to 10 commits
  glo log --format=json                      # Output as JSON
  glo log --format=markdown --table          # Output as markdown table`,
	Run: runLogCommand,
}

func runLogCommand(cmd *cobra.Command, args []string) {
	gitExec := gitexec.NewGitExecutor()
	if !gitExec.IsGitRepository() {
		fmt.Fprintf(os.Stderr, "Error: Not a git repository\n")
		os.Exit(1)
	}

	author, _ := cmd.Flags().GetString("author")
	since, _ := cmd.Flags().GetString("since")
	until, _ := cmd.Flags().GetString("until")
	message, _ := cmd.Flags().GetString("message")
	limit, _ := cmd.Flags().GetInt("limit")
	format, _ := cmd.Flags().GetString("format")
	table, _ := cmd.Flags().GetBool("table")
	summary, _ := cmd.Flags().GetBool("summary")
	
	if format == "" {
		format, _ = cmd.Parent().PersistentFlags().GetString("format")
	}

	commits, err := gitExec.GetGitLogs(author, since, until, limit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching git logs: %v\n", err)
		os.Exit(1)
	}

	if message != "" {
		commits = filterCommitsByMessage(commits, message)
	}

	if len(commits) == 0 {
		fmt.Println("No commits found matching the criteria.")
		return
	}

	switch strings.ToLower(format) {
	case "json":
		jsonFormatter := formatter.NewJSONFormatter(true)
		if summary {
			metadata := map[string]interface{}{
				"author": author,
				"since":  since,
				"until":  until,
			}
			fmt.Println(jsonFormatter.FormatSummary(commits, metadata))
		} else {
			fmt.Println(jsonFormatter.FormatList(commits))
		}
	case "markdown", "md":
		mdFormatter := formatter.NewMarkdownFormatter()
		if summary {
			fmt.Println(mdFormatter.FormatSummary(commits))
		} else if table {
			fmt.Println(mdFormatter.FormatTable(commits))
		} else {
			fmt.Println(mdFormatter.FormatList(commits))
		}
	case "color", "":
		colorFormatter := formatter.NewColorFormatter()
		if summary {
			displayColorSummary(commits, colorFormatter)
		} else {
			fmt.Println(colorFormatter.FormatList(commits))
		}
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown format '%s'. Use: color, json, or markdown\n", format)
		os.Exit(1)
	}
}

func filterCommitsByMessage(commits []models.Commit, message string) []models.Commit {
	var filtered []models.Commit
	messageLower := strings.ToLower(message)
	
	for _, commit := range commits {
		if strings.Contains(strings.ToLower(commit.Message), messageLower) {
			filtered = append(filtered, commit)
		}
	}
	
	return filtered
}

func displayColorSummary(commits []models.Commit, colorFormatter *formatter.ColorFormatter) {
	fmt.Println(colorFormatter.FormatHeader("Git Repository Summary"))
	fmt.Printf("Total commits: %d\n\n", len(commits))
	
	authorCount := make(map[string]int)
	for _, commit := range commits {
		authorCount[commit.Author]++
	}
	
	fmt.Println(colorFormatter.FormatHeader("Commits by Author:"))
	for author, count := range authorCount {
		fmt.Printf("  %s: %d commits\n", author, count)
	}
	
	fmt.Println(colorFormatter.FormatHeader("\nRecent Commits:"))
	limit := 5
	if len(commits) < limit {
		limit = len(commits)
	}
	
	for i := 0; i < limit; i++ {
		fmt.Printf("%d. %s\n", i+1, colorFormatter.Format(commits[i]))
	}
}

func init() {
	rootCmd.AddCommand(logCmd)

	logCmd.Flags().StringP("author", "a", "", "Filter commits by author name")
	logCmd.Flags().StringP("since", "s", "", "Show commits since date (YYYY-MM-DD)")
	logCmd.Flags().StringP("until", "u", "", "Show commits until date (YYYY-MM-DD)")
	logCmd.Flags().StringP("message", "m", "", "Filter commits containing message")
	logCmd.Flags().IntP("limit", "l", 0, "Limit number of commits (0 = no limit)")
	logCmd.Flags().StringP("format", "f", "", "Output format: color, json, markdown")
	logCmd.Flags().BoolP("table", "t", false, "Output markdown as table format")
	logCmd.Flags().BoolP("summary", "", false, "Show summary with statistics")
}