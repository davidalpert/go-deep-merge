package cmd

import (
	"fmt"
	"github.com/davidalpert/go-deep-merge/internal/app"
	"github.com/davidalpert/go-deep-merge/internal/cfgset"
	"github.com/davidalpert/go-printers/v1"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"path"
)

type SyncFolderOptions struct {
	*printers.PrinterOptions
	cfgset.MergeOptions
	OutFolder string
	OutFormat string
}

func NewSyncFolderOptions(ioStreams printers.IOStreams) *SyncFolderOptions {
	return &SyncFolderOptions{
		PrinterOptions: printers.NewPrinterOptions().WithStreams(ioStreams).WithDefaultOutput("text"),
		OutFormat:      "yaml",
	}
}

func NewCmdSyncFolder(ioStreams printers.IOStreams) *cobra.Command {
	o := NewSyncFolderOptions(ioStreams)
	var cmd = &cobra.Command{
		Use:     "folder <source_folder>",
		Short:   "merge two config files together",
		Aliases: []string{"f", "fs"},
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

	cmd.Flags().BoolVarP(&o.Debug, "debug", "d", false, "enable debug output")
	cmd.Flags().StringVar(&o.OutFolder, "out-folder", "out", "folder to place output")
	//cmd.Flags().StringVar(&o.OutFormat, "out-format", "out", "format for output")

	return cmd
}

// Complete the options
func (o *SyncFolderOptions) Complete(cmd *cobra.Command, args []string) error {
	o.SourceFolder = args[0]
	return nil
}

// Validate the options
func (o *SyncFolderOptions) Validate() error {
	return o.PrinterOptions.Validate()
}

// Run the command
func (o *SyncFolderOptions) Run() error {
	result, err := cfgset.Merge(o.MergeOptions)
	if err != nil {
		return err
	}

	if err = app.Fs.MkdirAll(o.OutFolder, os.ModePerm); err != nil {
		return fmt.Errorf("making %#v: %#v", o.OutFolder, err)
	}

	for _, appResult := range result {
		appOutDir := path.Join(o.OutFolder, appResult.AppDir)
		if err = app.Fs.MkdirAll(appOutDir, os.ModePerm); err != nil {
			return fmt.Errorf("making %#v: %#v", appOutDir, err)
		}

		for slug, mergeResult := range appResult.MergeBySlug {
			outFile := path.Join(appOutDir, fmt.Sprintf("%s.%s", slug, o.OutFormat))
			b, err := yaml.Marshal(mergeResult)
			if err != nil {
				return fmt.Errorf("marshalling %#v to %#v: %#v", mergeResult, outFile, err)
			}

			if err = afero.WriteFile(app.Fs, outFile, b, os.ModePerm); err != nil {
				return fmt.Errorf("writing %#v: %#v", outFile, err)
			}

			// TODO: collect errors into an error result rather than failing out on the first one and write to STDERR
		}
	}

	//return o.WithDefaultOutput("json").WriteOutput(result)
	return nil
}
