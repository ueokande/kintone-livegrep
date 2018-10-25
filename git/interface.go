package git

import "context"

// Interface is an interface of git engine
type Interface interface {
	// Update updates from remove with url and branch to local, and returns HEAD\s commat
	Update(ctx context.Context, url string, branch string) (string, error)

	// Remove removes local repository
	Remove(ctx context.Context, url string, branch string) error
}
