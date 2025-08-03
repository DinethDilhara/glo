package models

type FileStatus struct {
	Path       string `json:"path"`
	Status     string `json:"status"`     
	StatusCode string `json:"statusCode"` 
}

type RepositoryStatus struct {
	Branch       string       `json:"branch"`
	Ahead        int          `json:"ahead"`
	Behind       int          `json:"behind"`
	Staged       []FileStatus `json:"staged"`
	Modified     []FileStatus `json:"modified"`
	Untracked    []FileStatus `json:"untracked"`
	Conflicts    []FileStatus `json:"conflicts"`
	IsClean      bool         `json:"isClean"`
	RemoteBranch string       `json:"remoteBranch,omitempty"`
}

func (f FileStatus) GetStatusDescription() string {
	switch f.Status {
	case "M":
		return "Modified"
	case "A":
		return "Added"
	case "D":
		return "Deleted"
	case "R":
		return "Renamed"
	case "C":
		return "Copied"
	case "U":
		return "Unmerged"
	case "?":
		return "Untracked"
	default:
		return "Unknown"
	}
}
