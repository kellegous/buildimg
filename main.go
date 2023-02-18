package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/kellegous/buildimg/pkg"
)

func createBuilder(root string) (func() error, error) {
	cmd := exec.Command("docker", "buildx", "create", "--use")
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return func() error {
		cmd := exec.Command("docker", "buildx", "rm")
		cmd.Dir = root
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}, nil
}

func Build(
	root string,
	name string,
	dockerFile string,
	platforms []string,
	push bool,
) error {
	done, err := createBuilder(root)
	if err != nil {
		return err
	}
	defer done()

	cmd := []string{
		"buildx",
		"build",
		fmt.Sprintf("--platform=%s", strings.Join(platforms, ",")),
		fmt.Sprintf("--file=%s", dockerFile),
	}

	if push {
		cmd = append(cmd, "--push")
	} else {
		cmd = append(cmd, "--load")
	}

	cmd = append(cmd, "-t", name, ".")

	cb := exec.Command("docker", cmd...)
	cb.Dir = root
	cb.Stdout = os.Stdout
	cb.Stderr = os.Stderr

	return cb.Run()
}

func main() {
	var flags pkg.Flags
	name := flags.Parse()

	fullName := fmt.Sprintf("%s:%s", name, flags.Version)
	if err := Build(
		flags.Root,
		fullName,
		flags.Dockerfile,
		flags.Platforms,
		flags.Push,
	); err != nil {
		fmt.Fprintf(os.Stderr,
			"Build failed %s", fullName)
		os.Exit(2)
	}

	fmt.Printf("Published %s\n", fullName)
}
