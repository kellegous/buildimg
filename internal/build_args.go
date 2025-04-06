package internal

import "strings"

type BuildArgs []string

func (a *BuildArgs) Set(v string) error {
	*a = append(*a, v)
	return nil
}

func (a *BuildArgs) String() string {
	return strings.Join(*a, ", ")
}

func (a *BuildArgs) Type() string {
	return "build args"
}
