package cmd

import (
	"github.com/spf13/cobra"
)

func storeCmd(rf *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "store",
		Short: "store commands are for debug/admin use only",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(storeEditCmd(rf))

	return cmd
}
