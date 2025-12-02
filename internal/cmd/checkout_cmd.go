package cmd

import (
	"github.com/kellegous/poop"
	"github.com/spf13/cobra"
)

func checkoutCmd(rf *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "checkout",
		Short:   "checkout a branch",
		Aliases: []string{"co", "switch"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := runCheckout(cmd, args[0]); err != nil {
				poop.HitFan(err)
			}
		},
	}
	return cmd
}

func runCheckout(cmd *cobra.Command, name string) error {
	return nil
}
