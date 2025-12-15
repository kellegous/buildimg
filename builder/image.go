package builder

import "iter"

type Image struct {
	Path       string
	Dockerfile string
	Name       string
	Targets    []*Target
	BuildArgs  []string
	Labels     []string
	Secrets    []string
}

func (i *Image) toBuildCmds() iter.Seq[*buildCmd] {
	return func(yield func(*buildCmd) bool) {
		var platformsToPush []string

		for _, target := range i.Targets {
			if target.push() {
				platformsToPush = append(platformsToPush, target.Platform)
			} else {
				if !yield(&buildCmd{
					Image:     i,
					Platforms: []string{target.Platform},
					Dest:      target.Output,
				}) {
					return
				}
			}
		}

		if len(platformsToPush) > 0 {
			if !yield(&buildCmd{
				Image:     i,
				Platforms: platformsToPush,
			}) {
				return
			}
		}
	}
}

type buildCmd struct {
	*Image
	Platforms []string
	Dest      string
}
