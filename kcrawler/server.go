package kcrawler

import (
	"context"

	"github.com/ueokande/livegreptone"
	"github.com/ueokande/livegreptone/db"
	"github.com/ueokande/livegreptone/kintone/rest"
)

// Server is an kcrawler server
type Server struct {
	// Kintone is an interface of the Kintone REST client
	Kintone rest.Interface

	// DB is an interface of the database
	DB db.Interface
}

// Run sync project informations from Kintone
func (s *Server) Run(ctx context.Context) error {
	records, err := s.Kintone.GetRecords(ctx)
	if err != nil {
		return err
	}

	projects := make([]livegreptone.Project, len(records))
	for i, r := range records {
		repos := make([]livegreptone.Repository, len(r.Repositories.Value))
		for j, r := range r.Repositories.Value {
			repos[j] = livegreptone.Repository{
				URL:    r.Value.URL.Value,
				Branch: r.Value.Branch.Value,
			}
		}
		projects[i].ID = r.ID.Value
		projects[i].Name = r.Name.Value
		projects[i].Repositories = repos
		projects[i].Revision = r.Revision.Value
	}

	err = s.cleanProjects(ctx, projects)
	if err != nil {
		return err
	}
	err = s.addProjects(ctx, projects)
	return nil
}

func (s *Server) cleanProjects(ctx context.Context, projects []livegreptone.Project) error {
	ids, err := s.DB.GetProjectIDs(ctx)
	if err != nil {
		return err
	}
	idsSet := make(map[string]struct{})
	for _, id := range ids {
		idsSet[id] = struct{}{}
	}
	for _, p := range projects {
		delete(idsSet, p.ID)
	}
	for id := range idsSet {
		err := s.DB.RemoveProject(ctx, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) addProjects(ctx context.Context, projects []livegreptone.Project) error {
	for _, p := range projects {
		err := s.DB.UpdateProject(ctx, p)
		if err != nil {
			return err
		}
	}
	return nil
}
