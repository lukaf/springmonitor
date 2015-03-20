package main

import (
	"fmt"
	"os"
	//	"reflect"
	"testing"
)

const (
	cwAccessId           = "testId"
	cwSecretKey          = "testSecret"
	cwRegion             = "eu-west-1"
	testUnit             = "Count"
	testNameSpace        = "testNameSpace"
	testInstanceId       = "testInstanceId"
	testAutoScalingGroup = "testAutoScalingGroup"
	metricsCount         = 75
)

func TestNewCloudWatch(t *testing.T) {
	aws := &AWS{Region: cwRegion}

	if err := os.Setenv("AWS_ACCESS_KEY_ID", cwAccessId); err != nil {
		t.Errorf("Unable to set AWS_ACCESS_KEY_ID")
	}

	if err := os.Setenv("AWS_SECRET_ACCESS_KEY", cwSecretKey); err != nil {
		t.Errorf("Unable to set AWS_SECRET_ACCESS_KEY")
	}

	// just a type check, nothing more
	var _ *CloudWatch = NewCloudWatch(aws)
}

func TestBuildMetrics(t *testing.T) {
	testMetrics := &TestMetric{}
	testMetrics.FetchData()

	// metricsLimit is a constant defined in cloudwatch.go which limits the number of metric sent with one API call.
	// Limit is enforced by AWS.
	bucketsNumber := metricsCount / metricsLimit
	if metricsCount%metricsLimit > 0 {
		bucketsNumber += 1
	}

	buckets := BuildMetrics(testMetrics, testUnit, testNameSpace, testInstanceId, testAutoScalingGroup)

	if len(buckets) != bucketsNumber {
		t.Errorf("Number of returned buckets don't match: got %d expected %d\n", len(buckets), bucketsNumber)
	}

	for bucketIndex, bucket := range buckets {
		if len(bucket.MetricData) > metricsLimit {
			t.Errorf(
				"Too many metrics in bucket %d: limit is %d bucket contains %d\n",
				bucketIndex,
				metricsLimit,
				len(bucket.MetricData),
			)
		}

		if *bucket.Namespace != testNameSpace {
			t.Errorf("Wrong Namespace value in bucket %d: got %s expected %s", bucketIndex, *bucket.Namespace, testNameSpace)
		}

		for metricIndex, metric := range bucket.MetricData {
			if *metric.Unit != testUnit {
				t.Errorf("Wrong value for unit in metric %d in bucket %d: got %s expected %s\n", metricIndex, bucketIndex, *metric.Unit, testUnit)
			}

			if *metric.MetricName != fmt.Sprintf("%.0f", *metric.Value) {
				t.Errorf("Wrong value for metric %d in bucket %d: got %.0f expected %s\n", metricIndex, bucketIndex, *metric.Value, *metric.MetricName)
			}

			for _, dimension := range metric.Dimensions {
				if *dimension.Name == "AutoScalingGroupName" && *dimension.Value != testAutoScalingGroup {
					t.Errorf("Wrong value for AutoScalingGroupName in metric %d in bucket %d: got %s expected %s\n", metricIndex, bucketIndex, *dimension.Value, testAutoScalingGroup)
				}

				if *dimension.Name == "InstanceId" && *dimension.Value != testInstanceId {
					t.Errorf("Wrong value for InstanceId in metric %d in bucket %d: got %s expected %s\n", metricIndex, bucketIndex, *dimension.Value, testInstanceId)
				}
			}
		}
	}
}

// Metric implementation for testing
type TestMetric struct {
	Data map[string]interface{}
	Url  string
}

func (m *TestMetric) GetData() map[string]interface{} {
	return m.Data
}

func (m *TestMetric) FetchData() error {

	content := make(map[string]interface{})
	var x float64 = 0

	for ; x < metricsCount; x++ {
		content[fmt.Sprintf("%.0f", x)] = x
	}

	m.Data = content

	return nil
}
