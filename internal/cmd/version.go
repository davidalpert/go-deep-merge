package cmd

import (
	"fmt"
	"github.com/davidalpert/go-deep-merge/internal/version"
	"github.com/davidalpert/go-printers/v1"
	"github.com/spf13/cobra"
	"strings"
)

type VersionOptions struct {
	*printers.PrinterOptions
	VersionDetails *version.DetailStruct
}

func NewVersionOptions(ioStreams printers.IOStreams) *VersionOptions {
	return &VersionOptions{
		PrinterOptions: printers.NewPrinterOptions().WithStreams(ioStreams).WithDefaultOutput("text"),
		VersionDetails: &version.Detail,
	}
}

func NewCmdVersion(ioStreams printers.IOStreams) *cobra.Command {
	o := NewVersionOptions(ioStreams)
	var cmd = &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Args:  cobra.NoArgs,
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

	return cmd
}

// Complete the options
func (o *VersionOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}

// Validate the options
func (o *VersionOptions) Validate() error {
	return o.PrinterOptions.Validate()
}

// Run the command
func (o *VersionOptions) Run() error {
	if strings.EqualFold(*o.OutputFormat, "text") {
		s := fmt.Sprintf("%s %s", o.VersionDetails.AppName, o.VersionDetails.Version)
		_, err := fmt.Fprintln(o.Out, s)
		return err
	}
	if o.FormatCategory() == "table" || o.FormatCategory() == "csv" {
		o.OutputFormat = printers.StringPointer("json")
	}

	s, _, err := o.PrinterOptions.FormatOutput(o.VersionDetails)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(o.Out, s)
	return err
}
