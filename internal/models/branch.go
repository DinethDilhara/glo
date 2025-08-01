package models

type Branch struct {
	Name               string `json:"name"`
	IsCurrent          bool   `json:"is_current"`
	IsRemote           bool   `json:"is_remote"`
	LastCommitHash     string `json:"last_commit_hash"`
	LastCommitMessage  string `json:"last_commit_message"`
	LastCommitAuthor   string `json:"last_commit_author"`
	LastCommitDate     string `json:"last_commit_date"`
}
