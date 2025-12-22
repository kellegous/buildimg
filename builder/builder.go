package builder

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
)

type Builder struct {
	lck          sync.Mutex
	didShutdown  bool
	name         string
	commander    Commander
	outputFormat OutputFormat
}

func (b *Builder) Shutdown(ctx context.Context) error {
	b.lck.Lock()
	defer b.lck.Unlock()

	if b.didShutdown {
		return nil
	}

	if err := b.commander.Command(
		ctx,
		"docker",
		"buildx",
		"rm",
		b.name,
	).Run(); err != nil {
		return err
	}

	b.didShutdown = true
	return nil
}

func (b *Builder) Build(
	ctx context.Context,
	img *Image,
) error {
	for cmd := range img.toBuildCmds() {
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

	if b.outputFormat != "" {
		args = append(args, fmt.Sprintf("--progress=%s", b.outputFormat))
	}

	for _, arg := range c.BuildArgs {
		args = append(args, "--build-arg", arg)
	}

	for _, label := range c.Labels {
		args = append(args, "--label", label)
	}

	for _, secret := range c.Secrets {
		args = append(args, "--secret", secret)
	}

	if c.Dest == "" {
		args = append(args, "--push")
	} else {
		args = append(args,
			"-o",
			fmt.Sprintf("type=docker,dest=%s", c.Dest))
	}

	args = append(args, "-t", c.Name, c.Path)

	return b.commander.Command(ctx, "docker", args...).Run()
}

func Start(
	ctx context.Context,
	opts ...BuilderOption,
) (*Builder, error) {
	var o BuilderOptions
	for _, opt := range opts {
		opt(&o)
	}

	b := Builder{
		commander: o.getCommander(),
		name:      o.getNameGenerator()(),
	}

	if err := b.start(ctx); err != nil {
		return nil, err
	}

	return &b, nil
}

func (b *Builder) start(ctx context.Context) error {
	return b.commander.Command(
		ctx,
		"docker",
		"buildx",
		"create",
		"--name",
		b.name,
	).Run()
}

func defaultIdGenerator() string {
	var key [8]byte
	rand.Read(key[:])
	return fmt.Sprintf("buildimg-%s", hex.EncodeToString(key[:]))
}
