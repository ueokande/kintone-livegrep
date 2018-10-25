package lgrunner

import (
	"context"

	"github.com/ueokande/livegreptone/db"
)

// Server is a livegrep runner server
type Server struct {
	GitRootFS string
	DB        db.Interface
}

// Run watches etcd and runs a livegrep containers
func (s *Server) Run(ctx context.Context) error {
	runner := runnerImpl{
		gitRootFS: s.GitRootFS,
	}
	for st := range s.DB.WatchStatus(ctx) {
		projects, err := s.DB.GetOwnedProjects(ctx, st.URL, st.Branch)
		if err != nil {
			return err
		}
		for _, p := range projects {
			manifest := ManifestFromProject(p)
			err := runner.CreateIndex(ctx, manifest)
			if err != nil {
				return err
			}
			// Ignore stopping index db
			runner.StopIndexDB(ctx, p.Name)
			err = runner.RunIndexDB(ctx, p.Name)
			if err != nil {
				return err
			}
		}

		projects, err = s.DB.GetAllProjects(ctx)
		if err != nil {
			return err
		}

		// Ignore stopping web
		runner.StopWeb(ctx)

		config := WebConfigFromProjects(projects)
		err = runner.RunWeb(ctx, config)
		if err != nil {
			return err
		}
	}

	return nil
}
