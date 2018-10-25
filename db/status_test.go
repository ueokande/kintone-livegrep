package db

import (
	"context"
	"testing"

	"github.com/ueokande/livegreptone"
)

func testGetAndUpdate(t *testing.T) {
	ctx := context.Background()

	etcd, err := NewEtcdClient(t)
	if err != nil {
		t.Fatal(err)
	}
	m := New(etcd)
	err = m.UpdateStatus(ctx, "https://github.com/kubernetes/kubernetes", "master",
		livegreptone.RepositoryStatus{Commit: "1234ABCD"})
	if err != nil {
		t.Fatal(err)
	}
	err = m.UpdateStatus(ctx, "https://github.com/kubernetes/kubernetes", "release-1.11",
		livegreptone.RepositoryStatus{Commit: "1000AAAA"})
	if err != nil {
		t.Fatal(err)
	}

	s, err := m.GetStatus(ctx, "https://github.com/kubernetes/kubernetes", "master")
	if err != nil {
		t.Fatal(err)
	}
	if s.Commit != "1234ABCD" {
		t.Errorf(`status.Commit != "1234ABCD": %v`, s.Commit)
	}

	s, err = m.GetStatus(ctx, "https://github.com/kubernetes/kubernetes", "release-1.11")
	if err != nil {
		t.Fatal(err)
	}
	if s.Commit != "1000AAAA" {
		t.Errorf(`status.Commit != "1000AAAA": %v`, s.Commit)
	}
}

func testUpdateIfNoUpdates(t *testing.T) {
	ctx := context.Background()
	etcd, err := NewEtcdClient(t)
	if err != nil {
		t.Fatal(err)
	}
	m := New(etcd)
	err = m.UpdateStatus(ctx, "https://github.com/kubernetes/kubernetes", "master",
		livegreptone.RepositoryStatus{Commit: "1234ABCD"})
	if err != nil {
		t.Fatal(err)
	}

	err = m.UpdateStatus(ctx, "https://github.com/kubernetes/kubernetes", "master",
		livegreptone.RepositoryStatus{Commit: "1000AAAA"})
	if err != nil {
		t.Fatal(err)
	}
	s, err := m.GetStatus(ctx, "https://github.com/kubernetes/kubernetes", "master")
	if err != nil {
		t.Fatal(err)
	}
	if s.Commit != "1000AAAA" {
		t.Errorf(`status.Commit != "1000AAAA": %v`, s.Commit)
	}

	resp1, err := etcd.Get(ctx, "/")
	if err != nil {
		t.Fatal(err)
	}
	err = m.UpdateStatus(ctx, "https://github.com/kubernetes/kubernetes", "master",
		livegreptone.RepositoryStatus{Commit: "1000AAAA"})
	if err != nil {
		t.Fatal(err)
	}
	resp2, err := etcd.Get(ctx, "/")
	if err != nil {
		t.Fatal(err)
	}
	if resp1.Header.Revision != resp2.Header.Revision {
		t.Error("resp1.Header.Revision != resp2.Header.Revision")
	}

}

func TestStatus(t *testing.T) {
	t.Run("GetAndUpdate", testGetAndUpdate)
	t.Run("testUpdateIfNoUpdates", testUpdateIfNoUpdates)
}
