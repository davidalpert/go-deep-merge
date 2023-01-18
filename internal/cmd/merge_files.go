package cmd

import (
	"fmt"
	"github.com/davidalpert/go-deep-merge/internal/app"
	v1 "github.com/davidalpert/go-deep-merge/v1"
	"github.com/davidalpert/go-printers/v1"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type MergeFilesOptions struct {
	*printers.PrinterOptions
	Source      []byte
	Destination []byte
	Debug       bool
}

func NewMergeFilesOptions(ioStreams printers.IOStreams) *MergeFilesOptions {
	return &MergeFilesOptions{
		PrinterOptions: printers.NewPrinterOptions().WithStreams(ioStreams).WithDefaultOutput("text"),
	}
}

func NewCmdMergeFiles(ioStreams printers.IOStreams) *cobra.Command {
	o := NewMergeFilesOptions(ioStreams)
	var cmd = &cobra.Command{
		Use:     "files <src_file> <dest_file>",
		Short:   "merge two config files together",
		Aliases: []string{"f", "fs", "file"},
		Args:    cobra.ExactArgs(2),
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
	cmd.Flags().BoolVarP(&o.Debug, "debug", "d", false, "enable debug output")

	return cmd
}

// Complete the options
func (o *MergeFilesOptions) Complete(cmd *cobra.Command, args []string) error {
	if b, err := afero.ReadFile(app.Fs, args[0]); err != nil {
		return fmt.Errorf("reading %#v: %#v", args[0], err)
	} else {
		o.Source = b
	}
	if b, err := afero.ReadFile(app.Fs, args[1]); err != nil {
		return fmt.Errorf("reading %#v: %#v", args[1], err)
	} else {
		o.Destination = b
	}
	return nil
}

// Validate the options
func (o *MergeFilesOptions) Validate() error {
	return o.PrinterOptions.Validate()
}

// Run the command
func (o *MergeFilesOptions) Run() error {
	var src, dest map[string]interface{}
	if err := yaml.Unmarshal(o.Source, &src); err != nil {
		return fmt.Errorf("unmarshalling src: %#v", err)
	}

	if err := yaml.Unmarshal(o.Destination, &dest); err != nil {
		return fmt.Errorf("unmarshalling dest: %#v", err)
	}

	//fmt.Fprintln(o.Out, "source:")
	//o.WriteOutput(src)
	//fmt.Fprintln(o.Out, "dest:")
	//o.WriteOutput(dest)
	//fmt.Fprintln(o.Out, "result:")

	r, err := v1.MergeWithOptions(src, dest, v1.NewConfigDeeperMergeBang().WithMergeHashArrays(true).WithDebug(o.Debug))
	if err != nil {
		return fmt.Errorf("merging files: %#v", err)
	}

	return o.WriteOutput(r)
}
