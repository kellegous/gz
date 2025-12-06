package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/kellegous/poop"
	"github.com/spf13/cobra"

	"github.com/kellegous/gz/internal"
	"github.com/kellegous/gz/internal/client"
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

func toCommitOptions(flags *commitFlags) *client.CommitOptions {
	var msg client.MessageOption
	if flags.message != "" {
		msg = client.WithMessage(flags.message)
	} else if !flags.append && !flags.edit {
		msg = client.KeepExistingMessage()
	}

	return &client.CommitOptions{
		All:     true,
		Append:  flags.append,
		Message: msg,
	}
}

func runCommit(cmd *cobra.Command, flags *commitFlags) error {
	ctx := cmd.Context()

	c, err := client.Open(ctx, flags.root)
	if err != nil {
		return poop.Chain(err)
	}
	defer c.Close()

	branch, err := c.Commit(ctx, toCommitOptions(flags))
	if err != nil {
		return poop.Chain(err)
	}

	b, err := json.MarshalIndent(internal.BranchFromProto(branch), "", "  ")
	if err != nil {
		return poop.Chain(err)
	}
	fmt.Println(string(b))

	return nil
}
