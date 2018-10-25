package lgrunner

import (
	"errors"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/ueokande/livegreptone"
)

// IndexManifest repserents a manifest of the codesearch
type IndexManifest struct {
	Name    string   `json:"name"`
	FSPaths []FSPath `json:"fs_paths"`
}

// FSPath represents a fs_path in the manifest of the codesearch
type FSPath struct {
	Name     string            `json:"name"`
	Path     string            `json:"path"`
	Metadata map[string]string `json:"metadata"`
}

// ManifestFromProject creates a IndexManifest from Project
func ManifestFromProject(project livegreptone.Project) IndexManifest {
	var m IndexManifest
	m.Name = project.Name
	m.FSPaths = make([]FSPath, len(project.Repositories))
	for i, r := range project.Repositories {
		host, owner, repo, err := ParseGitHubURL(r.URL)
		if err != nil {
			// Ignore non github repos
			continue

		}
		m.FSPaths[i] = FSPath{
			Name: repo,
			Path: filepath.Join("/mnt/livegrep-repos", host, owner, repo, r.Branch),
			Metadata: map[string]string{
				"url-pattern": host + "/{name}/blob/HEAD/{path}#L{lno}",
			},
		}
	}
	return m
}

// ParseGitHubURL parses url of the github, and returns host, repository owner, and repository's name
func ParseGitHubURL(urlstr string) (host, owner, repo string, err error) {
	u, err := url.Parse(urlstr)
	if err != nil {
		return "", "", "", err
	}
	paths := strings.Split(u.Path, "/")
	if len(paths) < 3 {
		return "", "", "", errors.New("imvalid URL")
	}
	return u.Host, paths[1], paths[2], nil
}

// WebConfig represents a config of the livegrep web server
type WebConfig struct {
	// Backends is backend configurations
	Backends []BackendConfig `json:"backends"`
	// Listen is a comma-separated of the listen address and port
	Listen string `json:"listen"`
}

// BackendConfig represents a background config of the livegrep web server
type BackendConfig struct {
	ID      string `json:"id"`
	Address string `json:"addr"`
}

// WebConfigFromProjects creates a WebConfig from projects
func WebConfigFromProjects(projects []livegreptone.Project) WebConfig {
	backends := make([]BackendConfig, len(projects))
	for i, p := range projects {
		backends[i] = BackendConfig{
			ID:      p.Name,
			Address: IndexHostName(p.Name) + ":9999",
		}
	}
	return WebConfig{
		Backends: backends,
		Listen:   "0.0.0.0:8910",
	}
}
