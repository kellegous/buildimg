package builder

import (
	"context"
	"os"
	"os/exec"
)

type defaultCommander struct{}

func (c *defaultCommander) Command(
	ctx context.Context,
	name string,
	args ...string,
) Command {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
