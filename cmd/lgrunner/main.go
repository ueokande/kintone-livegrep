package main

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/namespace"
	"github.com/cybozu-go/well"
	"github.com/ueokande/livegreptone/db"
	"github.com/ueokande/livegreptone/lgrunner"
)

var (
	etcdEndpoints = os.Getenv("ETCD_ENDPOINTS")
	etcdPrefix    = os.Getenv("ETCD_PREFIX")
	gitRootFS     = os.Getenv("GIT_ROOT_FS")
)

func main() {
	if len(etcdEndpoints) == 0 {
		log.Fatal(errors.New("ETCD_ENDPOINTS is not set"))
	}
	if len(gitRootFS) == 0 {
		log.Fatal(errors.New("GIT_ROOT_FS is not set"))
	}

	etcdcfg := clientv3.Config{
		Endpoints:   []string{"http://127.0.0.1:2379"},
		DialTimeout: 2 * time.Second,
	}
	etcd, err := clientv3.New(etcdcfg)
	if err != nil {
		log.Fatalf("Failed to launch etcd: %v", err)
	}
	if len(etcdPrefix) > 0 {
		etcd.KV = namespace.NewKV(etcd.KV, etcdPrefix)
		etcd.Watcher = namespace.NewWatcher(etcd.Watcher, etcdPrefix)
		etcd.Lease = namespace.NewLease(etcd.Lease, etcdPrefix)
	}
	db := db.New(etcd)
	s := lgrunner.Server{
		GitRootFS: gitRootFS,
		DB:        db,
	}
	well.Go(s.Run)
	well.Stop()
	err = well.Wait()
	if err != nil {
		log.Fatal(err)
	}

}
