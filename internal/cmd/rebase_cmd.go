package cmd

import (
	"github.com/kellegous/poop"
	"github.com/spf13/cobra"
)

func rebaseCmd(rf *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rebase",
		Short: "rebase the current branch",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := runRebase(cmd); err != nil {
				poop.HitFan(err)
			}
		},
	}
	return cmd
}

func runRebase(cmd *cobra.Command) error {
	return nil
}
