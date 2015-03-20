package main

import (
	"os"
	"testing"
)

const (
	ec2AccessId  = "ec2AccessId"
	ec2SecretKey = "ec2SecretKey"
	ec2Region    = "eu-west-1"
)

func TestNewEC2(t *testing.T) {
	aws := &AWS{Region: ec2Region}

	if err := os.Setenv("AWS_ACCESS_KEY_ID", ec2AccessId); err != nil {
		t.Errorf("Unable to set AWS_ACCESS_KEY_ID")
	}

	if err := os.Setenv("AWS_SECRET_ACCESS_KEY", ec2SecretKey); err != nil {
		t.Errorf("Unable to set AWS_SECRET_ACCESS_KEY")
	}

	var _ *EC2 = NewEC2(aws)
}
