package mock

import (
	"sync"

	"github.com/ueokande/livegreptone"
	"github.com/ueokande/livegreptone/db"
)

type mock struct {
	mu       *sync.Mutex
	projects map[string]livegreptone.Project
	statuses map[string]livegreptone.RepositoryStatus
}

// New returns mock implementation of db.Interface
func New() db.Interface {
	return &mock{
		mu:       new(sync.Mutex),
		projects: make(map[string]livegreptone.Project),
		statuses: make(map[string]livegreptone.RepositoryStatus),
	}
}
