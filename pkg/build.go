package pkg

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type Image struct {
	Root     string
	Dockfile string
	Name     string
	Targets  []*Target
}

func (i *Image) toBuildCmds() []*buildCmd {
	var toPush []string
	var cmds []*buildCmd
	for _, target := range i.Targets {
		if target.push() {
			toPush = append(toPush, target.Platform)
			continue
		}

		cmds = append(cmds, &buildCmd{
			Root:       i.Root,
			Dockerfile: i.Dockfile,
			Name:       i.Name,
			Platforms:  []string{target.Platform},
			Dest:       target.Output,
		})
	}

	if len(toPush) > 0 {
		cmds = append(cmds, &buildCmd{
			Root:       i.Root,
			Dockerfile: i.Dockfile,
			Name:       i.Name,
			Platforms:  toPush,
		})
	}

	return cmds
}

type buildCmd struct {
	Root       string
	Dockerfile string
	Name       string
	Platforms  []string
	Dest       string
}

func (c *buildCmd) Build(ctx context.Context) error {
	args := []string{
		"buildx",
		"build",
		fmt.Sprintf("--platform=%s", strings.Join(c.Platforms, ",")),
		fmt.Sprintf("--file=%s", c.Dockerfile),
	}

	if c.Dest == "" {
		args = append(args, "--push")
	} else {
		args = append(args,
			"-o",
			fmt.Sprintf("type=docker,dest=%s", c.Dest))
	}

	args = append(args, "-t", c.Name, c.Root)

	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (i *Image) Build(ctx context.Context) error {
	done, err := createBuilder(ctx, i.Root)
	if err != nil {
		return err
	}
	defer done()

	for _, cmd := range i.toBuildCmds() {
		if err := cmd.Build(ctx); err != nil {
			return err
		}
	}

	return done()
}

func createBuilder(
	ctx context.Context,
	root string,
) (func() error, error) {
	cmd := exec.CommandContext(
		ctx,
		"docker",
		"buildx",
		"create",
		"--use")
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	didShutdown := false
	var lock sync.Mutex
	shutdown := func() error {
		lock.Lock()
		defer lock.Unlock()
		if didShutdown {
			return nil
		}
		didShutdown = true

		cmd := exec.Command("docker", "buildx", "rm")
		cmd.Dir = root
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	return shutdown, cmd.Run()
}
