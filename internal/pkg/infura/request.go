package infura

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type Client struct {
	url        string
	projectID  string
	privateKey string
	httpClient *http.Client
}

func New(projectID string, privateKey string) *Client {
	return &Client{
		url:        "https://mainnet.infura.io/v3/" + projectID,
		projectID:  projectID,
		privateKey: privateKey,
		httpClient: &http.Client{},
	}
}

func (c *Client) Request(json string) (string, error) {
	req, err := http.NewRequest("POST", c.url, bytes.NewBufferString(json))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth("", c.privateKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
