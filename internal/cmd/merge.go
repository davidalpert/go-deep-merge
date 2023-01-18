package cmd

import (
	"github.com/davidalpert/go-printers/v1"
	"github.com/spf13/cobra"
)

func NewCmdMerge(ioStreams printers.IOStreams) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "merge",
		Aliases: []string{"m"},
		Short:   "merge subcommands",
		Args:    cobra.NoArgs,
	}

	cmd.AddCommand(NewCmdMergeFiles(ioStreams))

	return cmd
}
