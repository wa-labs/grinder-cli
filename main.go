package grindercli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Client ...
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient ...
func NewClient(baseURL string) (*Client, error) {

	client := &Client{
		httpClient: &http.Client{},
		baseURL:    baseURL,
	}
	return client, nil
}

func (c *Client) request(method, path string, content io.Reader, headers map[string]string) ([]byte, error) {
	fullPath := fmt.Sprintf("%s%s", c.baseURL, path)
	url, err := c.buildURL(fullPath)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, url.String(), content)
	if err != nil {
		return nil, err
	}

	c.setDefaultHeaders(request)
	for k, h := range headers {
		request.Header.Set(k, h)
	}
	response, err := c.httpClient.Do(request)

	if err != nil {
		return nil, err
	}

	fmt.Printf("Status : %d \n", response.StatusCode)

	if response.StatusCode >= 400 && response.StatusCode < 600 {
		return nil, errors.New("HTTP Status " + response.Status)
	}

	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()

	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) buildURL(pathOrURL string) (*url.URL, error) {
	fmt.Println(pathOrURL)
	u, err := url.ParseRequestURI(pathOrURL)
	if err != nil {
		u, err = url.Parse(c.baseURL)
		if err != nil {
			return nil, err
		}

		return u.Parse(pathOrURL)
	}

	return u, nil
}

func (c *Client) setDefaultHeaders(request *http.Request) {
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
}

// Get ...
func (c *Client) Get(path string, readOnlyDb bool) ([]byte, error) {
	var h map[string]string
	if !readOnlyDb {
		h["Read-Only-DB"] = "false"
	}
	return c.request("GET", path, nil, h)
}

// Post ...
func (c *Client) Post(path string, obj interface{}) ([]byte, error) {
	payloadJSON, _ := json.Marshal(obj)
	return c.request("POST", path, bytes.NewBuffer(payloadJSON), nil)
}

// Put ...
func (c *Client) Put(path string, obj interface{}) ([]byte, error) {
	payloadJSON, _ := json.Marshal(obj)
	return c.request("PUT", path, bytes.NewBuffer(payloadJSON), nil)
}

// Delete ...
func (c *Client) Delete(path string) ([]byte, error) {
	return c.request("DELETE", path, nil, nil)
}
