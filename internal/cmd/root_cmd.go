package cmd

import (
	"context"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/kellegous/gz/internal/git"
	"github.com/kellegous/gz/internal/store"
)

type rootFlags struct {
	root string
}

func (r *rootFlags) store(ctx context.Context) (*store.Store, error) {
	return store.Open(ctx, filepath.Join(r.root, ".git/gz.db"))
}

func (r *rootFlags) workDir() (*git.WorkDir, error) {
	return git.WorkDirAt(r.root)
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

	cmd.AddCommand(checkoutCmd(&flags))
	cmd.AddCommand(createCmd(&flags))
	cmd.AddCommand(commitCmd(&flags))
	cmd.AddCommand(rebaseCmd(&flags))

	return cmd
}

func Execute() {
	if err := rootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
