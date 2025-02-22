package pkg

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type Builder struct {
	lck         sync.Mutex
	name        string
	didShutdown bool
}

func (b *Builder) Shutdown(ctx context.Context) error {
	b.lck.Lock()
	defer b.lck.Unlock()

	if b.didShutdown {
		return nil
	}

	cmd := exec.CommandContext(
		ctx,
		"docker",
		"buildx",
		"rm",
		b.name)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}
	b.didShutdown = true
	return nil
}

func (b *Builder) Build(
	ctx context.Context,
	img *Image,
) error {
	for _, cmd := range img.toBuildCmds() {
		if err := b.build(ctx, cmd); err != nil {
			return err
		}
	}
	return nil
}

func (b *Builder) build(
	ctx context.Context,
	c *buildCmd,
) error {
	args := []string{
		"buildx",
		"build",
		fmt.Sprintf("--platform=%s", strings.Join(c.Platforms, ",")),
		fmt.Sprintf("--file=%s", c.Dockerfile),
		fmt.Sprintf("--builder=%s", b.name),
	}

	for _, arg := range c.BuildArgs {
		args = append(args, "--build-arg", arg)
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

func StartBuilder(
	ctx context.Context,
	root string,
	name string,
) (*Builder, error) {
	name, err := nameForBuilder(name)
	if err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(
		ctx,
		"docker",
		"buildx",
		"create",
		"--name",
		name)

	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return &Builder{name: name}, nil
}

func nameForBuilder(name string) (string, error) {
	if name != "" {
		return name, nil
	}

	var key [8]byte
	if _, err := rand.Read(key[:]); err != nil {
		return "", err
	}

	return fmt.Sprintf("buildimg-%s", hex.EncodeToString(key[:])), nil
}
