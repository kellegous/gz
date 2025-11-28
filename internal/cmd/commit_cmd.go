package cmd

import (
	"github.com/kellegous/gz"
	"github.com/kellegous/gz/internal/git"
	"github.com/kellegous/poop"
	"github.com/spf13/cobra"
)

type commitFlags struct {
	*rootFlags
	message string
	edit    bool
	append  bool
}

func commitCmd(rf *rootFlags) *cobra.Command {
	flags := commitFlags{
		rootFlags: rf,
	}

	cmd := &cobra.Command{
		Use:     "commit",
		Short:   "commit current changes into branch",
		Args:    cobra.NoArgs,
		Aliases: []string{"save"},
		Run: func(cmd *cobra.Command, args []string) {
			if err := runCommit(cmd, &flags); err != nil {
				poop.HitFan(err)
			}
		},
	}

	// -a is not a good shortcut
	cmd.Flags().BoolVarP(
		&flags.append,
		"append",
		"a",
		false,
		"append the commit to the branch",
	)

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

	wd, err := flags.workDir()
	if err != nil {
		return poop.Chain(err)
	}

	s, err := flags.store(ctx)
	if err != nil {
		return poop.Chain(err)
	}
	defer s.Close()

	repo := wd.Repository()

	head, err := repo.Head()
	if err != nil {
		return poop.Chain(err)
	}

	branch, err := s.GetBranch(ctx, head.Name().Short())
	if err != nil {
		return poop.Chain(err)
	}

	// TODO(kellegous): validate flags because some flags are mutually exclusive

	amend := len(branch.Commits) > 0 && !flags.append

	var msg *git.Msg
	if flags.message != "" {
		msg = git.Message(flags.message)
	} else if amend && !flags.edit {
		msg = git.NoEdit()
	}

	// we are in append mode
	ref, err := wd.Commit(ctx, git.CommitOptions{
		All:     true,
		Message: msg,
		Amend:   amend,
	})
	if err != nil {
		return poop.Chain(err)
	}

	branch, err = s.UpdateBranch(ctx, &gz.Branch{
		Name:        branch.Name,
		Commits:     append(branch.Commits, ref.Hash().Bytes()),
		Parent:      branch.Parent,
		Description: branch.Description,
	})

	return nil
}
