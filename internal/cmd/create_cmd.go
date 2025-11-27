package cmd

import (
	"encoding/json"
	"fmt"

	git "github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/kellegous/gz"
	"github.com/kellegous/poop"
	"github.com/spf13/cobra"
)

func createCmd(rf *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "create a new stacked feature branch",
		Args:    cobra.RangeArgs(1, 2),
		Aliases: []string{"+", "push"},
		Run: func(cmd *cobra.Command, args []string) {
			if err := runCreate(cmd, rf, args); err != nil {
				poop.HitFan(err)
			}
		},
	}

	return cmd
}

func runCreate(cmd *cobra.Command, rf *rootFlags, args []string) error {
	ctx := cmd.Context()

	name := args[0]

	r, err := rf.repo()
	if err != nil {
		return poop.Chain(err)
	}

	wt, err := r.Worktree()
	if err != nil {
		return poop.Chain(err)
	}

	var ref *plumbing.Reference
	if len(args) == 1 {
		ref, err = r.Head()
		if err != nil {
			return poop.Chain(err)
		}
	} else {
		ref, err = r.Reference(plumbing.NewBranchReferenceName(args[1]), true)
		if err != nil {
			return poop.Chain(err)
		}
	}

	if err := wt.Checkout(&git.CheckoutOptions{
		Hash:   ref.Hash(),
		Create: true,
		Keep:   true,
		Branch: plumbing.NewBranchReferenceName(name),
	}); err != nil {
		return poop.Chain(err)
	}

	s, err := rf.store(cmd.Context())
	if err != nil {
		return poop.Chain(err)
	}
	defer s.Close()

	branch, err := s.UpsertBranch(
		ctx,
		&gz.Branch{
			Name:   name,
			Sha:    ref.Hash().Bytes(),
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
