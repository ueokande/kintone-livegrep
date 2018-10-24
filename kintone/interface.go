package kintone

import (
	"context"

	"github.com/ueokande/livegreptone"
)

type ProjectsModel interface {
	DumpProjects(ctx context.Context) ([]livegreptone.Project, error)
}
