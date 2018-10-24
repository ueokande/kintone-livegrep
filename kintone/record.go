package kintone

type SingleLineTextField struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type RevisionField struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type IdField struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type RepositoriesField struct {
	Type  string          `json:"type"`
	Value []RepositoryRow `json:"value"`
}

type RepositoryRow struct {
	Id    string `json:"id"`
	Value struct {
		URL    SingleLineTextField `json:"url"`
		Branch SingleLineTextField `json:"branch"`
	} `json:"value"`
}

type Record struct {
	Repositories RepositoriesField   `json:"Repositories"`
	Name         SingleLineTextField `json:"name"`
	Revision     RevisionField       `json:"$revision"`
	Id           IdField             `json:"$id"`
}

type RecordResponse struct {
	Record Record `json:"record"`
}

type RecordListResponse struct {
	Records []Record `json:"records"`
}
