package lgrunner

import (
	"context"
	"strings"

	"github.com/docker/docker/client"
)

const (
	// Image is a docker image form livegrep
	Image = "livegrep/base:ac5078ceb8-0"
)

// Runner is an interface of the livegrep runner
type Runner interface {
	CreateIndex(ctx context.Context, manifest IndexManifest) error

	RerunIndexDB(ctx context.Context, project string) error
	RerunWeb(ctx context.Context, config WebConfig) error
}

// NewRunner creates a Runner
func NewRunner(gitRootFS string) Runner {
	return &runnerImpl{
		gitRootFS: gitRootFS,
	}
}

type runnerImpl struct {
	gitRootFS string
	docker    *client.Client
}

// IsErrNoSuchContainer returns true if err container "No such container"
func IsErrNoSuchContainer(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "No such container")
}
