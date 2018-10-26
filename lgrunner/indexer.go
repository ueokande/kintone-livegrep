package lgrunner

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"
)

func (d *runnerImpl) CreateIndex(ctx context.Context, manifest IndexManifest) error {
	log.Printf("indexing %s...", manifest.Name)
	f, err := ioutil.TempFile("", "livegreptone-index")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())

	err = json.NewEncoder(f).Encode(manifest)
	if err != nil {
		log.Printf("failed to write indexer manifeest: %v", err)
		return err
	}

	cmd := []string{
		"/livegrep/bin/codesearch", "-index_only", "-dump_index", "/mnt" + IndexPath(manifest.Name), "/mnt/manifest.json",
	}
	containerConfig := &container.Config{
		Image: Image,
		Cmd:   strslice.StrSlice(cmd),
	}
	hostConfig := &container.HostConfig{
		ReadonlyRootfs: true,
		AutoRemove:     true,
		Mounts: []mount.Mount{
			{Type: mount.TypeBind, Source: d.gitRootFS, Target: "/mnt/livegrep-repos", ReadOnly: true},
			{Type: mount.TypeVolume, Source: IndexVolumeName, Target: "/mnt/livegrep-index"},
			{Type: mount.TypeBind, Source: f.Name(), Target: "/mnt/manifest.json", ReadOnly: true},
		},
	}
	resp, err := d.docker.ContainerCreate(ctx, containerConfig, hostConfig, nil, "")
	if err != nil {
		log.Printf("failed to create a container for indexing %s: %v", manifest.Name, err)
	}

	err = d.docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		log.Printf("failed to start a container for indexing %s: %v", manifest.Name, err)
		return err
	}
	_, err = d.docker.ContainerWait(ctx, resp.ID)
	if err != nil {
		log.Printf("failed to wait a container for indexing %s: %v", manifest.Name, err)
		return err
	}
	log.Printf("Indexing %s completed", manifest.Name)
	return nil
}

func (d *runnerImpl) RerunIndexDB(ctx context.Context, project string) error {
	log.Printf("restarting index server for %s", project)
	name := IndexContainerName(project)
	err := d.docker.ContainerKill(ctx, name, "SIGTERM")
	if IsErrNoSuchContainer(err) {
		err := d.createIndexDB(ctx, project)
		if err != nil {
			log.Printf("failed to create indexing container for %s: %v", project, err)
			return err
		}
	} else if err != nil {
		log.Printf("Failed to kill index server for %s: %v", project, err)
		return err
	} else {
		log.Printf("Killed index server for %s", project)
	}

	err = d.docker.ContainerStart(ctx, name, types.ContainerStartOptions{})
	if err != nil {
		log.Printf("failed to start index server for %s: %v", project, err)
		return err
	}
	log.Printf("started index server for %s", project)
	return nil
}

// IndexPath returns an index path of the project
func IndexPath(project string) string {
	return "/livegrep-index/" + project + ".idx"
}

func (d *runnerImpl) createIndexDB(ctx context.Context, project string) error {
	name := IndexContainerName(project)
	cmd := []string{
		"/livegrep/bin/codesearch", "-load_index", "/mnt" + IndexPath(project), "-grpc", "0.0.0.0:9999",
	}
	containerConfig := &container.Config{
		Image:    Image,
		Cmd:      strslice.StrSlice(cmd),
		Hostname: IndexHostName(project),
	}
	hostConfig := &container.HostConfig{
		ReadonlyRootfs: true,
		Mounts: []mount.Mount{
			{Type: mount.TypeVolume, Source: IndexVolumeName, Target: "/mnt/livegrep-index", ReadOnly: true},
		},
	}
	_, err := d.docker.ContainerCreate(ctx, containerConfig, hostConfig, nil, name)
	return err
}
