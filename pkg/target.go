package pkg

import (
	"fmt"
	"strings"
)

type Target struct {
	Platform string
	Output   string
}

func (t *Target) Set(v string) error {
	p, o, _ := strings.Cut(v, ":")
	t.Platform = strings.TrimSpace(p)
	t.Output = strings.TrimSpace(o)
	return nil
}

func (t *Target) String() string {
	if t.Output == "" {
		return t.Platform
	}
	return fmt.Sprintf("%s:%s", t.Platform, t.Output)
}

func (t *Target) Type() string {
	return "target"
}

func (t *Target) push() bool {
	return t.Output == ""
}
