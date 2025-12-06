package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kellegous/gz/internal/client"
	"github.com/kellegous/gz/internal/store"
	"github.com/kellegous/poop"
	"github.com/spf13/cobra"
)

func storeGetCmd(rf *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get a branch from the database",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := runStoreGet(cmd, rf); err != nil {
				poop.HitFan(err)
			}
		},
	}
	return cmd
}

func runStoreGet(cmd *cobra.Command, flags *rootFlags) error {
	ctx := cmd.Context()

	c, err := client.Open(ctx, flags.root)
	if err != nil {
		return poop.Chain(err)
	}
	defer c.Close()

	branch, err := c.GetBranch(ctx)
	if errors.Is(err, store.ErrNotFound) {
		return nil
	} else if err != nil {
		return poop.Chain(err)
	}

	b, err := json.MarshalIndent(branch, "", "  ")
	if err != nil {
		return poop.Chain(err)
	}
	fmt.Println(string(b))

	return nil
}
