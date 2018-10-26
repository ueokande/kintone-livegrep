package lgrunner

import (
	"context"
	"io/ioutil"
	"log"

	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	"github.com/ueokande/livegreptone/db"
)

// Server is a livegrep runner server
type Server struct {
	GitRootFS string
	DB        db.Interface
}

// Run watches etcd and runs a livegrep containers
func (s *Server) Run(ctx context.Context) error {
	d, err := docker.NewEnvClient()
	if err != nil {
		return err
	}
	r, err := d.ImagePull(ctx, Image, types.ImagePullOptions{})
	if err != nil {
		log.Printf("failed to pull image: %v", err)
		return err
	}
	out, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	log.Printf("Pulled image: %v", string(out))

	runner := runnerImpl{
		gitRootFS: s.GitRootFS,
		docker:    d,
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
			err = runner.RerunIndexDB(ctx, p.Name)
			if err != nil {
				return err
			}
		}

		projects, err = s.DB.GetAllProjects(ctx)
		if err != nil {
			return err
		}

		config := WebConfigFromProjects(projects)
		err = runner.RerunWeb(ctx, config)
		if err != nil {
			return err
		}
	}

	return nil
}
