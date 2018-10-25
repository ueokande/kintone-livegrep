package lgrunner

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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

	indexPath := IndexPath(manifest.Name)
	stderr := new(bytes.Buffer)
	cmd := exec.CommandContext(ctx, "docker", "run", "--rm", "--read-only",
		"-v", d.gitRootFS+":/mnt/livegrep-repos:ro",
		"-v", IndexVolumeName+":/mnt/livegrep-index:rw",
		"-v", f.Name()+":/mnt/manifest.json:ro",
		Image, "/livegrep/bin/codesearch",
		"-index_only", "-dump_index", "/mnt"+indexPath, "/mnt/manifest.json")
	cmd.Stdout = os.Stdout
	cmd.Stderr = stderr
	err = cmd.Run()
	if err != nil {
		log.Printf("failed to index %s: %v: args=%v, stderr=%s", manifest.Name, err, cmd.Args, stderr.String())
		return err
	}
	log.Printf("indexed %s", manifest.Name)
	return nil
}

func (d *runnerImpl) RunIndexDB(ctx context.Context, project string) error {
	hostname := IndexHostName(project)
	indexPath := IndexPath(project)
	stderr := new(bytes.Buffer)
	name := IndexContainerName(project)
	cmd := exec.CommandContext(ctx, "docker", "run", "--read-only", "-d", "--rm", "--name="+name, "--hostname="+hostname,
		"-v", IndexVolumeName+":/mnt/livegrep-index:ro",
		Image, "/livegrep/bin/codesearch", "-load_index", "/mnt"+indexPath, "-grpc", "0.0.0.0:9999")
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		log.Printf("failed to run index server: project=%s: %v: stderr=%s", project, err, stderr.String())
		return err
	}
	log.Printf("started index server: project=%s", project)
	return nil
}

func (d *runnerImpl) StopIndexDB(ctx context.Context, project string) error {
	name := IndexContainerName(project)
	stderr := new(bytes.Buffer)
	cmd := exec.CommandContext(ctx, "docker", "stop", name)
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		log.Printf("failed to stop index server %s: %v: stderr=%s", project, err, stderr.String())
		return err
	}
	log.Printf("stopped index server: project=%s", project)
	return nil
}

// IndexPath returns an index path of the project
func IndexPath(project string) string {
	return "/livegrep-index/" + project + ".idx"
}
