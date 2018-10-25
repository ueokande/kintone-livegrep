package db

import (
	"context"
	"reflect"
	"sort"
	"testing"

	"github.com/ueokande/livegreptone"
)

func TestProject(t *testing.T) {
	ctx := context.Background()

	p1 := livegreptone.Project{Id: "0", Name: "Kubernetes"}
	p2 := livegreptone.Project{Id: "1", Name: "Ceph"}
	p3 := livegreptone.Project{Id: "2", Name: "Mesos"}

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
	ids, err := m.GetProjectIds(ctx)
	sort.Strings(ids)
	if !reflect.DeepEqual(ids, []string{"0", "2"}) {
		t.Errorf(`reflect.DeepEqual(ids, []string{"0", "2"}): %v`, ids)
	}
}
