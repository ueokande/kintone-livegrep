package livegreptone

// Project represents a project
type Project struct {
	Id           string       `json:"name"`
	Name         string       `json:"id"`
	Repositories []Repository `json:"repositories"`
	Revision     string       `json:"revision"`
}

// Repository represents VCS repository and branch
type Repository struct {
	URL    string `json:"url"`
	Branch string `json:"branch"`
}
