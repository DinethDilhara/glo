package parser

import (
	"strings"
	"time"

	"github.com/DinethDilhara/glo/internal/models"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseGitLogOutput(output string) ([]models.Commit, error) {
	var commits []models.Commit
	lines := strings.Split(strings.TrimSpace(output), "\n")
	
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		
		commit, err := p.parseCommitLine(line)
		if err != nil {
			continue
		}
		commits = append(commits, commit)
	}
	
	return commits, nil
}

func (p *Parser) parseCommitLine(line string) (models.Commit, error) {
	parts := strings.SplitN(line, "|", 4)
	if len(parts) < 4 {
		return models.Commit{}, nil
	}
	
	return models.Commit{
		Hash:    strings.TrimSpace(parts[0]),
		Author:  strings.TrimSpace(parts[1]),
		Date:    strings.TrimSpace(parts[2]),
		Message: strings.TrimSpace(parts[3]),
	}, nil
}

func (p *Parser) FilterCommits(commits []models.Commit, author, message string, since time.Time) []models.Commit {
	var filtered []models.Commit
	
	for _, commit := range commits {

		if author != "" && !strings.Contains(strings.ToLower(commit.Author), strings.ToLower(author)) {
			continue
		}
		
		if message != "" && !strings.Contains(strings.ToLower(commit.Message), strings.ToLower(message)) {
			continue
		}
		
		if !since.IsZero() {
			commitTime, err := time.Parse("2006-01-02 15:04:05 -0700", commit.Date)
			if err == nil && commitTime.Before(since) {
				continue
			}
		}
		
		filtered = append(filtered, commit)
	}
	
	return filtered
}