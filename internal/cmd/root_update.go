package cmd

import "fmt"

type RootUpdate int

const (
	RootUpdateNothing RootUpdate = iota
	RootUpdateFetchAndRebase
	RootUpdateRebase
)

func (r *RootUpdate) Set(v string) error {
	switch v {
	case "fetch-and-rebase", "f":
		*r = RootUpdateFetchAndRebase
		return nil
	case "rebase", "r":
		*r = RootUpdateRebase
		return nil
	case "nothing", "n":
		*r = RootUpdateNothing
		return nil
	}
	return fmt.Errorf("invalid root update: %s", v)
}

func (r *RootUpdate) String() string {
	switch *r {
	case RootUpdateFetchAndRebase:
		return "fetch-and-rebase"
	case RootUpdateRebase:
		return "rebase"
	case RootUpdateNothing:
		return "nothing"
	}
	return ""
}

func (r *RootUpdate) Type() string {
	return "string"
}
