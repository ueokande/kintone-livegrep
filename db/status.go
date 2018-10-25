package db

import (
	"context"

	"github.com/ueokande/livegreptone"
)

// GetStatus returns an status of the repository from etcd
func (d *model) GetStatus(ctx context.Context, repo string, branch string) (livegreptone.RepositoryStatus, error) {
	// TODO
	return livegreptone.RepositoryStatus{}, nil
}

// UpdateStatus creates or update status of the repository from etcd
func (d *model) UpdateStatus(ctx context.Context, repo string, branch string, status livegreptone.RepositoryStatus) error {
	// TODO
	return nil
}
