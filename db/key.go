package db

const (
	// ProjectKeyPrefix is an prefix of the projects
	ProjectKeyPrefix = "projects/"

	// StatusesKeyPrefix is an prefix of the status
	StatusesKeyPrefix = "statuses/"
)

// ProjectKey returns an etcd key of id
func ProjectKey(id string) string {
	return ProjectKeyPrefix + id
}

// StatusKey returns an etcd key of repo and branch
func StatusKey(repo, branch string) string {
	return StatusesKeyPrefix + repo + "/" + branch
}
