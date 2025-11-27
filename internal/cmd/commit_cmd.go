package cmd

import (
	git "github.com/go-git/go-git/v6"
	"github.com/kellegous/gz"
	"github.com/kellegous/gz/internal/editor"
	"github.com/kellegous/poop"
	"github.com/spf13/cobra"
)

type commitFlags struct {
	*rootFlags
	message string
	edit    bool
}

func commitCmd(rf *rootFlags) *cobra.Command {
	flags := commitFlags{
		rootFlags: rf,
	}

	cmd := &cobra.Command{
		Use:   "commit",
		Short: "commit current changes into branch",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := runCommit(cmd, &flags); err != nil {
				poop.HitFan(err)
			}
		},
	}

	cmd.Flags().StringVarP(
		&flags.message,
		"message",
		"m",
		"",
		"the message for the commit",
	)
	cmd.Flags().BoolVarP(
		&flags.edit,
		"edit",
		"e",
		false,
		"edit the commit message",
	)

	return cmd
}

func runCommit(cmd *cobra.Command, flags *commitFlags) error {
	ctx := cmd.Context()

	r, err := flags.repo()
	if err != nil {
		return poop.Chain(err)
	}

	s, err := flags.store(ctx)
	if err != nil {
		return poop.Chain(err)
	}
	defer s.Close()

	head, err := r.Head()
	if err != nil {
		return poop.Chain(err)
	}

	branch, err := s.GetBranch(ctx, head.Name().Short())
	if err != nil {
		return poop.Chain(err)
	}

	wt, err := r.Worktree()
	if err != nil {
		return poop.Chain(err)
	}

	if branch.Sha == nil {
		// we need to make a new commit
		message, err := editor.EditFrom(ctx, r, "")
		if err != nil {
			return poop.Chain(err)
		}

		commit, err := wt.Commit(message, &git.CommitOptions{
			All: true,
		})
		if err != nil {
			return poop.Chain(err)
		}

		branch.Sha = commit.Bytes()
		branch, err = s.UpdateBranch(ctx, &gz.Branch{
			Name:   branch.Name,
			Sha:    commit.Bytes(),
			Parent: branch.Parent,
		})
	} else {
		// we are amending the existing commit
	}

	return nil
}
