package main

import (
	"encoding/json"
)

type Metric interface {
	FetchData() error
	GetData() map[string]interface{}
}

// Spring's metric endpoint
type Metrics struct {
	Data map[string]interface{}
	Url  string
	User string
	Pass string
}

func (m *Metrics) FetchData() error {
	content, err := HttpGet(m.Url, m.User, m.Pass)
	if err != nil {
		return err
	}

	return json.Unmarshal(content, &m.Data)
}

func (m *Metrics) GetData() map[string]interface{} {
	return m.Data
}
