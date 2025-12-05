package cmd

import (
	"github.com/kellegous/poop"
	"github.com/spf13/cobra"

	"github.com/kellegous/gz/internal/client"
)

func checkoutCmd(rf *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "checkout",
		Short:   "checkout a branch",
		Aliases: []string{"co"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := runCheckout(cmd, rf, args[0]); err != nil {
				poop.HitFan(err)
			}
		},
	}
	return cmd
}

// TODO(kellegous): There should be fast ways do to the following:
// 1. checkout the root branch (co :r)
// 2. checkout the parent branch of the current branch (co :p)
// 3. checkout the previous branch you were on (co :-)
func runCheckout(cmd *cobra.Command, flags *rootFlags, name string) error {
	ctx := cmd.Context()

	client, err := client.Open(ctx, flags.root)
	if err != nil {
		return poop.Chain(err)
	}
	defer client.Close()

	return poop.Chain(client.Checkout(ctx, name))
}
