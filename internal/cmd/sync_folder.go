package cmd

import (
	"fmt"
	"github.com/davidalpert/go-deep-merge/internal/app"
	v1 "github.com/davidalpert/go-deep-merge/v1"
	"github.com/davidalpert/go-printers/v1"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"sort"
	"strings"
)

type SyncFolderOptions struct {
	*printers.PrinterOptions
	SourceFolder string
	OutFolder    string
	OutFormat    string
	Debug        bool
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

type AppResult struct {
	AppDir      string
	MergeBySlug map[string]map[string]interface{}
}

// Run the command
func (o *SyncFolderOptions) Run() error {
	fis, err := afero.ReadDir(app.Fs, o.SourceFolder)
	if err != nil {
		return fmt.Errorf("reading source folder %#v: %#v", o.SourceFolder, err)
	}

	appDirs := make([]string, 0)
	for _, fi := range fis {
		if fi.IsDir() {
			appDirs = append(appDirs, path.Join(o.SourceFolder, fi.Name()))
		}
	}

	result := make([]AppResult, 0)
	for _, appDir := range appDirs {
		var defaultFile string
		var overrideFiles = make([]string, 0)

		fis, err = afero.ReadDir(app.Fs, appDir)
		for _, fi := range fis {
			name := fi.Name()
			if strings.HasSuffix(name, "default.yaml") {
				defaultFile = path.Join(appDir, name)
			} else if strings.HasSuffix(name, ".yaml") {
				overrideFiles = append(overrideFiles, path.Join(appDir, name))
			}
		}

		sort.Slice(overrideFiles, func(i, j int) bool {
			return len(overrideFiles[i]) < len(overrideFiles[j])
		})

		destFile, err := afero.ReadFile(app.Fs, defaultFile)
		if err != nil {
			return fmt.Errorf("read dest file %#v: %#v", defaultFile, err)
		}

		mergeResultBySlug := make(map[string]map[string]interface{})
		for _, override := range overrideFiles {
			var dest map[string]interface{}
			if err := yaml.Unmarshal(destFile, &dest); err != nil {
				return fmt.Errorf("unmarshalling dest: %#v", err)
			}

			slug := strings.TrimSuffix(path.Base(override), path.Ext(override))
			if strings.ContainsAny(slug, ".") {
				baseSlug := strings.Split(slug, ".")[0]
				// merge on top of another

				r, err := v1.MergeWithOptions(mergeResultBySlug[baseSlug], dest, v1.NewConfigDeeperMergeBang().WithMergeHashArrays(true).WithDebug(o.Debug))
				if err != nil {
					return fmt.Errorf("merging files %#v -> %#v: %#v", override, defaultFile, err)
				}

				dest = r
			}

			sourceFile, err := afero.ReadFile(app.Fs, override)
			if err != nil {
				return fmt.Errorf("read source file %#v: %#v", override, err)
			}

			var src map[string]interface{}
			if err := yaml.Unmarshal(sourceFile, &src); err != nil {
				return fmt.Errorf("unmarshalling src: %#v", err)
			}

			r, err := v1.MergeWithOptions(src, dest, v1.NewConfigDeeperMergeBang().WithMergeHashArrays(true).WithDebug(o.Debug))
			if err != nil {
				return fmt.Errorf("merging files %#v -> %#v: %#v", override, defaultFile, err)
			}

			mergeResultBySlug[slug] = r
		}

		result = append(result, AppResult{
			path.Base(appDir),
			mergeResultBySlug,
		})
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
