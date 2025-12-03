package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/kellegous/gz/internal/client"
	"github.com/kellegous/poop"
	"github.com/spf13/cobra"
)

type createFlags struct {
	*rootFlags
	from string
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

	return cmd
}

func runCreate(
	cmd *cobra.Command,
	flags *createFlags,
	name string,
) error {
	ctx := cmd.Context()

	client, err := client.Open(ctx, flags.root)
	if err != nil {
		return poop.Chain(err)
	}
	defer client.Close()

	branch, err := client.CreateBranch(ctx, name, flags.from)
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
