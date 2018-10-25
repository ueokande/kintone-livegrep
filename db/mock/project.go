package mock

import (
	"context"
	"strings"

	"github.com/ueokande/livegreptone"
)

func (d *mock) GetProject(ctx context.Context, id string) (livegreptone.Project, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.projects[id], nil
}

func (d *mock) GetProjectIDs(ctx context.Context) ([]string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	ids := make([]string, len(d.projects))
	i := 0
	for id := range d.projects {
		ids[i] = id
		i++
	}
	return ids, nil
}

func (d *mock) UpdateProject(ctx context.Context, project livegreptone.Project) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.projects[project.ID] = project
	return nil
}

func (d *mock) RemoveProject(ctx context.Context, id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	delete(d.projects, id)

	return nil
}

func (d *mock) GetRepositories(ctx context.Context) ([]livegreptone.Repository, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	repoSet := make(map[string]struct{})
	for _, p := range d.projects {
		for _, r := range p.Repositories {
			key := r.URL + "\x00" + r.Branch
			repoSet[key] = struct{}{}
		}
	}
	repos := make([]livegreptone.Repository, len(repoSet))
	i := 0
	for key := range repoSet {
		kv := strings.Split(key, "\x00")
		repos[i] = livegreptone.Repository{URL: kv[0], Branch: kv[1]}
		i++
	}
	return repos, nil
}

func (d *mock) GetOwnedProjects(ctx context.Context, repo string, branch string) ([]livegreptone.Project, error) {
	var projects []livegreptone.Project
	for _, p := range d.projects {
		for _, r := range p.Repositories {
			if r.URL == repo && r.Branch == branch {
				projects = append(projects, p)
			}
		}
	}
	return projects, nil
}
