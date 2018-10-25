package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/namespace"
	"github.com/cybozu-go/well"
	"github.com/ueokande/livegreptone/db"
	"github.com/ueokande/livegreptone/kcrawler"
	"github.com/ueokande/livegreptone/kintone/rest"
)

var (
	KintoneOrigin = os.Getenv("KINTONE_ORIGIN")
	KintoneAppID  = os.Getenv("KINTONE_APP_ID")
	KintoneToken  = os.Getenv("KINTONE_API_TOKEN")
	EtcdEndpoints = os.Getenv("ETCD_ENDPOINTS")
	EtcdPrefix    = os.Getenv("ETCD_PREFIX")
)

func main() {
	if len(KintoneOrigin) == 0 {
		log.Fatal(errors.New("KINTONE_ORIGIN is not set"))
	}
	if len(KintoneAppID) == 0 {
		log.Fatal(errors.New("KINTONE_APP_ID is not set"))
	}
	if len(KintoneToken) == 0 {
		log.Fatal(errors.New("KINTONE_API_TOKEN is not set"))
	}
	if len(EtcdEndpoints) == 0 {
		log.Fatal(errors.New("ETCD_ENDPOINTS is not set"))
	}

	appID, err := strconv.Atoi(KintoneAppID)
	if err != nil {
		log.Fatalf("invalid KINTONE_APP_ID: %v", err)
	}

	k := rest.NewClient(&http.Client{}, KintoneOrigin, appID, KintoneToken)
	etcdcfg := clientv3.Config{
		Endpoints:   []string{"http://127.0.0.1:2379"},
		DialTimeout: 2 * time.Second,
	}

	etcd, err := clientv3.New(etcdcfg)
	if err != nil {
		log.Fatalf("Failed to launch etcd: %v", err)
	}
	if len(EtcdPrefix) > 0 {
		etcd.KV = namespace.NewKV(etcd.KV, EtcdPrefix)
		etcd.Watcher = namespace.NewWatcher(etcd.Watcher, EtcdPrefix)
		etcd.Lease = namespace.NewLease(etcd.Lease, EtcdPrefix)
	}
	db := db.New(etcd)
	s := kcrawler.Server{
		Kintone: k,
		DB:      db,
	}

	well.Go(s.Run)
	well.Stop()
	err = well.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
