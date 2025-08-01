package formatter

import (
	"encoding/json"

	"github.com/DinethDilhara/glo/internal/models"
)

type JSONFormatter struct {
	Indent bool
}

func NewJSONFormatter(indent bool) *JSONFormatter {
	return &JSONFormatter{Indent: indent}
}

func (jf *JSONFormatter) Format(commit models.Commit) string {
	var data []byte
	var err error
	
	if jf.Indent {
		data, err = json.MarshalIndent(commit, "", "  ")
	} else {
		data, err = json.Marshal(commit)
	}
	
	if err != nil {
		return "{}"
	}
	
	return string(data)
}

func (jf *JSONFormatter) FormatList(commits []models.Commit) string {
	var data []byte
	var err error
	
	if jf.Indent {
		data, err = json.MarshalIndent(commits, "", "  ")
	} else {
		data, err = json.Marshal(commits)
	}
	
	if err != nil {
		return "[]"
	}
	
	return string(data)
}

func (jf *JSONFormatter) FormatSummary(commits []models.Commit, metadata map[string]interface{}) string {
	summary := map[string]interface{}{
		"total_commits": len(commits),
		"commits":       commits,
		"metadata":      metadata,
	}
	
	var data []byte
	var err error
	
	if jf.Indent {
		data, err = json.MarshalIndent(summary, "", "  ")
	} else {
		data, err = json.Marshal(summary)
	}
	
	if err != nil {
		return "{}"
	}
	
	return string(data)
}

func (jf *JSONFormatter) FormatBranches(branches []models.Branch) string {
	var data []byte
	var err error
	
	if jf.Indent {
		data, err = json.MarshalIndent(branches, "", "  ")
	} else {
		data, err = json.Marshal(branches)
	}
	
	if err != nil {
		return "[]"
	}
	
	return string(data)
}