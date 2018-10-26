package livegreptone

// RepositoryStatus represents a repository status
type RepositoryStatus struct {
	URL    string `json:"url"`
	Branch string `json:"branch"`
	Commit string `json:"commit"`
}
