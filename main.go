package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/kellegous/buildimg/builder"
	"github.com/kellegous/buildimg/internal"
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
	var targets builder.Targets
	var path, dockerfile, tag, builderName string
	buildArgs := internal.NewStringsFlag("build args")
	labels := internal.NewStringsFlag("labels")
	secrets := internal.NewStringsFlag("secrets")

	cmd := &cobra.Command{
		Use:   "buildimg [flags] name",
		Short: "automation for building and pushing images",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx, done := signal.NotifyContext(
				cmd.Context(),
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

			if path == "" {
				path = filepath.Dir(dockerfile)
			}

			b, err := builder.Start(ctx, builder.WithName(builderName))
			if err != nil {
				cmd.PrintErrf("start builder: %s", err)
				os.Exit(1)
			}
			defer b.Shutdown(context.Background())

			image := builder.Image{
				Path:       path,
				Dockerfile: dockerfile,
				Name:       fmt.Sprintf("%s:%s", args[0], tag),
				Targets:    targets,
				BuildArgs:  buildArgs.Vals,
				Labels:     labels.Vals,
				Secrets:    secrets.Vals,
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

	cmd.Flags().Var(
		&labels,
		"label",
		"a label (i.e. FOO=bar)")

	cmd.Flags().Var(
		&secrets,
		"secret",
		"a secret (i.e. id=github-token,src=github-token.txt)",
	)

	cmd.Flags().StringVar(
		&path,
		"path",
		"",
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
		&builderName,
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
