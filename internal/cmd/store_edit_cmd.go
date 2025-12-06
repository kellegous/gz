package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/kellegous/gz/internal/client"
	"github.com/kellegous/poop"
	"github.com/spf13/cobra"
)

func storeEditCmd(rf *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "edit the current branch in the database",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := runStoreEdit(cmd, rf); err != nil {
				poop.HitFan(err)
			}
		},
	}
	return cmd
}

func runStoreEdit(cmd *cobra.Command, flags *rootFlags) error {
	ctx := cmd.Context()

	c, err := client.Open(ctx, flags.root)
	if err != nil {
		return poop.Chain(err)
	}
	defer c.Close()

	branch, err := c.EditBranch(ctx)
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
