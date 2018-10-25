package mock

import (
	"context"
	"strconv"
	"sync"

	"github.com/ueokande/livegreptone/git"
	"golang.org/x/exp/rand"
)

type mock struct {
	mu    *sync.Mutex
	repos map[string]string
}

// New returns a mock implementation of git.Interface
func New() git.Interface {
	return &mock{
		mu:    new(sync.Mutex),
		repos: make(map[string]string),
	}
}

func (m *mock) Update(ctx context.Context, url string, branch string) (string, error) {
	key := url + "/" + branch
	commit := strconv.Itoa(rand.Int())
	m.repos[key] = commit
	return commit, nil
}

func (m *mock) Remove(ctx context.Context, url string, branch string) error {
	key := url + "/" + branch
	delete(m.repos, key)
	return nil
}
