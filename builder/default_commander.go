package builder

import (
	"context"
	"io"
	"os/exec"
)

type defaultCommander struct {
	stdout io.Writer
	stderr io.Writer
}

func (c *defaultCommander) Command(
	ctx context.Context,
	name string,
	args ...string,
) Command {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = c.stdout
	cmd.Stderr = c.stderr
	return cmd
}
