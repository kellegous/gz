package cmd

import (
	"github.com/kellegous/poop"
	"github.com/spf13/cobra"

	"github.com/kellegous/gz/internal/client"
)

type rebaseFlags struct {
	*rootFlags
	Root client.RootUpdate
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
	ctx := cmd.Context()

	c, err := client.Open(ctx, flags.root)
	if err != nil {
		return poop.Chain(err)
	}
	defer c.Close()

	return poop.Chain(c.Rebase(ctx, &client.RebaseOptions{
		Root: flags.Root,
	}))
}
