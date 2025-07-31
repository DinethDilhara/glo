package formatter

import (
	"fmt"
	"strings"

	"github.com/DinethDilhara/glo/internal/models"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
)

type ColorFormatter struct{}

func NewColorFormatter() *ColorFormatter {
	return &ColorFormatter{}
}

func (cf *ColorFormatter) Format(commit models.Commit) string {
	var result strings.Builder
	
	result.WriteString(fmt.Sprintf("%s%s%s ", ColorYellow, commit.Hash[:8], ColorReset))
	
	result.WriteString(fmt.Sprintf("%s%s%s ", ColorGreen, commit.Author, ColorReset))
	
	result.WriteString(fmt.Sprintf("%s%s%s ", ColorCyan, commit.Date, ColorReset))
	
	result.WriteString(commit.Message)
	
	return result.String()
}

func (cf *ColorFormatter) FormatList(commits []models.Commit) string {
	var result strings.Builder
	
	for i, commit := range commits {
		result.WriteString(cf.Format(commit))
		if i < len(commits)-1 {
			result.WriteString("\n")
		}
	}
	
	return result.String()
}

func (cf *ColorFormatter) FormatHeader(text string) string {
	return fmt.Sprintf("%s%s%s%s", ColorBold, ColorBlue, text, ColorReset)
}