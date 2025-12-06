package cmd

import (
	"os"
	"slices"

	"github.com/kellegous/poop"
	"github.com/spf13/cobra"

	"github.com/kellegous/gz/internal/client"
)

type createFlags struct {
	*rootFlags
	from  string
	alias StringSet
}

func createCmd(rf *rootFlags) *cobra.Command {
	flags := createFlags{
		rootFlags: rf,
	}

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "create a new stacked feature branch",
		Args:    cobra.RangeArgs(1, 2),
		Aliases: []string{"+", "push"},
		Run: func(cmd *cobra.Command, args []string) {
			if err := runCreate(cmd, &flags, args[0]); err != nil {
				poop.HitFan(err)
			}
		},
	}

	cmd.Flags().StringVarP(
		&flags.from,
		"from",
		"f",
		"",
		"the branch to create from",
	)

	cmd.Flags().VarP(
		&flags.alias,
		"alias",
		"a",
		"the alias for the branch",
	)

	return cmd
}

func runCreate(
	cmd *cobra.Command,
	flags *createFlags,
	name string,
) error {
	ctx := cmd.Context()

	c, err := client.Open(ctx, flags.root)
	if err != nil {
		return poop.Chain(err)
	}
	defer c.Close()

	branch, err := c.CreateBranch(
		ctx,
		name,
		flags.from,
		slices.Collect(flags.alias.Values()),
	)
	if err != nil {
		return poop.Chain(err)
	}

	if err := branch.WriteJSONTo(os.Stdout); err != nil {
		return poop.Chain(err)
	}

	return nil
}
