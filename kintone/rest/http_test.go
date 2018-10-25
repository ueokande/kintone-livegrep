package rest

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"testing"
)

func testGetRecord(t *testing.T) {
	path := filepath.Join("..", "testdata", "record.json")
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(bytes)
	}))
	defer ts.Close()

	url, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	origin := fmt.Sprintf("%s://%s", url.Scheme, url.Host)
	c := restHTTPClient{origin: origin, http: &http.Client{}}
	record, err := c.GetRecord(context.Background(), 10)
	if err != nil {
		t.Fatal(err)
	}
	if record.Name.Value != "Kubernetes" {
		t.Errorf("record.Name.Value != \"Kubernetes\"; %v", record.Name.Value)
	}
}

func testGetRecords(t *testing.T) {
	path := filepath.Join("..", "testdata", "records.json")
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(bytes)
	}))
	defer ts.Close()

	url, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	origin := fmt.Sprintf("%s://%s", url.Scheme, url.Host)
	c := restHTTPClient{origin: origin, http: &http.Client{}}
	records, err := c.GetRecords(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 {
		t.Fatalf("len(records) != 1: %v", len(records))
	}
	if len(records[0].Repositories.Value) != 1 {
		t.Errorf("len(records[0].Repositories.Value) != 1: %v", len(records))
	}
}

func TestRSTHTTPClient(t *testing.T) {
	t.Run("GetRecord", testGetRecord)
	t.Run("GetRecords", testGetRecords)
}
