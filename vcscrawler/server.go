package vcscrawler

import (
	"context"

	"github.com/ueokande/livegreptone"
	"github.com/ueokande/livegreptone/db"
	"github.com/ueokande/livegreptone/git"
)

// Server is an vcs crawl server
type Server struct {
	// Git is an interface of the git engine
	Git git.Interface

	// DB is an interface of the database
	DB db.Interface
}

// Run runs vcs crawler
func (s *Server) Run(ctx context.Context) error {
	repos, err := s.DB.GetRepositories(ctx)
	if err != nil {
		return err
	}
	for _, r := range repos {
		commit, err := s.Git.Update(ctx, r.URL, r.Branch)
		if err != nil {
			return err
		}
		st := livegreptone.RepositoryStatus{Commit: commit}
		err = s.DB.UpdateStatus(ctx, r.URL, r.Branch, st)
		if err != nil {
			return err
		}

	}
	return nil
}
