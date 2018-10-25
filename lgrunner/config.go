package lgrunner

import (
	"errors"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/ueokande/livegreptone"
)

type IndexManifest struct {
	Name    string   `json:"name"`
	FSPaths []FSPath `json:"fs_paths"`
}

type FSPath struct {
	Name     string            `json:"name"`
	Path     string            `json:"path"`
	Metadata map[string]string `json:"metadata"`
}

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

type WebConfig struct {
	Backends []BackendConfig `json"backends"`
	Listen   string          `json:"listen"`
}

type BackendConfig struct {
	ID      string `json:"id"`
	Address string `json:"addr"`
}
