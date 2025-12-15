package builder

import (
	"context"
	"reflect"
	"testing"
)

type testCommand func()

func (c testCommand) Run() error {
	c()
	return nil
}

type testCommander struct {
	commands [][]string
}

func (c *testCommander) Command(ctx context.Context, name string, args ...string) Command {
	return testCommand(func() {
		cmd := append([]string{name}, args...)
		c.commands = append(c.commands, cmd)
	})
}

func TestBuilder(t *testing.T) {
	tests := []struct {
		name             string
		run              func(t *testing.T, c Commander)
		expectedCommands [][]string
	}{
		{
			name: "named builder",
			run: func(t *testing.T, c Commander) {
				ctx := t.Context()
				b, err := Start(ctx, WithCommander(c), WithName("test"))
				if err != nil {
					t.Fatal(err)
				}
				defer b.Shutdown(ctx)
			},
			expectedCommands: [][]string{
				{"docker", "buildx", "create", "--name", "test"},
				{"docker", "buildx", "rm", "test"},
			},
		},

		{
			name: "unnamed builder",
			run: func(t *testing.T, c Commander) {
				ctx := t.Context()
				b, err := Start(ctx, WithCommander(c), WithNameGenerator(func() string {
					return "buildimg-12345678"
				}))
				if err != nil {
					t.Fatal(err)
				}
				defer b.Shutdown(ctx)
			},
			expectedCommands: [][]string{
				{"docker", "buildx", "create", "--name", "buildimg-12345678"},
				{"docker", "buildx", "rm", "buildimg-12345678"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := &testCommander{}
			test.run(t, c)
			if !reflect.DeepEqual(c.commands, test.expectedCommands) {
				t.Errorf("expected commands %v, got %v", test.expectedCommands, c.commands)
			}
		})
	}
}
