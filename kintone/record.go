package kintone

// SingleLineTextField represents SINGLE_LINE_TEXT field
type SingleLineTextField struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// RevisionField represents __REVISION__ field
type RevisionField struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// IDField represents __ID__ field
type IDField struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// RepositoriesField represents SUBTABLE field of the Repositories
type RepositoriesField struct {
	Type  string          `json:"type"`
	Value []RepositoryRow `json:"value"`
}

// RepositoryRow represents a row of the Repositories
type RepositoryRow struct {
	ID    string `json:"id"`
	Value struct {
		URL    SingleLineTextField `json:"url"`
		Branch SingleLineTextField `json:"branch"`
	} `json:"value"`
}

// Record represents App's record
type Record struct {
	Repositories RepositoriesField   `json:"Repositories"`
	Name         SingleLineTextField `json:"name"`
	Revision     RevisionField       `json:"$revision"`
	ID           IDField             `json:"$id"`
}

// RecordResponse represents a response of the recpord from REST API
type RecordResponse struct {
	Record Record `json:"record"`
}

// RecordListResponse represents a response of records from REST API
type RecordListResponse struct {
	Records []Record `json:"records"`
}
