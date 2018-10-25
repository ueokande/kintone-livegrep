package db

import (
	"context"
	"reflect"
	"sort"
	"testing"

	"github.com/ueokande/livegreptone"
)

func testUpdateGet(t *testing.T) {
	ctx := context.Background()

	p1 := livegreptone.Project{ID: "0", Name: "Kubernetes"}
	p2 := livegreptone.Project{ID: "1", Name: "Ceph"}
	p3 := livegreptone.Project{ID: "2", Name: "Mesos"}

	etcd, err := NewEtcdClient(t)
	if err != nil {
		t.Fatal(err)
	}
	m := New(etcd)
	for _, p := range []livegreptone.Project{p1, p2, p3} {
		err := m.UpdateProject(ctx, p)
		if err != nil {
			t.Fatal(err)
		}
	}

	p, err := m.GetProject(ctx, "1")
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "Ceph" {
		t.Errorf(`p.Name != "Ceph": %v`, p.Name)
	}

	err = m.RemoveProject(ctx, "1")
	if err != nil {
		t.Fatal(err)
	}
}

func testGetIDs(t *testing.T) {
	ctx := context.Background()

	p1 := livegreptone.Project{ID: "0", Name: "Kubernetes"}
	p2 := livegreptone.Project{ID: "1", Name: "Ceph"}
	p3 := livegreptone.Project{ID: "2", Name: "Mesos"}

	etcd, err := NewEtcdClient(t)
	if err != nil {
		t.Fatal(err)
	}
	m := New(etcd)
	for _, p := range []livegreptone.Project{p1, p2, p3} {
		err := m.UpdateProject(ctx, p)
		if err != nil {
			t.Fatal(err)
		}
	}

	ids, err := m.GetProjectIDs(ctx)
	sort.Strings(ids)
	if !reflect.DeepEqual(ids, []string{"0", "1", "2"}) {
		t.Errorf(`reflect.DeepEqual(ids, []string{"0", "2"}): %v`, ids)
	}
}

func testGetRepositories(t *testing.T) {
	ctx := context.Background()

	p1 := livegreptone.Project{ID: "0", Name: "Kubernetes", Repositories: []livegreptone.Repository{
		livegreptone.Repository{URL: "https://github.com/kubernetes/kubernetes", Branch: "master"},
		livegreptone.Repository{URL: "https://github.com/kubernetes/kubernetes", Branch: "release-1.11"},
	}}
	p2 := livegreptone.Project{ID: "1", Name: "Ceph"}
	p3 := livegreptone.Project{ID: "2", Name: "Mesos", Repositories: []livegreptone.Repository{
		livegreptone.Repository{URL: "https://github.com/apache/mesos", Branch: "master"},
	}}

	etcd, err := NewEtcdClient(t)
	if err != nil {
		t.Fatal(err)
	}
	m := New(etcd)
	for _, p := range []livegreptone.Project{p1, p2, p3} {
		err := m.UpdateProject(ctx, p)
		if err != nil {
			t.Fatal(err)
		}
	}

	repos, err := m.GetRepositories(ctx)
	if len(repos) != 3 {
		t.Fatalf("len(repos) != 3: %v", len(repos))
	}
}

func testGetOwnedProjects(t *testing.T) {
	ctx := context.Background()

	p1 := livegreptone.Project{ID: "0", Name: "Kubernetes1", Repositories: []livegreptone.Repository{
		livegreptone.Repository{URL: "https://github.com/kubernetes/kubernetes", Branch: "master"},
	}}
	p2 := livegreptone.Project{ID: "1", Name: "Kubernetes2", Repositories: []livegreptone.Repository{
		livegreptone.Repository{URL: "https://github.com/kubernetes/kubernetes", Branch: "release-1.11"},
	}}
	p3 := livegreptone.Project{ID: "2", Name: "Kubernetes3", Repositories: []livegreptone.Repository{
		livegreptone.Repository{URL: "https://github.com/kubernetes/kubernetes", Branch: "master"},
	}}

	etcd, err := NewEtcdClient(t)
	if err != nil {
		t.Fatal(err)
	}
	m := New(etcd)
	for _, p := range []livegreptone.Project{p1, p2, p3} {
		err := m.UpdateProject(ctx, p)
		if err != nil {
			t.Fatal(err)
		}
	}

	repos, err := m.GetOwnedProjects(ctx, "https://github.com/kubernetes/kubernetes", "master")
	if len(repos) != 2 {
		t.Fatalf("len(repos) != 3: %v", len(repos))
	}
	names := []string{repos[0].Name, repos[1].Name}
	sort.Strings(names)
	if !reflect.DeepEqual(names, []string{"Kubernetes1", "Kubernetes3"}) {
		t.Errorf(`reflect.DeepEqual(names, []string{"kubernetes1", "Kubernetes3"}): %v`, names)
	}
}

func TestProject(t *testing.T) {
	t.Run("UpdateGet", testUpdateGet)
	t.Run("GetIDs", testGetIDs)
	t.Run("GetRepositories", testGetRepositories)
	t.Run("GetOwnedProjects", testGetOwnedProjects)
}
