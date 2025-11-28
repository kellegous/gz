package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/go-git/go-git/v6/plumbing"
	"github.com/kellegous/gz"
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

	wd, err := flags.workDir()
	if err != nil {
		return poop.Chain(err)
	}

	repo := wd.Repository()

	var ref *plumbing.Reference
	if flags.from != "" {
		ref, err = repo.Reference(plumbing.NewBranchReferenceName(flags.from), true)
		if err != nil {
			return poop.Chain(err)
		}
	} else {
		ref, err = repo.Head()
		if err != nil {
			return poop.Chain(err)
		}
	}

	if err := wd.CreateBranch(ctx, name, flags.from); err != nil {
		return poop.Chain(err)
	}

	s, err := flags.store(ctx)
	if err != nil {
		return poop.Chain(err)
	}
	defer s.Close()

	branch, err := s.UpsertBranch(
		ctx,
		&gz.Branch{
			Name:   name,
			Parent: ref.Name().Short(),
		},
	)
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
