package formatters

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/DinethDilhara/glo/internal/formatter"
	"github.com/DinethDilhara/glo/internal/models"
)

type LogFormatter struct{}

func NewLogFormatter() *LogFormatter {
	return &LogFormatter{}
}

func (lf *LogFormatter) FormatJSON(commits []models.Commit) string {
	data, err := json.MarshalIndent(commits, "", "  ")
	if err != nil {
		return "[]"
	}
	return string(data)
}

func (lf *LogFormatter) FormatJSONSummary(commits []models.Commit, metadata map[string]interface{}) string {
	summary := map[string]interface{}{
		"total_commits": len(commits),
		"commits":       commits,
		"metadata":      metadata,
	}
	
	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(data)
}

func (lf *LogFormatter) FormatMarkdown(commits []models.Commit) string {
	var result strings.Builder
	
	result.WriteString("# Git Commit History\n\n")
	result.WriteString(fmt.Sprintf("*Generated on %s*\n\n", time.Now().Format("2006-01-02 15:04:05")))
	result.WriteString(fmt.Sprintf("**Total Commits:** %d\n\n", len(commits)))
	result.WriteString("---\n\n")
	
	for i, commit := range commits {
		result.WriteString(fmt.Sprintf("### %d. %s\n\n", i+1, commit.Message))
		result.WriteString(fmt.Sprintf("- **Hash:** `%s`\n", commit.Hash[:8]))
		result.WriteString(fmt.Sprintf("- **Author:** %s\n", commit.Author))
		result.WriteString(fmt.Sprintf("- **Date:** %s\n\n", commit.Date))
		
		if i < len(commits)-1 {
			result.WriteString("---\n\n")
		}
	}
	
	return result.String()
}

func (lf *LogFormatter) FormatColor(commits []models.Commit) string {
	var result strings.Builder
	
	for i, commit := range commits {
		result.WriteString(fmt.Sprintf("%s%s%s ", formatter.ColorYellow, commit.Hash[:8], formatter.ColorReset))
		result.WriteString(fmt.Sprintf("%s%s%s ", formatter.ColorGreen, commit.Author, formatter.ColorReset))
		result.WriteString(fmt.Sprintf("%s%s%s ", formatter.ColorCyan, commit.Date, formatter.ColorReset))
		result.WriteString(commit.Message)
		
		if i < len(commits)-1 {
			result.WriteString("\n")
		}
	}
	
	return result.String()
}

func (lf *LogFormatter) FormatColorSummary(commits []models.Commit) string {
	var result strings.Builder
	
	result.WriteString(fmt.Sprintf("%sGit Repository Summary%s\n\n", formatter.ColorBold+formatter.ColorBlue, formatter.ColorReset))
	result.WriteString(fmt.Sprintf("Total commits: %d\n\n", len(commits)))
	
	authorCount := make(map[string]int)
	for _, commit := range commits {
		authorCount[commit.Author]++
	}
	
	result.WriteString(fmt.Sprintf("%sCommits by Author:%s\n", formatter.ColorBold+formatter.ColorBlue, formatter.ColorReset))
	for author, count := range authorCount {
		result.WriteString(fmt.Sprintf("  %s: %d commits\n", author, count))
	}
	
	result.WriteString(fmt.Sprintf("\n%sRecent Commits:%s\n", formatter.ColorBold+formatter.ColorBlue, formatter.ColorReset))
	limit := 5
	if len(commits) < limit {
		limit = len(commits)
	}
	
	for i := 0; i < limit; i++ {
		commit := commits[i]
		result.WriteString(fmt.Sprintf("%d. %s%s%s %s%s%s %s%s%s %s\n",
			i+1,
			formatter.ColorYellow, commit.Hash[:8], formatter.ColorReset,
			formatter.ColorGreen, commit.Author, formatter.ColorReset,
			formatter.ColorCyan, commit.Date, formatter.ColorReset,
			commit.Message))
	}
	
	return result.String()
}
