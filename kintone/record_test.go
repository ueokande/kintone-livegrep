package kintone

import (
	"encoding/json"
	"testing"
)

func TestRecord(t *testing.T) {
	body := `{
  "record": {
    "Updated_datetime": {
      "type": "UPDATED_TIME",
      "value": "2018-10-24T01:35:00Z"
    },
    "Created_datetime": {
      "type": "CREATED_TIME",
      "value": "2018-10-24T01:35:00Z"
    },
    "Repositories": {
      "type": "SUBTABLE",
      "value": [
        {
          "id": "375434423",
          "value": {
            "branch": {
              "type": "SINGLE_LINE_TEXT",
              "value": "master"
            },
            "url": {
              "type": "SINGLE_LINE_TEXT",
              "value": "https://github.com/kubernetes/kubernetes"
            }
          }
        }
      ]
    },
    "Record_number": {
      "type": "RECORD_NUMBER",
      "value": "1"
    },
    "name": {
      "type": "SINGLE_LINE_TEXT",
      "value": "Kubernetes"
    },
    "Created_by": {
      "type": "CREATOR",
      "value": {
        "code": "ueokande",
        "name": "Shin'ya Ueoka"
      }
    },
    "$revision": {
      "type": "__REVISION__",
      "value": "2"
    },
    "Updated_by": {
      "type": "MODIFIER",
      "value": {
        "code": "ueokande",
        "name": "Shin'ya Ueoka"
      }
    },
    "$id": {
      "type": "__ID__",
      "value": "1"
    }
  }
}`

	var resp RecordResponse
	err := json.Unmarshal([]byte(body), &resp)
	if err != nil {
		t.Fatal(err)
	}
	record := resp.Record
	if record.Name.Value != "Kubernetes" {
		t.Errorf("record.Name.Value != \"Kubernetes\"; %v", record.Name.Value)
	}
	if record.Id.Value != "1" {
		t.Errorf("record.Id.Value != \"1\"; %v", record.Id.Value)
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
	body := `{
  "records": [
    {
      "Updated_datetime": {
        "type": "UPDATED_TIME",
        "value": "2018-10-24T01:35:00Z"
      },
      "Created_datetime": {
        "type": "CREATED_TIME",
        "value": "2018-10-24T01:35:00Z"
      },
      "Repositories": {
        "type": "SUBTABLE",
        "value": [
          {
            "id": "375434423",
            "value": {
              "branch": {
                "type": "SINGLE_LINE_TEXT",
                "value": "master"
              },
              "url": {
                "type": "SINGLE_LINE_TEXT",
                "value": "https://github.com/kubernetes/kubernetes"
              }
            }
          }
        ]
      },
      "Record_number": {
        "type": "RECORD_NUMBER",
        "value": "1"
      },
      "name": {
        "type": "SINGLE_LINE_TEXT",
        "value": "Kubernetes"
      },
      "Created_by": {
        "type": "CREATOR",
        "value": {
          "code": "uekande",
          "name": "Shin'ya Ueoka"
        }
      },
      "$revision": {
        "type": "__REVISION__",
        "value": "2"
      },
      "Updated_by": {
        "type": "MODIFIER",
        "value": {
          "code": "uekande",
          "name": "Shin'ya Ueoka"
        }
      },
      "$id": {
        "type": "__ID__",
        "value": "1"
      }
    }
  ],
  "totalCount": null
}`

	var resp RecordListResponse
	err := json.Unmarshal([]byte(body), &resp)
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
	if record.Id.Value != "1" {
		t.Errorf("record.Id.Value != \"1\"; %v", record.Id.Value)
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
