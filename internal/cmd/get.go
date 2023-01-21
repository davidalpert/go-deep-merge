package cmd

import (
	"fmt"
	"github.com/davidalpert/go-deep-merge/internal/app"
	"github.com/davidalpert/go-deep-merge/internal/paramstore"
	"github.com/davidalpert/go-printers/v1"
	"github.com/spf13/cobra"
	"sort"
	"strings"
)

type GetOptions struct {
	*printers.PrinterOptions
	Client       app.ConfigProvider
	ProviderName string
	Key          string
	Debug        bool
	Recursive    bool
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
	//cmd.Flags().BoolVar(&o.Decrypt, "decrypt", false, "decrypt result?")
	cmd.Flags().BoolVarP(&o.Debug, "debug", "d", false, "debug")
	cmd.Flags().BoolVarP(&o.Recursive, "recursive", "r", false, "recursively get all values under that path")

	return cmd
}

// Complete the options
func (o *GetOptions) Complete(cmd *cobra.Command, args []string) error {
	o.ProviderName = args[0]

	// providers-by-key pattern
	//supportedProviders := map[string]func() app.ConfigProvider{
	//	"aws": func() app.ConfigProvider { return paramstore.NewSSMClient(o.Debug) },
	//}
	//
	//supportedProviderKeys := make([]string, 0)
	//for k, _ := range supportedProviders {
	//	supportedProviderKeys = append(supportedProviderKeys, k)
	//}
	//sort.Strings(supportedProviderKeys)
	//
	//if clientFunc, ok := supportedProviders[o.ProviderName]; ok {
	//	o.Client = clientFunc()
	//} else {
	//	return fmt.Errorf("unrecognized provider %#v; supported providers are: %#v", o.ProviderName, strings.Join(supportedProviderKeys, ", "))
	//}

	// providers-as-list pattern
	supportedProviderNames := []string{"aws"}
	sort.Strings(supportedProviderNames)
	if strings.EqualFold(o.ProviderName, "aws") {
		o.Client = paramstore.NewSSMClient(o.Debug)
	} else {
		return fmt.Errorf("unrecognized provider %#v; supported providers are: %#v", o.ProviderName, strings.Join(supportedProviderNames, ", "))
	}
	o.Key = args[1]

	return nil
}

// Validate the options
func (o *GetOptions) Validate() error {
	return o.PrinterOptions.Validate()
}

// Run the command
func (o *GetOptions) Run() error {
	if o.Recursive {
		return o.getMany()
	}
	return o.getOne()
}

func (o *GetOptions) getMany() error {
	result, err := o.Client.GetValueTree(o.Key)

	if err != nil {
		return fmt.Errorf("get many %s %#v: %#v", o.ProviderName, o.Key, err)
	}

	return o.WriteOutput(result)
}

func (o *GetOptions) getOne() error {
	result, err := o.Client.GetValue(o.Key)
	if err != nil {
		return fmt.Errorf("get %s %#v: %#v", o.ProviderName, o.Key, err)
	}

	return o.WriteOutput(result)
}
