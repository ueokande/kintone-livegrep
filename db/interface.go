package db

import (
	"context"
	"errors"

	"github.com/ueokande/livegreptone"
)

var (
	ErrNotFound = errors.New("not found")
)

type Interface interface {
	GetProject(ctx context.Context, id string) (livegreptone.Project, error)

	GetProjectIds(ctx context.Context) ([]string, error)

	UpdateProject(ctx context.Context, project livegreptone.Project) error

	RemoveProject(ctx context.Context, id string) error

	GetStatus(ctx context.Context, repo string, branch string) (livegreptone.RepositoryStatus, error)

	UpdateStatus(ctx context.Context, repo string, branch string, status livegreptone.RepositoryStatus) error
}
