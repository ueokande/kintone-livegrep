package lgrunner

import (
	"context"
)

const (
	// Image is a docker image form livegrep
	Image = "livegrep/base:ac5078ceb8-0"
)

// Runner is an interface of the livegrep runner
type Runner interface {
	CreateIndex(ctx context.Context, manifest IndexManifest) error

	RunIndexDB(ctx context.Context, project string) error
	StopIndexDB(ctx context.Context, project string) error

	RunWeb(ctx context.Context, config WebConfig) error
	StopWeb(ctx context.Context) error
}

// NewRunner creates a Runner
func NewRunner(gitRootFS string) Runner {
	return &runnerImpl{
		gitRootFS: gitRootFS,
	}
}

type runnerImpl struct {
	gitRootFS string
}
