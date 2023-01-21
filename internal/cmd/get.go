package cmd

import (
	"fmt"
	"github.com/davidalpert/go-deep-merge/internal/cfgset"
	"github.com/davidalpert/go-deep-merge/internal/provider"
	"github.com/davidalpert/go-printers/v1"
	"github.com/spf13/cobra"
)

type GetOptions struct {
	*printers.PrinterOptions
	provider.Options
	Key       string
	Recursive bool
	//DecryptResult bool
}

func NewAwsGetOptions(ioStreams printers.IOStreams) *GetOptions {
	return &GetOptions{
		PrinterOptions: printers.NewPrinterOptions().WithStreams(ioStreams).WithDefaultOutput("text"),
	}
}

func NewCmdGet(ioStreams printers.IOStreams) *cobra.Command {
	o := NewAwsGetOptions(ioStreams)
	var cmd = &cobra.Command{
		Use:     "get <provider> <path>",
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
	o.Options.AddProviderOptions(cmd.Flags())
	cmd.Flags().BoolVarP(&o.Debug, "debug", "d", false, "enable debug output")
	cmd.Flags().BoolVarP(&o.Recursive, "recursive", "r", false, "recursively get all values under that path")

	return cmd
}

// Complete the options
func (o *GetOptions) Complete(cmd *cobra.Command, args []string) error {
	o.ProviderName = args[0]
	if p, err := provider.New(o.Options); err != nil {
		return fmt.Errorf("building provider: %#v", err)
	} else {
		o.Provider = p
	}

	if len(args) > 1 {
		o.Key = args[1]
	} else {
		o.Key = "/"
	}

	return nil
}

// Validate the options
func (o *GetOptions) Validate() error {
	return o.PrinterOptions.Validate()
}

// Run the command
func (o *GetOptions) Run() error {
	flattened := map[string]string{}
	if o.Recursive {
		if many, err := o.Provider.GetValueTree(o.Key); err != nil {
			return err
		} else {
			for k, v := range many {
				flattened[k] = v
			}
		}
	} else {
		if one, err := o.Provider.GetValue(o.Key); err != nil {
			return err
		} else {
			flattened[o.Key] = one
		}
	}

	if o.FormatCategory() == "text" {
		return o.WriteOutput(cfgset.FlattenedToString(flattened))
	}
	return o.WriteOutput(flattened)
}
