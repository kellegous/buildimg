package internal

import "strings"

type Targets []*Target

func (t *Targets) Set(v string) error {
	var target Target
	if err := target.Set(v); err != nil {
		return err
	}
	*t = append(*t, &target)
	return nil
}

func (t *Targets) String() string {
	vals := make([]string, 0, len(*t))
	for _, target := range *t {
		vals = append(vals, target.String())
	}
	return strings.Join(vals, ", ")
}

func (t *Targets) Type() string {
	return "targets"
}
