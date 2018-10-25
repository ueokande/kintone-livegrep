package rest

import (
	"context"

	"github.com/ueokande/livegreptone/kintone"
)

type Interface interface {
	GetRecord(ctx context.Context, id int) (*kintone.Record, error)

	GetRecords(ctx context.Context) ([]kintone.Record, error)
}
