package main

import (
	"github.com/awslabs/aws-sdk-go/aws"
)

type AWS struct {
	AccessId    string
	SecretKey   string
	Region      string
	Credentials aws.CredentialsProvider
}

func (a *AWS) Auth() {
	a.Credentials = aws.DetectCreds(a.AccessId, a.SecretKey, "")
}

func awsString(s string) aws.StringValue {
	return aws.String(s)
}

func awsDouble(f float64) aws.DoubleValue {
	return aws.Double(f)
}
