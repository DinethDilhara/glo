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

func (ge *GitExecutor) GetBranches(all, remoteOnly bool) ([]models.Branch, error) {
	args := []string{"branch"}
	
	if all {
		args = append(args, "-a")
	} else if remoteOnly {
		args = append(args, "-r")
	}
	
	args = append(args, "-v", "--format=%(refname:short)|%(HEAD)|%(objectname:short)|%(authordate:short)|%(authorname)|%(contents:subject)")
	
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return nil, err
	}
	
	var branches []models.Branch
	lines := strings.Split(string(out), "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		parts := strings.SplitN(line, "|", 6)
		if len(parts) >= 4 {
			branch := models.Branch{
				Name:             parts[0],
				IsCurrent:        parts[1] == "*",
				IsRemote:         strings.Contains(parts[0], "origin/") || strings.Contains(parts[0], "remote/"),
				LastCommitHash:   parts[2],
				LastCommitDate:   parts[3],
			}
			
			if len(parts) >= 5 {
				branch.LastCommitAuthor = parts[4]
			}
			if len(parts) >= 6 {
				branch.LastCommitMessage = parts[5]
			}
			
			branches = append(branches, branch)
		}
	}
	
	if len(branches) == 0 {
		return ge.getBranchesBasic(all, remoteOnly)
	}
	
	return branches, nil
}

func (ge *GitExecutor) getBranchesBasic(all, remoteOnly bool) ([]models.Branch, error) {
	args := []string{"branch"}
	
	if all {
		args = append(args, "-a")
	} else if remoteOnly {
		args = append(args, "-r")
	}
	
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return nil, err
	}
	
	var branches []models.Branch
	lines := strings.Split(string(out), "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		isCurrent := strings.HasPrefix(line, "*")
		if isCurrent {
			line = strings.TrimSpace(line[1:])
		}
		
		line = strings.TrimSpace(line)
		
		branch := models.Branch{
			Name:      line,
			IsCurrent: isCurrent,
			IsRemote:  strings.Contains(line, "origin/") || strings.Contains(line, "remotes/"),
		}
		
		if commitInfo, err := ge.getLastCommitForBranch(line); err == nil {
			branch.LastCommitHash = commitInfo.Hash
			branch.LastCommitMessage = commitInfo.Message
			branch.LastCommitAuthor = commitInfo.Author
			branch.LastCommitDate = commitInfo.Date
		}
		
		branches = append(branches, branch)
	}
	
	return branches, nil
}

func (ge *GitExecutor) getLastCommitForBranch(branchName string) (*models.Commit, error) {
	out, err := exec.Command("git", "log", "-1", "--pretty=format:%H|%an|%ad|%s", "--date=short", branchName).Output()
	if err != nil {
		return nil, err
	}
	
	line := strings.TrimSpace(string(out))
	parts := strings.SplitN(line, "|", 4)
	if len(parts) == 4 {
		return &models.Commit{
			Hash:    parts[0][:8], 
			Author:  parts[1],
			Date:    parts[2],
			Message: parts[3],
		}, nil
	}
	
	return nil, fmt.Errorf("invalid commit format")
}

func (ge *GitExecutor) GetCommitGraph(limit int) ([]models.Commit, error) {
	args := []string{"log", "--graph", "--oneline", "--decorate", "--all", "--pretty=format:%H|%an|%ad|%s", "--date=short"}
	
	if limit > 0 {
		args = append(args, "--max-count="+strconv.Itoa(limit))
	}
	
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return nil, err
	}
	
	var commits []models.Commit
	lines := strings.Split(string(out), "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		cleanLine := strings.TrimLeft(line, "* |\\/_")
		parts := strings.SplitN(cleanLine, "|", 4)
		
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