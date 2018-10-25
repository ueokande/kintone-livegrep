package db

import (
	"context"
	"encoding/json"
	"log"

	"github.com/coreos/etcd/clientv3"
	"github.com/ueokande/livegreptone"
)

// GetStatus returns an status of the repository from etcd
func (d *model) GetStatus(ctx context.Context, repo string, branch string) (livegreptone.RepositoryStatus, error) {
	key := StatusKey(repo, branch)
	resp, err := d.etcd.Get(ctx, key)
	if err != nil {
		return livegreptone.RepositoryStatus{}, err
	}
	if resp.Count == 0 {
		return livegreptone.RepositoryStatus{}, ErrNotFound
	}
	var status livegreptone.RepositoryStatus
	err = json.Unmarshal(resp.Kvs[0].Value, &status)
	if err != nil {
		return livegreptone.RepositoryStatus{}, err
	}
	return status, nil
}

// UpdateStatus creates or update status of the repository from etcd
func (d *model) UpdateStatus(ctx context.Context, status livegreptone.RepositoryStatus) error {
	key := StatusKey(status.URL, status.Branch)
	value, err := json.Marshal(status)
	if err != nil {
		return err
	}
	_, err = d.etcd.Txn(ctx).
		If(clientv3.Compare(clientv3.Value(key), "=", string(value))).
		Then().
		Else(clientv3.OpPut(key, string(value))).
		Commit()
	if err != nil {
		return err
	}
	return nil
}

func (d *model) WatchStatus(ctx context.Context) <-chan livegreptone.RepositoryStatus {
	rch := d.etcd.Watch(ctx, StatusesKeyPrefix,
		clientv3.WithPrefix(),
	)
	ch := make(chan livegreptone.RepositoryStatus)
	go func() {
		for resp := range rch {
			for _, ev := range resp.Events {
				if ev.Type != clientv3.EventTypePut {
					// TODO support delete
					continue
				}
				var s livegreptone.RepositoryStatus
				err := json.Unmarshal(ev.Kv.Value, &s)
				if err != nil {
					log.Printf("failed parse json for %s", ev.Kv.Key)
					goto exit
				}
				ch <- s
			}
		}
	exit:
		close(ch)
	}()
	return ch
}
