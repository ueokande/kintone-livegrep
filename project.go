package livegreptone

// Project represents a project
type Project struct {
	Repositories []Repository `json:"repositories"`
}

// Repository represents VCS repository and branch
type Repository struct {
	URL    string `json:"url"`
	branch string `json:"url"`
}
