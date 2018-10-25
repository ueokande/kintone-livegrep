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
func (d *mock) UpdateStatus(ctx context.Context, repo string, branch string, status livegreptone.RepositoryStatus) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	key := repo + "/" + branch
	d.statuses[key] = status
	return nil
}
