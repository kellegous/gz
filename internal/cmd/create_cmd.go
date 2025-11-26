package cmd

import (
	"fmt"

	git "github.com/go-git/go-git/v6"
	"github.com/kellegous/poop"
	"github.com/spf13/cobra"
)

func createCmd(rf *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a new stacked feature branch",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			if err := runCreate(cmd, rf, args); err != nil {
				poop.HitFan(err)
			}
		},
	}

	return cmd
}

func runCreate(cmd *cobra.Command, rf *rootFlags, args []string) error {
	r, err := git.PlainOpen(rf.root)
	if err != nil {
		return poop.Chain(err)
	}

	fmt.Println(r)

	return nil
}
