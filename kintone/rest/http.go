package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ueokande/livegreptone/kintone"
)

type restHTTPClient struct {
	http   *http.Client
	origin string
	appID  int
	token  string
}

// NewClient creates an implementation of the Kintone Client
func NewClient(http *http.Client, origin string, appID int, token string) Interface {
	return &restHTTPClient{
		http:   http,
		origin: origin,
		appID:  appID,
		token:  token,
	}
}

// GetRecords gets a record
func (c *restHTTPClient) GetRecord(ctx context.Context, id int) (*kintone.Record, error) {
	url := getRecordURL(c.origin, c.appID, id)
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("X-Cybozu-API-Token", c.token)
	r.WithContext(ctx)

	resp, err := c.http.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data kintone.RecordResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data.Record, nil
}

// GetRecords gets record list
func (c *restHTTPClient) GetRecords(ctx context.Context) ([]kintone.Record, error) {
	// TODO support records more than 100 records
	url := getRecordsURL(c.origin, c.appID)
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("X-Cybozu-API-Token", c.token)
	r.WithContext(ctx)

	resp, err := c.http.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data kintone.RecordListResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data.Records, nil
}
