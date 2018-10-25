package kcrawler

import (
	"context"
	"reflect"
	"sort"
	"testing"

	"github.com/ueokande/livegreptone"
	db "github.com/ueokande/livegreptone/db/mock"
	"github.com/ueokande/livegreptone/kintone"
	rest "github.com/ueokande/livegreptone/kintone/rest/mock"
)

func testCleanProjects(t *testing.T) {
	ctx := context.Background()

	p1 := livegreptone.Project{ID: "0"}
	p2 := livegreptone.Project{ID: "1"}
	p3 := livegreptone.Project{ID: "2"}

	db := db.New()
	db.UpdateProject(ctx, p1)
	db.UpdateProject(ctx, p2)
	db.UpdateProject(ctx, p3)

	s := Server{DB: db}
	err := s.cleanProjects(ctx, []livegreptone.Project{p2, p3})
	if err != nil {
		t.Fatal(err)
	}
	ids, _ := db.GetProjectIDs(ctx)
	sort.Strings(ids)
	if !reflect.DeepEqual(ids, []string{"1", "2"}) {
		t.Errorf(`reflect.DeepEqual(ids, []string{"1", "2"}): %v`, ids)
	}

}

func testAddProjects(t *testing.T) {
	ctx := context.Background()

	db := db.New()
	db.UpdateProject(ctx, livegreptone.Project{ID: "0", Name: "Kubernetes"})
	db.UpdateProject(ctx, livegreptone.Project{ID: "1", Name: "Ceph"})

	s := Server{DB: db}
	err := s.addProjects(ctx, []livegreptone.Project{
		livegreptone.Project{ID: "0", Name: "K8s"},
		livegreptone.Project{ID: "1", Name: "Ceph"},
		livegreptone.Project{ID: "2", Name: "Mesos"},
	})
	if err != nil {
		t.Fatal(err)
	}

	if p, _ := db.GetProject(ctx, "0"); p.Name != "K8s" {
		t.Errorf(`p.Name != "K8s": %v`, p.Name)
	}
	if p, _ := db.GetProject(ctx, "1"); p.Name != "Ceph" {
		t.Errorf(`p.Name != "Ceph": %v`, p.Name)
	}
	if p, _ := db.GetProject(ctx, "2"); p.Name != "Mesos" {
		t.Errorf(`p.Name != "Mesos": %v`, p.Name)
	}
}

func testRun(t *testing.T) {
	ctx := context.Background()

	db := db.New()
	db.UpdateProject(ctx, livegreptone.Project{ID: "0"})
	db.UpdateProject(ctx, livegreptone.Project{ID: "1"})

	rest := rest.New([]kintone.Record{
		kintone.Record{ID: kintone.IDField{Value: "2"}},
		kintone.Record{ID: kintone.IDField{Value: "3"}},
	})

	s := Server{DB: db, Kintone: rest}
	err := s.Run(ctx)
	if err != nil {
		t.Fatal(err)
	}

	ids, _ := db.GetProjectIDs(ctx)
	sort.Strings(ids)
	if !reflect.DeepEqual(ids, []string{"2", "3"}) {
		t.Errorf(`reflect.DeepEqual(ids, []string{"2", "3"}): %v`, ids)
	}
}

func TestServer(t *testing.T) {
	t.Run("CleanProjects", testCleanProjects)
	t.Run("AddProjects", testAddProjects)
	t.Run("Run", testRun)
}
