package db

import "github.com/coreos/etcd/clientv3"

type model struct {
	etcd *clientv3.Client
}

func New(etcd *clientv3.Client) Interface {
	return &model{etcd}
}
