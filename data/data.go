package Data


// RepoData defines the structure for a Reposity (info about repository)
type RepoData struct {
	ID         int      `json:"id"`
	RepoURL    string   `json:"repository"`
	Commit     string   `json:"commit"`
	Dockerfile string   `json:"dockerfile"`
	Image      []string `json:"image"`
}

// Repositories is a group / Slice of Data which is a collection of main struct RepoData
type Repositories struct {
	Data []*RepoData
}