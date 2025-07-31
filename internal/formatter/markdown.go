package formatter

import (
	"fmt"
	"strings"
	"time"

	"github.com/DinethDilhara/glo/internal/models"
)

type MarkdownFormatter struct{}

func NewMarkdownFormatter() *MarkdownFormatter {
	return &MarkdownFormatter{}
}

func (mf *MarkdownFormatter) Format(commit models.Commit) string {
	var result strings.Builder
	
	result.WriteString(fmt.Sprintf("## %s\n\n", commit.Message))
	result.WriteString(fmt.Sprintf("**Hash:** `%s`\n\n", commit.Hash[:8]))
	result.WriteString(fmt.Sprintf("**Author:** %s\n\n", commit.Author))
	result.WriteString(fmt.Sprintf("**Date:** %s\n\n", commit.Date))
	result.WriteString("---\n\n")
	
	return result.String()
}

func (mf *MarkdownFormatter) FormatList(commits []models.Commit) string {
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

func (mf *MarkdownFormatter) FormatTable(commits []models.Commit) string {
	var result strings.Builder
	
	result.WriteString("# Git Commit History\n\n")
	result.WriteString("| Hash | Author | Date | Message |\n")
	result.WriteString("|------|--------|------|---------|\n")
	
	for _, commit := range commits {
		result.WriteString(fmt.Sprintf("| `%s` | %s | %s | %s |\n",
			commit.Hash[:8],
			commit.Author,
			commit.Date,
			strings.ReplaceAll(commit.Message, "|", "\\|"))) 
	}
	
	return result.String()
}

func (mf *MarkdownFormatter) FormatSummary(commits []models.Commit) string {
	var result strings.Builder
	
	authorCount := make(map[string]int)
	for _, commit := range commits {
		authorCount[commit.Author]++
	}
	
	result.WriteString("# Git Repository Summary\n\n")
	result.WriteString(fmt.Sprintf("**Total Commits:** %d\n\n", len(commits)))
	result.WriteString("## Commits by Author\n\n")
	
	for author, count := range authorCount {
		result.WriteString(fmt.Sprintf("- **%s:** %d commits\n", author, count))
	}
	
	result.WriteString("\n---\n\n")
	result.WriteString("## Recent Commits\n\n")
	
	limit := 5
	if len(commits) < limit {
		limit = len(commits)
	}
	
	for i := 0; i < limit; i++ {
		commit := commits[i]
		result.WriteString(fmt.Sprintf("%d. **%s** by %s (`%s`)\n",
			i+1, commit.Message, commit.Author, commit.Hash[:8]))
	}
	
	return result.String()
}