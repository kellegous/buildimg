package internal

type Image struct {
	Root      string
	Dockfile  string
	Name      string
	Targets   []*Target
	BuildArgs []string
	Labels    []string
}

func (i *Image) toBuildCmds() []*buildCmd {
	var toPush []string
	var cmds []*buildCmd
	for _, target := range i.Targets {
		if target.push() {
			toPush = append(toPush, target.Platform)
			continue
		}

		cmds = append(cmds, &buildCmd{
			Root:       i.Root,
			Dockerfile: i.Dockfile,
			Name:       i.Name,
			Platforms:  []string{target.Platform},
			Dest:       target.Output,
			BuildArgs:  i.BuildArgs,
			Labels:     i.Labels,
		})
	}

	if len(toPush) > 0 {
		cmds = append(cmds, &buildCmd{
			Root:       i.Root,
			Dockerfile: i.Dockfile,
			Name:       i.Name,
			Platforms:  toPush,
			BuildArgs:  i.BuildArgs,
			Labels:     i.Labels,
		})
	}

	return cmds
}

type buildCmd struct {
	Root       string
	Dockerfile string
	Name       string
	Platforms  []string
	Dest       string
	BuildArgs  []string
	Labels     []string
}
