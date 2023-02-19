package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"github.com/kellegous/buildimg/pkg"
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
	var targets pkg.Targets
	var root, dockerfile, tag string
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

			image := pkg.Image{
				Root:     root,
				Dockfile: dockerfile,
				Name:     fmt.Sprintf("%s:%s", args[0], tag),
				Targets:  targets,
			}

			if err := image.Build(ctx); err != nil {
				cmd.PrintErrf("build failed: %s", err)
				os.Exit(1)
			}

			fmt.Printf("%s\n", image.Name)
		},
	}

	cmd.Flags().Var(
		&targets,
		"target",
		"a build target (i.e. linux/amd64)")

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

	return cmd
}

func main() {
	if err := Command().Execute(); err != nil {
		os.Exit(1)
	}
}
