package rest

import (
	"context"

	"github.com/ueokande/livegreptone/kintone"
)

// Interface is an interface of the Kintone clinet
type Interface interface {
	GetRecord(ctx context.Context, id int) (*kintone.Record, error)

	GetRecords(ctx context.Context) ([]kintone.Record, error)
}
