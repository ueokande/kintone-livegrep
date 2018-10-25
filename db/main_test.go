package db

import (
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/namespace"
)

func NewEtcdClient(t *testing.T) (*clientv3.Client, error) {
	cfg := clientv3.Config{
		Endpoints:   []string{"http://127.0.0.1:2379"},
		DialTimeout: 2 * time.Second,
	}

	c, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}
	prefix := t.Name()
	c.KV = namespace.NewKV(c.KV, prefix)
	c.Watcher = namespace.NewWatcher(c.Watcher, prefix)
	c.Lease = namespace.NewLease(c.Lease, prefix)

	return c, nil
}
