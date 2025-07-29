package gitexec

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/DinethDilhara/glo/internal/models"
)

type GitExecutor struct{}

func NewGitExecutor() *GitExecutor {
	return &GitExecutor{}
}

func (ge *GitExecutor) GetGitLogs(author, since, until string, maxCount int) ([]models.Commit, error) {
	args := []string{"log", "--pretty=format:%H|%an|%ad|%s", "--date=iso"}
	
	if author != "" {
		args = append(args, "--author="+author)
	}
	if since != "" {
		args = append(args, "--since="+since)
	}
	if until != "" {
		args = append(args, "--until="+until)
	}
	if maxCount > 0 {
		args = append(args, "--max-count="+strconv.Itoa(maxCount))
	}
	
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return nil, err
	}
	
	var commits []models.Commit
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		
		parts := strings.SplitN(line, "|", 4)
		if len(parts) == 4 {
			commits = append(commits, models.Commit{
				Hash:    parts[0],
				Author:  parts[1],
				Date:    parts[2],
				Message: parts[3],
			})
		}
	}
	return commits, nil
}

func (ge *GitExecutor) GetCommitCount() (int, error) {
	out, err := exec.Command("git", "rev-list", "--count", "HEAD").Output()
	if err != nil {
		return 0, err
	}
	
	countStr := strings.TrimSpace(string(out))
	var count int
	_, err = fmt.Sscanf(countStr, "%d", &count)
	return count, err
}

func (ge *GitExecutor) GetCurrentBranch() (string, error) {
	out, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (ge *GitExecutor) IsGitRepository() bool {
	err := exec.Command("git", "status").Run()
	return err == nil
}