package builder

import (
	"io"
	"os"
)

type BuilderOptions struct {
	commander     Commander
	nameGenerator func() string
}

func (o *BuilderOptions) getCommander() Commander {
	if o.commander == nil {
		return &defaultCommander{
			stdout: os.Stdout,
			stderr: os.Stderr,
		}
	}
	return o.commander
}

func (o *BuilderOptions) getNameGenerator() func() string {
	if o.nameGenerator == nil {
		return defaultIdGenerator
	}
	return o.nameGenerator
}

type BuilderOption func(*BuilderOptions)

func WithCommander(commander Commander) BuilderOption {
	return func(o *BuilderOptions) {
		o.commander = commander
	}
}

func WithNameGenerator(generator func() string) BuilderOption {
	return func(o *BuilderOptions) {
		o.nameGenerator = generator
	}
}

func WithName(name string) BuilderOption {
	return func(o *BuilderOptions) {
		if name != "" {
			o.nameGenerator = func() string {
				return name
			}
		}
	}
}

func WithStdIO(stdout, stderr io.Writer) BuilderOption {
	return func(o *BuilderOptions) {
		o.commander = &defaultCommander{
			stdout: stdout,
			stderr: stderr,
		}
	}
}
