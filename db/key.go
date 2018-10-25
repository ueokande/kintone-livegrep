package db

const (
	ProjectKeyPrefix  = "projects/"
	StatusesKeyPrefix = "statuses/"
)

func ProjectKey(id string) string {
	return ProjectKeyPrefix + id
}
