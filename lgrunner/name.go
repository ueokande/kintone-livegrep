package lgrunner

const (
	// IndexVolumeName is a volume name for the index files in livegrep
	IndexVolumeName = "livegrep-index"

	// WebContainerName is a container's name of the web server
	WebContainerName = "livegrep-web"
	// WebHostName is a container's hostname of the web server
	WebHostName = "livegrep-web"

	// IndexContainerNamePrefix is a prefix in a container's name of the index server
	IndexContainerNamePrefix = "livegrep-index-"
	// IndexHostNamePrefix is a prefix in a container's hostname of the index server
	IndexHostNamePrefix = "livegrep-index-"
)

// IndexContainerName returns a container's name from the project name
func IndexContainerName(project string) string {
	return IndexContainerNamePrefix + project
}

// IndexHostName returns a container's host name from the project name
func IndexHostName(project string) string {
	return IndexHostNamePrefix + project
}
