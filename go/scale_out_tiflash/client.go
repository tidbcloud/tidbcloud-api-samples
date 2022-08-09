// Copyright(C) 2022 PingCAP. All Rights Reserved.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/icholy/digest"
)

var (
	clientInitOnce sync.Once

	restClient *resty.Client
)

func initClient(publicKey, privateKey string) error {
	clientInitOnce.Do(func() {
		restClient = resty.New()
		restClient.SetTransport(&digest.Transport{
			Username: publicKey,
			Password: privateKey,
		})
	})

	return nil
}

func waitFor(scene string, fn func() bool) {
	ticker := time.NewTicker(time.Second)

	for range ticker.C {
		if fn() {
			return
		}

		log.Printf("still waitting %s to get ready...", scene)
	}
}

// doRequest wraps resty request, it's a generic method to spawn a HTTP request
func doRequest(method, url string, payload, output interface{}) (*resty.Response, error) {
	request := restClient.R()

	// if payload is not nil, we'll put it on body
	if payload != nil {
		request.SetBody(payload)
	}

	// execute the request
	resp, err := request.Execute(method, url)
	b, _ := json.Marshal(payload)
	fmt.Printf("\npayload: %s\n", b)
	fmt.Printf("\nRequest: method %s, url %s, response %s\n\n", method, url, resp)
	if err != nil {
		return nil, err
	}

	// if the request return a non-200 response, wrap it with error
	if resp.StatusCode() != http.StatusOK {
		return resp, fmt.Errorf("Failed with status %d and resp %s", resp.StatusCode(), resp)
	}

	// if we need to unmarshal the response into a struct, we pass it here, otherwise pass nil in the argument
	if output != nil {
		return resp, json.Unmarshal(resp.Body(), output)
	}

	return resp, nil
}

func doGET(url string, payload, output interface{}) (*resty.Response, error) {
	return doRequest(resty.MethodGet, url, payload, output)
}

func doPOST(url string, payload, output interface{}) (*resty.Response, error) {
	return doRequest(resty.MethodPost, url, payload, output)
}

func doDELETE(url string, payload, output interface{}) (*resty.Response, error) {
	return doRequest(resty.MethodDelete, url, payload, output)
}

func doPATCH(url string, payload, output interface{}) (*resty.Response, error) {
	return doRequest(resty.MethodPatch, url, payload, output)
}
