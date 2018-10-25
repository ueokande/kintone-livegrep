package mock

import (
	"context"

	"github.com/ueokande/livegreptone"
)

func (d *mock) GetProject(ctx context.Context, id string) (livegreptone.Project, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.projects[id], nil
}

func (d *mock) GetProjectIds(ctx context.Context) ([]string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	ids := make([]string, len(d.projects))
	i := 0
	for id, _ := range d.projects {
		ids[i] = id
		i++
	}
	return ids, nil
}

func (d *mock) UpdateProject(ctx context.Context, project livegreptone.Project) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.projects[project.Id] = project
	return nil
}

func (d *mock) RemoveProject(ctx context.Context, id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	delete(d.projects, id)

	return nil
}
