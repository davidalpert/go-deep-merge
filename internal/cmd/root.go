package cmd

import (
	"fmt"
	"github.com/davidalpert/go-deep-merge/internal/version"
	"github.com/davidalpert/go-printers/v1"
	"os"

	"github.com/spf13/cobra"
)

// cfgFile is an optional path to a configuration file used to initialize viper
var cfgFile string

// Execute builds the default root command and invokes it with os.Args
func Execute() {
	rootCmd := NewRootCmd(printers.DefaultOSStreams())

	rootCmd.SetArgs(os.Args[1:]) // without program

	err := rootCmd.Execute()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// RootCmdOptions is a struct to support version command
type RootCmdOptions struct {
	printers.IOStreams
	Verbose bool
}

// NewRootCmdOptions returns initialized RootCmdOptions
func NewRootCmdOptions(ioStreams printers.IOStreams) *RootCmdOptions {
	return &RootCmdOptions{
		IOStreams: ioStreams,
	}
}

// NewRootCmd creates the 'root' command and configures it's nested children
func NewRootCmd(ioStreams printers.IOStreams) *cobra.Command {
	//o := NewRootCmdOptions(ioStreams)
	rootCmd := &cobra.Command{
		Use:           version.Detail.AppName,
		Short:         "A tool for managing, merging, and shipping config files.",
		Long:          ``,
		SilenceUsage:  true,
		SilenceErrors: true,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		//	Run: func(cmd *cobra.Command, args []string) { },}
		Aliases: []string{},
		//RunE: func(cmd *cobra.Command, args []string) error {
		//	if err := o.Complete(cmd, args); err != nil {
		//		return err
		//	}
		//	if err := o.Validate(); err != nil {
		//		return err
		//	}
		//	return o.Run()
		//},
	}

	// Register subcommands
	rootCmd.AddCommand(NewCmdAws(ioStreams))
	rootCmd.AddCommand(NewCmdMerge(ioStreams))
	rootCmd.AddCommand(NewCmdSync(ioStreams))
	rootCmd.AddCommand(NewCmdVersion(ioStreams))

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", fmt.Sprintf("config file (default is $HOME/.%s/config.yaml)", version.Detail.AppName))
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose output")

	return rootCmd
}

// Complete the options
func (o *RootCmdOptions) Complete(cmd *cobra.Command, args []string) error {
	verbose, _ := cmd.PersistentFlags().GetBool("verbose")
	o.Verbose = verbose
	return nil
}

// Validate the options
func (o *RootCmdOptions) Validate() error {
	return nil
}

// Run the command
func (o *RootCmdOptions) Run() error {
	return nil
}
