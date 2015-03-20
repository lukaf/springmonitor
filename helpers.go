package main

func InstanceId() (string, error) {
	id, err := HttpGet("http://169.254.169.254/latest/meta-data/instance-id", "", "")
	if err != nil {
		return "NA", err
	}
	return string(id), err
}
