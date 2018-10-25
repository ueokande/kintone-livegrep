package db

import (
	"context"

	"github.com/ueokande/livegreptone"
)

func (d *model) GetStatus(ctx context.Context, repo string, branch string) (livegreptone.RepositoryStatus, error) {
	// TODO
	return livegreptone.RepositoryStatus{}, nil
}

func (d *model) UpdateStatus(ctx context.Context, repo string, branch string, status livegreptone.RepositoryStatus) error {
	// TODO
	return nil
}
