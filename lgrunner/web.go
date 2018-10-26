package lgrunner

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
)

func (d *runnerImpl) RerunWeb(ctx context.Context, config WebConfig) error {
	log.Printf("restarting web server")
	f, err := ioutil.TempFile("", "livegreptone-web")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())

	err = json.NewEncoder(f).Encode(config)
	if err != nil {
		log.Printf("failed to write web config: %v", err)
		return err
	}

	links := make([]string, len(config.Backends))
	for i, b := range config.Backends {
		h := strings.Split(b.Address, ":")[0]
		links[i] = h
	}

	err = d.docker.ContainerStop(ctx, WebContainerName, nil)
	if IsErrNoSuchContainer(err) {
		err := d.createWeb(ctx, f.Name(), links)
		if err != nil {
			log.Printf("failed to create web container: %v", err)
			return err
		}
	} else if err != nil {
		log.Printf("Failed to kill web server: %v", err)
		return err
	} else {
		log.Printf("Killed web server")
	}

	// NOTE: Wait to ensure container removed
	time.Sleep(3 * time.Second)

	err = d.docker.ContainerStart(ctx, WebContainerName, types.ContainerStartOptions{})
	if err != nil {
		log.Printf("failed to start web: %v", err)
		return err
	}
	log.Printf("started web server")
	return nil
}

func (d *runnerImpl) createWeb(ctx context.Context, confpath string, links []string) error {
	cmd := []string{
		"/livegrep/bin/livegrep", "-docroot", "/livegrep/web/", "/etc/livegrep/livegrep.json",
	}

	containerConfig := &container.Config{
		Image:    Image,
		Cmd:      strslice.StrSlice(cmd),
		Hostname: WebHostName,
		ExposedPorts: nat.PortSet{
			"8910/tcp": struct{}{},
		},
	}
	hostConfig := &container.HostConfig{
		AutoRemove:     true,
		ReadonlyRootfs: true,
		Mounts: []mount.Mount{
			{Type: mount.TypeBind, Source: confpath, Target: "/etc/livegrep/livegrep.json", ReadOnly: true},
		},
		PortBindings: map[nat.Port][]nat.PortBinding{
			"8910/tcp": []nat.PortBinding{
				{HostIP: "0.0.0.0", HostPort: "8910"},
			},
		},
		Links: links,
	}

	_, err := d.docker.ContainerCreate(ctx, containerConfig, hostConfig, nil, WebContainerName)
	return err
}
