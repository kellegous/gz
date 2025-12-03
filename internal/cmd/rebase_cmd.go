package cmd

import (
	"fmt"

	"github.com/kellegous/poop"
	"github.com/spf13/cobra"
)

type rebaseFlags struct {
	*rootFlags
	Root RootUpdate
}

func rebaseCmd(rf *rootFlags) *cobra.Command {
	flags := rebaseFlags{
		rootFlags: rf,
	}

	cmd := &cobra.Command{
		Use:   "rebase",
		Short: "rebase the current branch",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := runRebase(cmd, &flags); err != nil {
				poop.HitFan(err)
			}
		},
	}

	cmd.Flags().VarP(
		&flags.Root,
		"root",
		"r",
		"the root update strategy",
	)
	return cmd
}

func runRebase(cmd *cobra.Command, flags *rebaseFlags) error {
	fmt.Println(flags)
	return nil
}
