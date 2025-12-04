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
		Aliases: []string{"co", "switch"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := runCheckout(cmd, rf, args[0]); err != nil {
				poop.HitFan(err)
			}
		},
	}
	return cmd
}

func runCheckout(cmd *cobra.Command, flags *rootFlags, name string) error {
	ctx := cmd.Context()

	client, err := client.Open(ctx, flags.root)
	if err != nil {
		return poop.Chain(err)
	}
	defer client.Close()

	return poop.Chain(client.Checkout(ctx, name))
}
