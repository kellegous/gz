package cmd

import (
	"github.com/kellegous/poop"
	"github.com/spf13/cobra"

	"github.com/kellegous/gz/internal/client"
)

func aliasCmd(rf *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alias",
		Short: "alias a branch",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if err := runAlias(cmd, rf, args[0], args[1:]); err != nil {
				poop.HitFan(err)
			}
		},
	}
	return cmd
}

func runAlias(
	cmd *cobra.Command,
	flags *rootFlags,
	name string,
	aliases []string,
) error {
	ctx := cmd.Context()

	c, err := client.Open(ctx, flags.root)
	if err != nil {
		return poop.Chain(err)
	}
	defer c.Close()

	return poop.Chain(c.Alias(ctx, name, aliases))
}
