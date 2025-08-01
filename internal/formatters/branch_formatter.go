package formatters

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/DinethDilhara/glo/internal/formatter"
	"github.com/DinethDilhara/glo/internal/models"
)

type BranchFormatter struct{}

func NewBranchFormatter() *BranchFormatter {
	return &BranchFormatter{}
}

func (bf *BranchFormatter) FormatJSON(branches []models.Branch) string {
	data, err := json.MarshalIndent(branches, "", "  ")
	if err != nil {
		return "[]"
	}
	return string(data)
}

func (bf *BranchFormatter) FormatTable(branches []models.Branch) string {
	var result strings.Builder
	
	result.WriteString(fmt.Sprintf("%sGit Branches%s\n\n", formatter.ColorBold+formatter.ColorBlue, formatter.ColorReset))
	
	result.WriteString(fmt.Sprintf("%-20s %-10s %-15s %-50s %s\n", "Branch", "Type", "Last Commit", "Message", "Author"))
	result.WriteString(strings.Repeat("-", 120) + "\n")
	
	for _, branch := range branches {
		typeColor := formatter.ColorGreen
		if branch.IsRemote {
			typeColor = formatter.ColorRed
		}
		if branch.IsCurrent {
			typeColor = formatter.ColorYellow
		}
		
		branchType := "local"
		if branch.IsRemote {
			branchType = "remote"
		}
		if branch.IsCurrent {
			branchType = "current"
		}
		
		result.WriteString(fmt.Sprintf("%-20s %s%-10s%s %-15s %-50s %s\n",
			branch.Name,
			typeColor, branchType, formatter.ColorReset,
			branch.LastCommitDate,
			truncateString(branch.LastCommitMessage, 48),
			branch.LastCommitAuthor))
	}
	
	return result.String()
}

func (bf *BranchFormatter) FormatTree(branches []models.Branch, withDates bool) string {
	var result strings.Builder
	
	result.WriteString(fmt.Sprintf("%sGit Branch Tree%s\n\n", formatter.ColorBold+formatter.ColorBlue, formatter.ColorReset))
	
	local := make([]models.Branch, 0)
	remote := make([]models.Branch, 0)
	current := ""
	
	for _, branch := range branches {
		if branch.IsCurrent {
			current = branch.Name
		}
		if branch.IsRemote {
			remote = append(remote, branch)
		} else {
			local = append(local, branch)
		}
	}
	
	result.WriteString(fmt.Sprintf("%s📁 Repository%s\n", formatter.ColorBold, formatter.ColorReset))
	result.WriteString("│\n")
	
	if len(local) > 0 {
		result.WriteString(fmt.Sprintf("├── %s🌿 Local Branches%s\n", formatter.ColorGreen, formatter.ColorReset))
		for i, branch := range local {
			isLast := i == len(local)-1 && len(remote) == 0
			prefix := "│   ├── "
			if isLast {
				prefix = "│   └── "
			}
			
			branchColor := formatter.ColorGreen
			indicator := "  "
			if branch.IsCurrent {
				branchColor = formatter.ColorYellow
				indicator = "* "
			}
			
			result.WriteString(fmt.Sprintf("%s%s%s%s%s", prefix, indicator, branchColor, branch.Name, formatter.ColorReset))
			if withDates {
				result.WriteString(fmt.Sprintf(" %s(%s)%s", formatter.ColorCyan, branch.LastCommitDate, formatter.ColorReset))
			}
			result.WriteString("\n")
		}
	}
	
	if len(remote) > 0 {
		if len(local) > 0 {
			result.WriteString("│\n")
		}
		result.WriteString(fmt.Sprintf("└── %s🌐 Remote Branches%s\n", formatter.ColorRed, formatter.ColorReset))
		for i, branch := range remote {
			isLast := i == len(remote)-1
			prefix := "    ├── "
			if isLast {
				prefix = "    └── "
			}
			
			result.WriteString(fmt.Sprintf("%s%s%s%s", prefix, formatter.ColorRed, branch.Name, formatter.ColorReset))
			if withDates {
				result.WriteString(fmt.Sprintf(" %s(%s)%s", formatter.ColorCyan, branch.LastCommitDate, formatter.ColorReset))
			}
			result.WriteString("\n")
		}
	}
	
	result.WriteString("\n")
	result.WriteString(fmt.Sprintf("%sCurrent branch: %s%s%s\n", formatter.ColorBold, formatter.ColorYellow, current, formatter.ColorReset))
	
	return result.String()
}

func (bf *BranchFormatter) FormatColor(branches []models.Branch, withDates bool) string {
	var result strings.Builder
	
	for _, branch := range branches {
		prefix := "  "
		color := formatter.ColorGreen
		
		if branch.IsCurrent {
			prefix = "* "
			color = formatter.ColorYellow
		} else if branch.IsRemote {
			color = formatter.ColorRed
		}
		
		result.WriteString(fmt.Sprintf("%s%s%s%s", prefix, color, branch.Name, formatter.ColorReset))
		
		if withDates {
			result.WriteString(fmt.Sprintf(" %s(%s)%s", formatter.ColorCyan, branch.LastCommitDate, formatter.ColorReset))
		}
		
		result.WriteString("\n")
	}
	
	return result.String()
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
