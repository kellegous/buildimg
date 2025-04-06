package internal

import "strings"

type StringsFlag struct {
	kind string
	Vals []string
}

func NewStringsFlag(kind string) StringsFlag {
	return StringsFlag{
		kind: kind,
	}
}

func (f *StringsFlag) Set(v string) error {
	f.Vals = append(f.Vals, v)
	return nil
}

func (f *StringsFlag) String() string {
	return strings.Join(f.Vals, ", ")
}

func (f *StringsFlag) Type() string {
	return f.kind
}
