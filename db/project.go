package db

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/coreos/etcd/clientv3"
	"github.com/ueokande/livegreptone"
)

// GetProject returns a project of id from etcd
func (d *model) GetProject(ctx context.Context, id string) (livegreptone.Project, error) {
	key := ProjectKey(id)
	resp, err := d.etcd.Get(ctx, key)
	if err != nil {
		return livegreptone.Project{}, err
	}
	if resp.Count == 0 {
		return livegreptone.Project{}, ErrNotFound
	}

	var p livegreptone.Project
	err = json.Unmarshal([]byte(resp.Kvs[0].Value), &p)
	if err != nil {
		return livegreptone.Project{}, err
	}
	return p, nil
}

// GetProjectIDs returns project IDs from etcd
func (d *model) GetProjectIDs(ctx context.Context) ([]string, error) {
	resp, err := d.etcd.Get(ctx, ProjectKeyPrefix,
		clientv3.WithPrefix(), clientv3.WithKeysOnly())
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(resp.Kvs))
	for i, kv := range resp.Kvs {
		ids[i] = string(kv.Key)[len(ProjectKeyPrefix):]
	}
	return ids, nil
}

// UpdateProject create or update a project of id from etcd
func (d *model) UpdateProject(ctx context.Context, project livegreptone.Project) error {
	key := ProjectKey(project.ID)
	value, err := json.Marshal(project)
	if err != nil {
		return err
	}
	_, err = d.etcd.Put(ctx, key, string(value))
	if err != nil {
		return err
	}
	return nil
}

// RemoveProject removes a project of id from metcd
func (d *model) RemoveProject(ctx context.Context, id string) error {
	key := ProjectKey(id)
	_, err := d.etcd.Delete(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

// GetRepositories gets repository list in projects
func (d *model) GetRepositories(ctx context.Context) ([]livegreptone.Repository, error) {
	resp, err := d.etcd.Get(ctx, ProjectKeyPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	repoSet := make(map[string]struct{})
	for _, kv := range resp.Kvs {
		var project livegreptone.Project
		err := json.Unmarshal(kv.Value, &project)
		if err != nil {
			return nil, err
		}

		for _, r := range project.Repositories {
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

// GetOwnedProjects returns owner projects of the repository from etcd
func (d *model) GetOwnedProjects(ctx context.Context, repo string, branch string) ([]livegreptone.Project, error) {
	resp, err := d.etcd.Get(ctx, ProjectKeyPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	var projects []livegreptone.Project
	for _, kv := range resp.Kvs {
		var project livegreptone.Project
		err := json.Unmarshal(kv.Value, &project)
		if err != nil {
			return nil, err
		}

		for _, r := range project.Repositories {
			if r.URL == repo && r.Branch == branch {
				projects = append(projects, project)
			}
		}
	}
	return projects, nil
}
