package cmd

import (
	"context"
	"os"
	"path/filepath"

	git "github.com/go-git/go-git/v6"
	"github.com/kellegous/gz/internal/store"
	"github.com/spf13/cobra"
)

type rootFlags struct {
	root string
}

func (r *rootFlags) repo() (*git.Repository, error) {
	return git.PlainOpen(r.root)
}

func (r *rootFlags) store(ctx context.Context) (*store.Store, error) {
	return store.Open(ctx, filepath.Join(r.root, ".git/gz.db"))
}

func rootCmd() *cobra.Command {
	var flags rootFlags

	cmd := &cobra.Command{
		Use:   "gz",
		Short: "gz is a tool for single commit, chained git branches",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.PersistentFlags().StringVarP(
		&flags.root,
		"root",
		"r",
		".",
		"the root directory of the project",
	)

	cmd.AddCommand(createCmd(&flags))

	return cmd
}

func Execute() {
	if err := rootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
