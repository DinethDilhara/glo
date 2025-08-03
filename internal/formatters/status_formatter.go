package formatters

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/DinethDilhara/glo/internal/formatter"
	"github.com/DinethDilhara/glo/internal/models"
)

type StatusFormatter struct {
	useColor bool
}

func NewStatusFormatter(useColor bool) *StatusFormatter {
	return &StatusFormatter{
		useColor: useColor,
	}
}

func (sf *StatusFormatter) FormatColor(status *models.RepositoryStatus) string {
	var result strings.Builder
	
	result.WriteString(sf.formatHeader(status))
	result.WriteString("\n")
	
	if status.RemoteBranch != "" {
		result.WriteString(sf.formatRemoteInfo(status))
		result.WriteString("\n")
	}
	
	if status.IsClean {
		cleanMsg := "Working tree clean"
		if sf.useColor {
			cleanMsg = sf.colorize(cleanMsg, formatter.ColorGreen)
		}
		result.WriteString(cleanMsg)
		result.WriteString("\n")
		return result.String()
	}
	
	if len(status.Conflicts) > 0 {
		result.WriteString(sf.formatFileSection("Conflicts", status.Conflicts, formatter.ColorRed))
		result.WriteString("\n")
	}
	
	if len(status.Staged) > 0 {
		result.WriteString(sf.formatFileSection("Staged", status.Staged, formatter.ColorGreen))
		result.WriteString("\n")
	}
	
	if len(status.Modified) > 0 {
		result.WriteString(sf.formatFileSection("Modified", status.Modified, formatter.ColorYellow))
		result.WriteString("\n")
	}
	
	if len(status.Untracked) > 0 {
		result.WriteString(sf.formatFileSection("Untracked", status.Untracked, formatter.ColorCyan))
		result.WriteString("\n")
	}
	
	result.WriteString(sf.formatNextSteps(status))
	
	return result.String()
}

func (sf *StatusFormatter) FormatJSON(status *models.RepositoryStatus) (string, error) {
	jsonData, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (sf *StatusFormatter) FormatTable(status *models.RepositoryStatus) string {
	var result strings.Builder
	
	result.WriteString(sf.formatHeader(status))
	result.WriteString("\n")
	result.WriteString(strings.Repeat("─", 50))
	result.WriteString("\n")
	
	if status.IsClean {
		result.WriteString("Status: Clean - no changes to commit\n")
		return result.String()
	}
	
	result.WriteString(fmt.Sprintf("%-12s %-8s %s\n", "CATEGORY", "STATUS", "FILE"))
	result.WriteString(strings.Repeat("─", 50))
	result.WriteString("\n")
	

	sf.addFilesToTable(&result, "Conflicts", status.Conflicts)
	sf.addFilesToTable(&result, "Staged", status.Staged)
	sf.addFilesToTable(&result, "Modified", status.Modified)
	sf.addFilesToTable(&result, "Untracked", status.Untracked)
	
	return result.String()
}

func (sf *StatusFormatter) FormatSummary(status *models.RepositoryStatus) string {
	if status.IsClean {
		return fmt.Sprintf("%s: clean", status.Branch)
	}
	
	parts := []string{}
	if len(status.Staged) > 0 {
		parts = append(parts, fmt.Sprintf("%d staged", len(status.Staged)))
	}
	if len(status.Modified) > 0 {
		parts = append(parts, fmt.Sprintf("%d modified", len(status.Modified)))
	}
	if len(status.Untracked) > 0 {
		parts = append(parts, fmt.Sprintf("%d untracked", len(status.Untracked)))
	}
	if len(status.Conflicts) > 0 {
		parts = append(parts, fmt.Sprintf("%d conflicts", len(status.Conflicts)))
	}
	
	return fmt.Sprintf("%s: %s", status.Branch, strings.Join(parts, ", "))
}

func (sf *StatusFormatter) colorize(text, color string) string {
	if !sf.useColor {
		return text
	}
	return fmt.Sprintf("%s%s%s", color, text, formatter.ColorReset)
}

func (sf *StatusFormatter) formatHeader(status *models.RepositoryStatus) string {
	header := fmt.Sprintf("Repository Status: %s", status.Branch)
	if sf.useColor {
		header = sf.colorize("Repository Status: ", formatter.ColorCyan) + sf.colorize(status.Branch, formatter.ColorWhite)
	}
	return header
}

func (sf *StatusFormatter) formatRemoteInfo(status *models.RepositoryStatus) string {
	var parts []string
	
	if status.Ahead > 0 {
		aheadStr := fmt.Sprintf("%d ahead", status.Ahead)
		if sf.useColor {
			aheadStr = sf.colorize(aheadStr, formatter.ColorGreen)
		}
		parts = append(parts, aheadStr)
	}
	
	if status.Behind > 0 {
		behindStr := fmt.Sprintf("%d behind", status.Behind)
		if sf.useColor {
			behindStr = sf.colorize(behindStr, formatter.ColorYellow)
		}
		parts = append(parts, behindStr)
	}
	
	if len(parts) == 0 {
		syncStr := "up to date"
		if sf.useColor {
			syncStr = sf.colorize(syncStr, formatter.ColorGreen)
		}
		parts = append(parts, syncStr)
	}
	
	remote := fmt.Sprintf("Remote: %s", status.RemoteBranch)
	if sf.useColor {
		remote = sf.colorize("Remote: ", formatter.ColorCyan) + sf.colorize(status.RemoteBranch, formatter.ColorWhite)
	}
	
	return fmt.Sprintf("%s (%s)", remote, strings.Join(parts, ", "))
}

func (sf *StatusFormatter) formatFileSection(title string, files []models.FileStatus, sectionColor string) string {
	var result strings.Builder
	
	sectionTitle := fmt.Sprintf("%s (%d file", title, len(files))
	if len(files) != 1 {
		sectionTitle += "s"
	}
	sectionTitle += ")"
	
	if sf.useColor {
		sectionTitle = sf.colorize(sectionTitle, sectionColor)
	}
	result.WriteString(sectionTitle)
	result.WriteString("\n")
	
	for _, file := range files {
		fileEntry := fmt.Sprintf("  %s    %s", 
			file.Status, 
			file.Path)
		
		if sf.useColor {
			fileEntry = fmt.Sprintf("  %s    %s", 
				sf.colorize(file.Status, sectionColor), 
				file.Path)
		}
		
		result.WriteString(fileEntry)
		result.WriteString("\n")
	}
	
	return result.String()
}

func (sf *StatusFormatter) addFilesToTable(result *strings.Builder, category string, files []models.FileStatus) {
	for i, file := range files {
		categoryCol := category
		if i > 0 {
			categoryCol = "" 
		}
		
		result.WriteString(fmt.Sprintf("%-12s %-8s %s\n", 
			categoryCol, 
			file.Status+" "+file.GetStatusDescription(), 
			file.Path))
	}
}

func (sf *StatusFormatter) formatNextSteps(status *models.RepositoryStatus) string {
	var suggestions []string
	
	if len(status.Conflicts) > 0 {
		suggestions = append(suggestions, "Resolve conflicts first")
	} else if len(status.Modified) > 0 || len(status.Untracked) > 0 {
		suggestions = append(suggestions, "git add . && git commit -m \"your message\"")
	} else if len(status.Staged) > 0 {
		suggestions = append(suggestions, "git commit -m \"your message\"")
	}
	
	if status.Ahead > 0 {
		suggestions = append(suggestions, "git push")
	}
	
	if len(suggestions) == 0 {
		return ""
	}
	
	nextSteps := "Next: " + strings.Join(suggestions, " → ")
	if sf.useColor {
		nextSteps = sf.colorize("Next: ", formatter.ColorCyan) + sf.colorize(strings.Join(suggestions, " → "), formatter.ColorWhite)
	}
	
	return nextSteps + "\n"
}
