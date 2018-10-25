package kcrawler

import (
	"context"

	"github.com/ueokande/livegreptone"
	"github.com/ueokande/livegreptone/db"
	"github.com/ueokande/livegreptone/kintone/rest"
)

type Server struct {
	KintoneOrigin string
	KintoneAppId  int
	KintoneToken  string

	Kintone rest.Interface
	DB      db.Interface
}

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
		projects[i].Id = r.Id.Value
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
	ids, err := s.DB.GetProjectIds(ctx)
	if err != nil {
		return err
	}
	idsSet := make(map[string]struct{})
	for _, id := range ids {
		idsSet[id] = struct{}{}
	}
	for _, p := range projects {
		delete(idsSet, p.Id)
	}
	for id, _ := range idsSet {
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
