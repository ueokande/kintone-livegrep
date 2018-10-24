package livegreptone

import "time"

// RepositoryStatus represents a repository status
type RepositoryStatus struct {
	UpdatedAt time.Time `json:"updated_at"`
	Commit    string    `json:"commit"`
}
