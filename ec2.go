package main

import (
	"github.com/awslabs/aws-sdk-go/gen/ec2"
)

type EC2 struct {
	*ec2.EC2
}

func NewEC2(a *AWS) *EC2 {
	return &EC2{ec2.New((*a).Credentials, (*a).Region, nil)}
}

func (service *EC2) GetAutoScalingGroup(instanceId string) (string, error) {
	request := &ec2.DescribeTagsRequest{
		Filters: []ec2.Filter{
			{Name: awsString("resource-id"), Values: []string{instanceId}},
		},
	}

	response, err := (*service).DescribeTags(request)
	if err != nil {
		return "", err
	}

	for _, tag := range response.Tags {
		if *tag.Key == "aws:autoscaling:groupName" {
			return *tag.Value, nil
		}
	}
	return "", nil
}
