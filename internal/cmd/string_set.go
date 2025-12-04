package cmd

import (
	"iter"
	"slices"
	"strings"
)

type StringSet struct {
	values map[string]struct{}
}

func (s *StringSet) Set(v string) error {
	if s.values == nil {
		s.values = make(map[string]struct{})
	}
	s.values[v] = struct{}{}
	return nil
}

func (s *StringSet) Values() iter.Seq[string] {
	return func(yield func(string) bool) {
		if s.values == nil {
			return
		}

		for v := range s.values {
			if !yield(v) {
				return
			}
		}
	}
}

func (s *StringSet) String() string {
	return strings.Join(slices.Collect(s.Values()), ", ")
}

func (s *StringSet) Type() string {
	return "string"
}
