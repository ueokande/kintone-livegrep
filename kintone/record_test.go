package kintone

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestRecord(t *testing.T) {
	path := filepath.Join("testdata", "record.json")
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	var resp RecordResponse
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		t.Fatal(err)
	}
	record := resp.Record
	if record.Name.Value != "Kubernetes" {
		t.Errorf("record.Name.Value != \"Kubernetes\"; %v", record.Name.Value)
	}
	if record.ID.Value != "1" {
		t.Errorf("record.ID.Value != \"1\"; %v", record.ID.Value)
	}
	if record.Revision.Value != "2" {
		t.Errorf("record.Revision.Value != \"1\"; %v", record.Revision.Value)
	}
	if len(record.Repositories.Value) != 1 {
		t.Fatalf("len(record.Repositories.Value) != 1: %v", len(record.Repositories.Value))
	}
	repo := record.Repositories.Value[0]
	if repo.Value.URL.Value != "https://github.com/kubernetes/kubernetes" {
		t.Errorf("repo.Value.URL.Value != \"https://github.com/kubernetes/kubernetes\": %v", repo.Value.URL.Value)
	}
	if repo.Value.Branch.Value != "master" {
		t.Errorf("repo.Value.Branch.Value != \"master\": %v", repo.Value.Branch.Value)
	}
}

func TestRecords(t *testing.T) {
	path := filepath.Join("testdata", "records.json")
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	var resp RecordListResponse
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		t.Fatal(err)
	}
	records := resp.Records
	if len(records) != 1 {
		t.Fatalf("len(records) != 1: %v", len(records))
	}
	record := records[0]
	if record.Name.Value != "Kubernetes" {
		t.Errorf("record.Name.Value != \"Kubernetes\"; %v", record.Name.Value)
	}
	if record.ID.Value != "1" {
		t.Errorf("record.ID.Value != \"1\"; %v", record.ID.Value)
	}
	if record.Revision.Value != "2" {
		t.Errorf("record.Revision.Value != \"1\"; %v", record.Revision.Value)
	}
	if len(record.Repositories.Value) != 1 {
		t.Fatalf("len(record.Repositories.Value) != 1: %v", len(record.Repositories.Value))
	}
	repo := record.Repositories.Value[0]
	if repo.Value.URL.Value != "https://github.com/kubernetes/kubernetes" {
		t.Errorf("repo.Value.URL.Value != \"https://github.com/kubernetes/kubernetes\": %v", repo.Value.URL.Value)
	}
	if repo.Value.Branch.Value != "master" {
		t.Errorf("repo.Value.Branch.Value != \"master\": %v", repo.Value.Branch.Value)
	}
}
