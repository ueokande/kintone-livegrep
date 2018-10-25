package db

const (
	// ProjectKeyPrefix is an prefix of the projects
	ProjectKeyPrefix = "projects/"

	// RepositoryKeyPrefix is an prefix of the repository list
	RepositoryKeyPrefix = "repositories/"

	// StatusesKeyPrefix is an prefix of the status
	StatusesKeyPrefix = "statuses/"
)

// ProjectKey returns an etcd key of id
func ProjectKey(id string) string {
	return ProjectKeyPrefix + id
}

// RepositoryKey returns an etcd key of repo and branch
func RepositoryKey(repo, branch string) string {
	return RepositoryKeyPrefix + repo + "/" + branch
}

// StatusKey returns an etcd key of repo and branch
func StatusKey(repo, branch string) string {
	return StatusesKeyPrefix + repo + "/" + branch
}
