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
	kintoneOrigin = os.Getenv("KINTONE_ORIGIN")
	kintoneAppID  = os.Getenv("KINTONE_APP_ID")
	kintoneToken  = os.Getenv("KINTONE_API_TOKEN")
	etcdEndpoints = os.Getenv("ETCD_ENDPOINTS")
	etcdPrefix    = os.Getenv("ETCD_PREFIX")
)

func main() {
	if len(kintoneOrigin) == 0 {
		log.Fatal(errors.New("KINTONE_ORIGIN is not set"))
	}
	if len(kintoneAppID) == 0 {
		log.Fatal(errors.New("KINTONE_APP_ID is not set"))
	}
	if len(kintoneToken) == 0 {
		log.Fatal(errors.New("KINTONE_API_TOKEN is not set"))
	}
	if len(etcdEndpoints) == 0 {
		log.Fatal(errors.New("ETCD_ENDPOINTS is not set"))
	}

	appID, err := strconv.Atoi(kintoneAppID)
	if err != nil {
		log.Fatalf("invalid KINTONE_APP_ID: %v", err)
	}

	k := rest.NewClient(&http.Client{}, kintoneOrigin, appID, kintoneToken)
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
