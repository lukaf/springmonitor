package main

import (
	"github.com/awslabs/aws-sdk-go/gen/cloudwatch"
)

type CloudWatch struct {
	*cloudwatch.CloudWatch
}

const metricsLimit = 20

func NewCloudWatch(a *AWS) *CloudWatch {
	return &CloudWatch{cloudwatch.New(a.Credentials, a.Region, nil)}
}

func BuildMetrics(m Metric, unit, namespace, instanceId, autoScalingGroup string) []*cloudwatch.PutMetricDataInput {

	var bucket []*cloudwatch.PutMetricDataInput
	var index int = 0

	dimensions := []cloudwatch.Dimension{
		{Name: awsString("InstanceId"), Value: awsString(instanceId)},
	}

	if autoScalingGroup != "" {
		dimensions = append(
			dimensions,
			cloudwatch.Dimension{
				Name: awsString("AutoScalingGroupName"), Value: awsString(autoScalingGroup),
			},
		)
	}

	metrics := []cloudwatch.MetricDatum{}
	for key, value := range m.GetData() {
		metrics = append(
			metrics,
			cloudwatch.MetricDatum{
				MetricName: awsString(key),
				Unit:       awsString(unit),
				Value:      awsDouble(value.(float64)),
				Dimensions: dimensions,
			},
		)
	}

	for ; index+metricsLimit <= len(metrics); index += metricsLimit {
		bucket = append(
			bucket,
			&cloudwatch.PutMetricDataInput{
				Namespace:  awsString(namespace),
				MetricData: metrics[index : index+metricsLimit],
			},
		)
	}

	if index <= len(metrics)-1 {
		bucket = append(
			bucket,
			&cloudwatch.PutMetricDataInput{
				Namespace:  awsString(namespace),
				MetricData: metrics[index:],
			},
		)
	}

	return bucket
}
