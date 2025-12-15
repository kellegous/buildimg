package builder

import (
	"context"
	"encoding/json"
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

func assertSameCommands(t *testing.T,
	expected [][]string,
	actual [][]string,
) {
	if reflect.DeepEqual(expected, actual) {
		return
	}

	if len(expected) != len(actual) {
		t.Fatalf(
			"expected %d commands, got %d",
			len(expected),
			len(actual),
		)
	}

	for i, cmd := range expected {
		if !reflect.DeepEqual(cmd, actual[i]) {
			t.Fatalf(
				"command at %d differs\nexpected:%s\ngot:%s",
				i,
				describeAsJSON(t, cmd),
				describeAsJSON(t, actual[i]),
			)
		}
	}
}

func describeAsJSON(t *testing.T, v any) string {
	json, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	return string(json)
}

func TestBuilder(t *testing.T) {
	tests := []struct {
		name             string
		run              func(t *testing.T, c Commander)
		expectedCommands [][]string
	}{
		{
			name: "NamedBuilder",
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
			name: "UnnamedBuilder",
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

		{
			name: "BuilderImage",
			run: func(t *testing.T, c Commander) {
				ctx := t.Context()
				b, err := Start(ctx, WithCommander(c), WithName("test"))
				if err != nil {
					t.Fatal(err)
				}
				defer b.Shutdown(ctx)

				if err := b.Build(ctx, &Image{
					Path:       "test",
					Dockerfile: "Dockerfile",
					Name:       "foo/test:latest",
					Targets: []*Target{
						{Platform: "linux/amd64", Output: "amd64.tar"},
						{Platform: "linux/arm64", Output: "arm64.tar"},
						{Platform: "linux/amd64"},
						{Platform: "linux/arm64"},
					},
					BuildArgs: []string{"VERSION=one"},
					Labels:    []string{"ALABEL=foo"},
					Secrets:   []string{"id=secret,src=secret.txt"},
				}); err != nil {
					t.Fatal(err)
				}
			},
			expectedCommands: [][]string{
				{"docker", "buildx", "create", "--name", "test"},
				{
					"docker", "buildx", "build",
					"--platform=linux/amd64",
					"--file=Dockerfile",
					"--builder=test",
					"--build-arg", "VERSION=one",
					"--label", "ALABEL=foo",
					"--secret", "id=secret,src=secret.txt",
					"-o", "type=docker,dest=amd64.tar",
					"-t", "foo/test:latest",
					"test",
				},
				{
					"docker", "buildx", "build",
					"--platform=linux/arm64",
					"--file=Dockerfile",
					"--builder=test",
					"--build-arg", "VERSION=one",
					"--label", "ALABEL=foo",
					"--secret", "id=secret,src=secret.txt",
					"-o", "type=docker,dest=arm64.tar",
					"-t", "foo/test:latest",
					"test",
				},
				{
					"docker", "buildx", "build",
					"--platform=linux/amd64,linux/arm64",
					"--file=Dockerfile",
					"--builder=test",
					"--build-arg", "VERSION=one",
					"--label", "ALABEL=foo",
					"--secret", "id=secret,src=secret.txt",
					"--push",
					"-t",
					"foo/test:latest",
					"test",
				},
				{"docker", "buildx", "rm", "test"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := &testCommander{}
			test.run(t, c)
			assertSameCommands(t, test.expectedCommands, c.commands)
		})
	}
}
