package mock

import (
	"context"
	"errors"
	"strconv"
	"sync"

	"github.com/ueokande/livegreptone/kintone"
	"github.com/ueokande/livegreptone/kintone/rest"
)

type mock struct {
	mu      *sync.Mutex
	records []kintone.Record
}

// New creates an mock implementation of the rest.Interface
func New(records []kintone.Record) rest.Interface {
	return &mock{
		mu:      new(sync.Mutex),
		records: records,
	}
}
func (m *mock) GetRecord(ctx context.Context, id int) (*kintone.Record, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	idstr := strconv.Itoa(id)
	for _, r := range m.records {
		if r.ID.Value == idstr {
			return &r, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mock) GetRecords(ctx context.Context) ([]kintone.Record, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.records, nil
}
