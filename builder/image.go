package builder

type Image struct {
	Path       string
	Dockerfile string
	Name       string
	Targets    []*Target
	BuildArgs  []string
	Labels     []string
	Secrets    []string
}

func (i *Image) toBuildCmds() []*buildCmd {
	var toPush []*Target
	var cmds []*buildCmd
	for _, target := range i.Targets {
		if target.push() {
			toPush = append(toPush, target)
			continue
		}

		cmds = append(cmds, &buildCmd{
			Image:     i,
			Platforms: []string{target.Platform},
			Dest:      target.Output,
		})
	}

	if len(toPush) > 0 {
		platforms := make([]string, 0, len(toPush))
		for _, target := range toPush {
			platforms = append(platforms, target.Platform)
		}
		cmds = append(cmds, &buildCmd{
			Image:     i,
			Platforms: platforms,
		})
	}

	return cmds
}

type buildCmd struct {
	*Image
	Platforms []string
	Dest      string
}
