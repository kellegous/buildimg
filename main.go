package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"github.com/kellegous/buildimg/internal"
	"github.com/spf13/cobra"
)

func getTagFromGit(ctx context.Context) (string, error) {
	var buf bytes.Buffer

	cmd := exec.CommandContext(ctx, "git", "rev-parse", "HEAD")
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		return "", err
	}

	sha := strings.TrimSpace(buf.String())
	if len(sha) > 8 {
		sha = sha[:8]
	}

	return sha, nil
}

func Command() *cobra.Command {
	var targets internal.Targets
	var root, dockerfile, tag, builder string
	var buildArgs internal.BuildArgs
	cmd := &cobra.Command{
		Use:   "buildimg [flags] name",
		Short: "automation for building and pushing images",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx, done := signal.NotifyContext(
				context.Background(),
				os.Interrupt)
			defer done()

			var err error

			if tag == "" {
				tag, err = getTagFromGit(ctx)
				if err != nil {
					cmd.PrintErrf("git rev-parse: %s", err)
					os.Exit(1)
				}
			}

			b, err := internal.StartBuilder(ctx, root, builder)
			if err != nil {
				cmd.PrintErrf("start builder: %s", err)
				os.Exit(1)
			}
			defer b.Shutdown(context.Background())

			image := internal.Image{
				Root:      root,
				Dockfile:  dockerfile,
				Name:      fmt.Sprintf("%s:%s", args[0], tag),
				Targets:   targets,
				BuildArgs: buildArgs,
			}

			if err := b.Build(ctx, &image); err != nil {
				cmd.PrintErrf("build: %s", err)
				os.Exit(1)
			}

			fmt.Printf("%s\n", image.Name)
		},
	}

	cmd.Flags().Var(
		&targets,
		"target",
		"a build target (i.e. linux/amd64)")

	cmd.Flags().Var(
		&buildArgs,
		"build-arg",
		"a build arg (i.e. FOO=bar)")

	cmd.Flags().StringVar(
		&root,
		"root",
		".",
		"the build context directory")

	cmd.Flags().StringVar(
		&dockerfile,
		"dockerfile",
		"./Dockerfile",
		"the Dockerfile to use for the build")

	cmd.Flags().StringVar(
		&tag,
		"tag",
		"",
		"the image tag for the build (default is based on git sha)")

	cmd.Flags().StringVar(
		&builder,
		"builder",
		"",
		"the name of the builder to use")

	return cmd
}

func main() {
	if err := Command().Execute(); err != nil {
		os.Exit(1)
	}
}
