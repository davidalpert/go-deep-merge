package cmd

import (
	"github.com/davidalpert/go-printers/v1"
	"github.com/spf13/cobra"
)

func NewCmdAws(ioStreams printers.IOStreams) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "aws",
		Aliases: []string{"a"},
		Short:   "aws parameter store subcommands",
		Args:    cobra.NoArgs,
	}

	cmd.AddCommand(NewCmdAwsGet(ioStreams))

	return cmd
}
