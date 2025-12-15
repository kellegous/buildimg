package builder

type BuilderOption func(*Builder)

func WithCommander(commander Commander) BuilderOption {
	return func(b *Builder) {
		b.commander = commander
	}
}

func WithName(name string) BuilderOption {
	return func(b *Builder) {
		b.name = name
	}
}
