package db

import (
	"context"
	"errors"

	"github.com/ueokande/livegreptone"
)

var (
	// ErrNotFound represents a not found error
	ErrNotFound = errors.New("not found")
)

// Interface is an interface to access database
type Interface interface {
	// GetProject returns a project of id
	GetProject(ctx context.Context, id string) (livegreptone.Project, error)

	// GetAllProjects returns all projects
	GetAllProjects(ctx context.Context) ([]livegreptone.Project, error)

	// GetProjectIDs returns project IDs
	GetProjectIDs(ctx context.Context) ([]string, error)

	// UpdateProject create or update a project of id
	UpdateProject(ctx context.Context, project livegreptone.Project) error

	// GetRepositories removes a repository
	GetRepositories(ctx context.Context) ([]livegreptone.Repository, error)

	// GetOwnedProjects returns owner projects of the repository
	GetOwnedProjects(ctx context.Context, repo string, branch string) ([]livegreptone.Project, error)

	// RemoveProject removes a project of id
	RemoveProject(ctx context.Context, id string) error

	// GetStatus returns an status of the repository
	GetStatus(ctx context.Context, repo string, branch string) (livegreptone.RepositoryStatus, error)

	// UpdateStatus creates if key is not existing or update if value is changed
	UpdateStatus(ctx context.Context, status livegreptone.RepositoryStatus) error

	// WatchStatus watches repository statuses
	WatchStatus(ctx context.Context) <-chan livegreptone.RepositoryStatus
}
