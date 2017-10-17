package service

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// HTTPClient is an interface of  client HTTP service.
type HTTPClient interface {
	// Post that performs sending of POST request to HTTP server.
	Post(method string,
		args map[string]interface{}) (map[string]interface{}, error)
}

// Client is an implementation of HTTPClient interface.
type Client struct {
	OpenViDuURL string
	Login       string
	Password    string
}

// Post sends HTTP post request to HTTP server.
func (c *Client) Post(
	method string,
	args map[string]interface{}) (map[string]interface{}, error) {
	var requestData io.Reader
	if args != nil {
		rawMessage, err := json.Marshal(args)
		if err != nil {
			return nil, err
		}
		requestData = bytes.NewBuffer(rawMessage)
	}
	req, err := http.NewRequest(
		http.MethodPost, fmt.Sprintf(
			"%s/%s", c.OpenViDuURL, method), requestData)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.Login, c.Password)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	if args != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
