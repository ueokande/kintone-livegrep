package rest

import "fmt"

func getRecordURL(origin string, app int, record int) string {
	return fmt.Sprintf("%s/k/v1/record.json?app=%d&id=%d", origin, app, record)
}

func getRecordsURL(origin string, app int) string {
	return fmt.Sprintf("%s/k/v1/records.json?app=%d", origin, app)
}
