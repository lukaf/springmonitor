package main

import (
	"os"
	"reflect"
	"testing"
)

const (
	testAccessId  = "AccessId"
	testSecretKey = "SecretKey"
	testRegion    = "eu-region"
)

func TestAWS(t *testing.T) {
	if err := os.Setenv("AWS_ACCESS_KEY_ID", testAccessId); err != nil {
		t.Errorf("Unable to set AWS_ACCESS_KEY_ID environment variable.\n")
	}

	if err := os.Setenv("AWS_SECRET_ACCESS_KEY", testSecretKey); err != nil {
		t.Errorf("Unable to set AWS_SECRET_ACCESS_KEY environment variable.\n")
	}

	services := &AWS{Region: testRegion}
	services.Auth()

	if reflect.TypeOf(services.Credentials).String() != "aws.staticCredentialsProvider" {
		t.Errorf("aws.Credentials of wrong type: got %s expected aws.Credentials\n", reflect.TypeOf(services.Credentials).String())
	}

	credentials, err := services.Credentials.Credentials()
	if err != nil {
		t.Errorf("No credentials could be provided: %s", err)
	}

	if credentials.AccessKeyID != testAccessId {
		t.Errorf("Wrong access id: got %s expected %s\n", credentials.AccessKeyID, testAccessId)
	}

	if credentials.SecretAccessKey != testSecretKey {
		t.Errorf("Wrong secret key: got %s expected %s\n", credentials.SecretAccessKey, testSecretKey)
	}
}

func TestAwsString(t *testing.T) {
	input := "TestAwsString"

	output := awsString(input)

	if reflect.TypeOf(output).String() != "aws.StringValue" {
		t.Errorf("output of wrong type: got %s expected aws.StringValue\n", reflect.TypeOf(output).String())
	}

	if *output != input {
		t.Errorf("Wrong value for ouput: got %s expected %s\n", *output, input)
	}
}

func TestAwsDouble(t *testing.T) {
	var input float64 = 10.0

	output := awsDouble(input)

	if reflect.TypeOf(output).String() != "aws.DoubleValue" {
		t.Errorf("output of wrong type: got %s expected aws.DoubleValue\n", reflect.TypeOf(output).String())
	}

	if *output != input {
		t.Errorf("Wrong value for output: got %f expected %f\n", *output, input)
	}
}
