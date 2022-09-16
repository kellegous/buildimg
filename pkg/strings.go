package pkg

import "strings"

type Strings []string

func (s *Strings) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func (s *Strings) String() string {
	return strings.Join(*s, ", ")
}
