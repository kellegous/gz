package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/kellegous/gz/internal/client"
	"github.com/kellegous/poop"
	"github.com/spf13/cobra"
)

func resetCmd(rf *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset",
		Short: "reset the current branch to point to the parent's HEAD",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := runReset(cmd, rf); err != nil {
				poop.HitFan(err)
			}
		},
	}
	return cmd
}

func runReset(cmd *cobra.Command, flags *rootFlags) error {
	ctx := cmd.Context()

	c, err := client.Open(ctx, flags.root)
	if err != nil {
		return poop.Chain(err)
	}
	defer c.Close()

	branch, err := c.Reset(ctx)
	if err != nil {
		return poop.Chain(err)
	}

	b, err := json.MarshalIndent(branch, "", "  ")
	if err != nil {
		return poop.Chain(err)
	}
	fmt.Println(string(b))

	return nil
}
