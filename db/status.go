package db

import (
	"context"
	"encoding/json"

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
func (d *model) UpdateStatus(ctx context.Context, repo string, branch string, status livegreptone.RepositoryStatus) error {
	key := StatusKey(repo, branch)
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
