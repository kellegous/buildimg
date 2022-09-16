package pkg

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultRoot      = "."
	defaultDockefile = "./Dockerfile"
)

var defaultPlatforms = []string{"linux/amd64"}

type Flags struct {
	Version    string
	Dockerfile string
	Root       string
	Push       bool
	Platforms  Strings
}

func (f *Flags) register(fs *flag.FlagSet) {
	fs.StringVar(
		&f.Version,
		"version",
		"latest",
		"the version tag for the image")

	fs.StringVar(
		&f.Root,
		"root",
		defaultRoot,
		"the root directory from which to build")

	fs.StringVar(
		&f.Dockerfile,
		"docker-file",
		defaultDockefile,
		"the docker file to use for building")

	fs.BoolVar(
		&f.Push,
		"push",
		false,
		"whether the image should be pushed")

	fs.Var(
		&f.Platforms,
		"platform",
		"the platforms to build built")
}

func (f *Flags) Parse() string {
	fs := flag.NewFlagSet(
		filepath.Base(os.Args[0]),
		flag.ExitOnError)

	usage := func() {
		fmt.Fprintf(fs.Output(), "usage: %s [options] image-name\n", fs.Name())
		fs.PrintDefaults()
		os.Exit(2)
	}

	fs.Usage = usage
	f.register(fs)
	fs.Parse(os.Args[1:])

	if fs.NArg() != 1 {
		fs.Usage()
	}

	if len(f.Platforms) == 0 {
		f.Platforms = defaultPlatforms
	}

	return fs.Arg(0)
}
