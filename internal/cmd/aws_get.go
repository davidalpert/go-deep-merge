package cmd

import (
	"fmt"
	"github.com/davidalpert/go-deep-merge/internal/paramstore"
	"github.com/davidalpert/go-printers/v1"
	"github.com/spf13/cobra"
)

type AwsGetOptions struct {
	*printers.PrinterOptions
	ParameterName string
	Decrypt       bool
	Debug         bool
}

func NewAwsGetOptions(ioStreams printers.IOStreams) *AwsGetOptions {
	return &AwsGetOptions{
		PrinterOptions: printers.NewPrinterOptions().WithStreams(ioStreams).WithDefaultOutput("text"),
	}
}

func NewCmdAwsGet(ioStreams printers.IOStreams) *cobra.Command {
	o := NewAwsGetOptions(ioStreams)
	var cmd = &cobra.Command{
		Use:     "get <path>",
		Short:   "get a config value",
		Aliases: []string{"g"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(cmd, args); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			return o.Run()
		},
	}

	o.PrinterOptions.AddPrinterFlags(cmd.Flags())
	cmd.Flags().BoolVar(&o.Decrypt, "decrypt", false, "decrypt result?")
	cmd.Flags().BoolVarP(&o.Debug, "debug", "d", false, "debug")

	return cmd
}

// Complete the options
func (o *AwsGetOptions) Complete(cmd *cobra.Command, args []string) error {
	o.ParameterName = args[0]
	return nil
}

// Validate the options
func (o *AwsGetOptions) Validate() error {
	return o.PrinterOptions.Validate()
}

// Run the command
func (o *AwsGetOptions) Run() error {
	ssmsvc := paramstore.NewSSMClient(o.Debug)
	result, err := ssmsvc.GetValue(o.ParameterName, o.Decrypt)
	if err != nil {
		return fmt.Errorf("get param %#v: %#v", o.ParameterName, err)
	}

	return o.WriteOutput(result)
}
