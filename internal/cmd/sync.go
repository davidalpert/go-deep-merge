package cmd

import (
	"github.com/davidalpert/go-printers/v1"
	"github.com/spf13/cobra"
)

func NewCmdSync(ioStreams printers.IOStreams) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "sync",
		Aliases: []string{"s"},
		Short:   "sync subcommands",
		Args:    cobra.NoArgs,
	}

	//cmd.AddCommand(NewCmdSyncAWSParameterStore(ioStreams))
	//cmd.AddCommand(NewCmdSyncConsul(ioStreams))
	cmd.AddCommand(NewCmdSyncFolder(ioStreams))

	return cmd
}
