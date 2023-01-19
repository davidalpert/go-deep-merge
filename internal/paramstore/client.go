package paramstore

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"log"
)

func Sessions(debug bool) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		CredentialsChainVerboseErrors: aws.Bool(debug),
	})
	svc := session.Must(sess, err)
	return svc, err
}

func NewSSMClient(debug bool) *Client {
	// Create AWS Session
	sess, err := Sessions(debug)
	if err != nil {
		log.Println(err)
		return nil
	}
	ssmsvc := &Client{ssm.New(sess)}
	// Return Client client
	return ssmsvc
}

// Client is a Client API client.
type Client struct {
	client ssmiface.SSMAPI
}

func (s *Client) GetValue(name string, decryption bool) (string, error) {
	ssmsvc := s.client
	parameter, err := ssmsvc.GetParameter(&ssm.GetParameterInput{
		Name:           &name,
		WithDecryption: &decryption,
	})
	if err != nil {
		return "", err
	}
	value := *parameter.Parameter.Value
	return value, nil
}
