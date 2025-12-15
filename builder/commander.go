package builder

import "context"

type Command interface {
	Run() error
}

type Commander interface {
	Command(
		ctx context.Context,
		name string,
		args ...string,
	) Command
}
