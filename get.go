package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func HttpGet(url, user, pass string) ([]byte, error) {
	if DEBUG {
		log.Printf("Making a request to %s\n", url)
	}

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if user != "" && pass != "" {
		req.SetBasicAuth(user, pass)
	}

	response, err := client.Do(req)

	if err != nil {
		if DEBUG && response != nil {
			log.Println("Error: ", err)
			log.Println("Response: ", response.Header)
			response.Body.Close()
		}
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		if DEBUG {
			log.Println("Unmarshalling error:", err)
		}
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, errors.New(string(body))
	}

	return body, nil
}
