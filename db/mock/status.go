package mock

import (
	"context"

	"github.com/ueokande/livegreptone"
	"github.com/ueokande/livegreptone/db"
)

func (d *mock) GetStatus(ctx context.Context, repo string, branch string) (livegreptone.RepositoryStatus, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	key := repo + "/" + branch
	status, ok := d.statuses[key]
	if !ok {
		return livegreptone.RepositoryStatus{}, db.ErrNotFound
	}
	return status, nil
}
func (d *mock) UpdateStatus(ctx context.Context, status livegreptone.RepositoryStatus) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	key := status.URL + "/" + status.Branch
	d.statuses[key] = status
	return nil
}

func (d *mock) WatchStatus(ctx context.Context) <-chan livegreptone.RepositoryStatus {
	// TODO
	return nil
}
