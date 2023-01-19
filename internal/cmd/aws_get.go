package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/davidalpert/go-printers/v1"
	"github.com/spf13/cobra"
	"log"
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

func Sessions(debug bool) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		CredentialsChainVerboseErrors: aws.Bool(debug),
	})
	svc := session.Must(sess, err)
	return svc, err
}

// SSM is a SSM API client.
type SSM struct {
	client ssmiface.SSMAPI
}

func NewSSMClient(debug bool) *SSM {
	// Create AWS Session
	sess, err := Sessions(debug)
	if err != nil {
		log.Println(err)
		return nil
	}
	ssmsvc := &SSM{ssm.New(sess)}
	// Return SSM client
	return ssmsvc
}

type Param struct {
	Name           string
	WithDecryption bool
	ssmsvc         *SSM
}

//Param creates the struct for querying the param store
func (s *SSM) Param(name string, decryption bool) *Param {
	return &Param{
		Name:           name,
		WithDecryption: decryption,
		ssmsvc:         s,
	}
}

func (p *Param) GetValue() (string, error) {
	ssmsvc := p.ssmsvc.client
	parameter, err := ssmsvc.GetParameter(&ssm.GetParameterInput{
		Name:           &p.Name,
		WithDecryption: &p.WithDecryption,
	})
	if err != nil {
		return "", err
	}
	value := *parameter.Parameter.Value
	return value, nil
}

// Run the command
func (o *AwsGetOptions) Run() error {
	ssmsvc := NewSSMClient(o.Debug)
	result, err := ssmsvc.Param(o.ParameterName, o.Decrypt).GetValue()
	if err != nil {
		return fmt.Errorf("get param %#v: %#v", o.ParameterName, err)
	}

	return o.WriteOutput(result)
}
