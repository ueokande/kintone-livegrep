package lgrunner

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func (d *runnerImpl) RunWeb(ctx context.Context, config WebConfig) error {
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

	args := []string{"run", "--read-only", "-d", "--rm",
		"--name=" + WebContainerName, "--hostname=" + WebHostName, "-p", "8910:8910",
		"-v", f.Name() + ":/etc/livegrep/livegrep.json:ro",
	}
	for _, b := range config.Backends {
		h := strings.Split(b.Address, ":")[0]
		args = append(args, "--link", h)
	}
	args = append(args, Image, "/livegrep/bin/livegrep", "-docroot", "/livegrep/web/", "/etc/livegrep/livegrep.json")

	stderr := new(bytes.Buffer)
	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Stderr = stderr
	err = cmd.Run()
	if err != nil {
		log.Printf("failed to run index server: %v: args=%v stderr=%s", err, args, stderr.String())
		return err
	}
	log.Printf("started web server")
	return nil
}

func (d *runnerImpl) StopWeb(ctx context.Context) error {
	stderr := new(bytes.Buffer)
	cmd := exec.CommandContext(ctx, "docker", "stop", WebContainerName)
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		log.Printf("failed to web server %v: stderr=%s", err, stderr.String())
		return err
	}
	log.Printf("stopped web server")
	return nil
}
