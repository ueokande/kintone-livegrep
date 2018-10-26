package git

import (
	"bytes"
	"context"
	"errors"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	// ErrInvalidURL represents invalid URL error
	ErrInvalidURL = errors.New("invalid URL")
)

type vcs struct {
	rootDir string
}

// New creates new vcs.Interface implementation
func New(rootDir string) Interface {
	return &vcs{rootDir}
}

// Update is an implementation of the Interface
func (v *vcs) Update(ctx context.Context, url string, branch string) (string, error) {
	path, err := v.path(url, branch)
	if err != nil {
		return "", ErrInvalidURL
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := v.clone(ctx, url, branch, path)
		if err != nil {
			return "", err
		}
	} else {
		err := v.checkoutUpdated(ctx, path, branch)
		if err != nil {
			return "", err
		}
	}
	commit, err := v.getHEAD(ctx, path)
	if err != nil {
		return "", err
	}
	return commit, err
}

func (v *vcs) clone(ctx context.Context, url string, branch string, path string) error {
	log.Printf("[git] Cloning repository url=%s branch=%s", url, branch)
	stderr := new(bytes.Buffer)

	cmd := exec.CommandContext(ctx, "git", "clone", "-b", branch, url, path)
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		log.Printf("[git] failed to git clone error=%v, url=%s, branch=%s, path=%s, stderr=%s", err, url, branch, path, stderr.String())
		return err
	}
	log.Printf("[git] cloned repository url=%s branch=%s", url, branch)
	return nil
}

func (v *vcs) checkoutUpdated(ctx context.Context, path string, branch string) error {
	stderr := new(bytes.Buffer)
	cmd := exec.CommandContext(ctx, "git", "remote", "update")
	cmd.Stderr = stderr
	cmd.Dir = path
	err := cmd.Run()
	if err != nil {
		log.Printf("[git] failed to git remote update, path=%s, stderr=%s", path, stderr.String())
		return err
	}

	stderr = new(bytes.Buffer)
	cmd = exec.CommandContext(ctx, "git", "checkout", "-qf", "origin/"+branch)
	cmd.Stderr = stderr
	cmd.Dir = path
	err = cmd.Run()
	if err != nil {
		log.Printf("[git] failed to git checkout, path=%s, branch=%s, stderr=%s", path, branch, stderr.String())
		return err
	}
	log.Printf("[git] checked out path=%s branch=%s", path, branch)
	return nil
}

func (v *vcs) getHEAD(ctx context.Context, path string) (string, error) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	cmd := exec.CommandContext(ctx, "git", "rev-parse", "HEAD")
	cmd.Dir = path
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		log.Printf("[git] failed to git rev-parse HEAD, path=%s, stderr=%s", path, stderr.String())
		return "", err
	}
	return strings.TrimSpace(stdout.String()), nil
}

// Remove is an implementation of the Interface
func (v *vcs) Remove(ctx context.Context, url string, branch string) error {
	path, err := v.path(url, branch)
	if err != nil {
		return err
	}
	return os.RemoveAll(path)
}

func (v *vcs) path(repoURL string, branch string) (string, error) {
	url, err := url.Parse(repoURL)
	if err != nil {
		return "", err
	}
	return filepath.Join(v.rootDir, url.Host, url.Path, branch), nil
}
