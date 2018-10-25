package livegreptone

// Project represents a project
type Project struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Repositories []Repository `json:"repositories"`
	Revision     string       `json:"revision"`
}

// Repository represents VCS repository and branch
type Repository struct {
	URL    string `json:"url"`
	Branch string `json:"branch"`
}
