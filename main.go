package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/kellegous/buildimg/pkg"
)

func Build(
	root string,
	name string,
	dockerFile string,
	platforms []string,
	push bool,
) error {
	ca := exec.Command("docker", "buildx", "create", "--use")
	ca.Dir = root
	ca.Stdout = os.Stdout
	ca.Stderr = os.Stderr
	if err := ca.Run(); err != nil {
		return err
	}

	cmd := []string{
		"buildx",
		"build",
		fmt.Sprintf("--platform=%s", strings.Join(platforms, ",")),
		fmt.Sprintf("--file=%s", dockerFile),
	}

	if push {
		cmd = append(cmd, "--push")
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
		log.Panic(err)
	}

	fmt.Printf("Published %s\n", fullName)
}
