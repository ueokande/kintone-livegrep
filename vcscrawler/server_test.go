package vcscrawler

import (
	"context"
	"testing"

	"github.com/ueokande/livegreptone"
	db "github.com/ueokande/livegreptone/db/mock"
	git "github.com/ueokande/livegreptone/git/mock"
)

func TestRun(t *testing.T) {
	ctx := context.Background()

	db := db.New()
	db.UpdateProject(ctx, livegreptone.Project{
		ID:   "0",
		Name: "Kubernetes",
		Repositories: []livegreptone.Repository{
			{URL: "https://github.com/kubernetes/kubernetes", Branch: "master"},
			{URL: "https://github.com/kubernetes/kubernetes", Branch: "release-1.11"},
		},
	})
	db.UpdateProject(ctx, livegreptone.Project{
		ID:   "1",
		Name: "Mesos",
		Repositories: []livegreptone.Repository{
			{URL: "https://github.com/apache/mesos", Branch: "master"},
		},
	})

	git := git.New()

	s := Server{
		Git: git,
		DB:  db,
	}
	err := s.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.GetStatus(ctx, "https://github.com/kubernetes/kubernetes", "master")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.GetStatus(ctx, "https://github.com/kubernetes/kubernetes", "release-1.11")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.GetStatus(ctx, "https://github.com/apache/mesos", "master")
	if err != nil {
		t.Fatal(err)
	}
}
