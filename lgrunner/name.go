package lgrunner

const (
	IndexVolumeName = "livegrep-index"

	WebContainerName = "livegrep-web"
	WebHostName      = "livegrep-web"

	IndexContainerNamePrefix = "livegrep-index-"
	IndexHostNamePrefix      = "livegrep-index-"
)

func IndexContainerName(project string) string {
	return IndexContainerNamePrefix + project
}

func IndexHostName(project string) string {
	return IndexHostNamePrefix + project
}
