package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

type SSM struct {
	client ssmiface.SSMAPI
}

type Param struct {
	Name           string
	WithDecryption bool
	ssmsvc         *SSM
}

func CreateAwsSession() (*session.Session, error) {

	// Initialize a session in <region> that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.

	sess, err := session.NewSessionWithOptions(session.Options{
		// Specify profile to load for the session's config
		Profile: "default",

		// Provide SDK Config options, such as Region.
		Config: aws.Config{
			Region: aws.String("eu-west-1"),
		},

		// Force enable Shared Config support
		SharedConfigState: session.SharedConfigEnable,
	})

	return sess, err

}

func NewSSMClient() *SSM {
	// Create AWS Session
	sess, err := CreateAwsSession()
	if err != nil {
		log.Println(err)
		return nil
	}
	ssmsvc := &SSM{ssm.New(sess)}
	// Return SSM client
	return ssmsvc
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

func main() {
	ssmsvc := NewSSMClient()
	usr, err := ssmsvc.Param("/babyshark/hardenize/user", true).GetValue()
	pass, err := ssmsvc.Param("/babyshark/hardenize/pass", true).GetValue()
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("User: %s\nPass: %s\n", usr, pass)
}
