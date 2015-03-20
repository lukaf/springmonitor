package main

import (
	"flag"
	"log"
)

var DEBUG bool = true

func main() {

	var url = flag.String("url", "http://localhost:8080/metrics", "Full URL for metrics.")
	var accessId = flag.String("id", "", "AWS access id.")
	var secretKey = flag.String("key", "", "AWS secret key.")
	var region = flag.String("region", "eu-west-1", "AWS region.")
	var namespace = flag.String("namespace", "", "CloudWatch namespace")
	var user = flag.String("user", "", "Username for BASIC auth.")
	var pass = flag.String("pass", "", "Password for BASIC auth.")
	flag.Parse()

	metrics := &Metrics{
		Url:  *url,
		User: *user,
		Pass: *pass,
	}

	if err := metrics.FetchData(); err != nil {
		log.Printf("Failed to fetch data: ", err)
		return
	}

	aws := &AWS{
		AccessId:  *accessId,
		SecretKey: *secretKey,
		Region:    *region,
	}

	log.Println("Fetching credentials.")
	aws.Auth()

	log.Println("Connecting to AWS services ...")
	ec2 := NewEC2(aws)
	cw := NewCloudWatch(aws)

	log.Println("Fetching instance ID.")
	instanceId, err := InstanceId()
	if err != nil {
		log.Println("Failed to fetch instance ID:", err)
	}

	log.Println("Fetching auto scaling group name.")
	scalingGroup, err := ec2.GetAutoScalingGroup(instanceId)
	if err != nil {
		log.Println("Failed to fetch auto scaling group name:", err)
	}

	log.Println("Building metrics.")
	datas := BuildMetrics(metrics, "Count", *namespace, instanceId, scalingGroup)

	for index, data := range datas {
		log.Println("Sending metrics batch: ", index)
		r := cw.PutMetricData(data)
		if r != nil {
			log.Println("Unexpected response while sending metrics:", r)
		}
	}

	log.Println("Done.")
}
