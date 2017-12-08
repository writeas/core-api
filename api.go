// Package as provides basic functions for any *.as API library
package as

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/writeas/impart"
	"io"
	"net/http"
)

type Client struct {
	Config *ClientConfig
	Token  string
}

func NewClient(cfg *ClientConfig) *Client {
	return &Client{
		Config: cfg,
		Token:  "",
	}
}

// SetToken sets the user token for all future Client requests. Setting this to
// an empty string will change back to unauthenticated requests.
func (c *Client) SetToken(token string) {
	c.Token = token
}

func (c *Client) Get(path string, r interface{}) (*impart.Envelope, error) {
	method := "GET"
	if method != "GET" && method != "HEAD" {
		return nil, fmt.Errorf("Method %s not currently supported by library (only HEAD and GET).\n", method)
	}

	return c.request(method, path, nil, r)
}

func (c *Client) Post(path string, data, r interface{}) (*impart.Envelope, error) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(data)
	return c.request("POST", path, b, r)
}

func (c *Client) Put(path string, data, r interface{}) (*impart.Envelope, error) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(data)
	return c.request("PUT", path, b, r)
}

func (c *Client) Delete(path string, data map[string]string) (*impart.Envelope, error) {
	r, err := c.buildRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	q := r.URL.Query()
	for k, v := range data {
		q.Add(k, v)
	}
	r.URL.RawQuery = q.Encode()

	return c.doRequest(r, nil)
}

func (c *Client) request(method, path string, data io.Reader, result interface{}) (*impart.Envelope, error) {
	r, err := c.buildRequest(method, path, data)
	if err != nil {
		return nil, err
	}

	return c.doRequest(r, result)
}

func (c *Client) buildRequest(method, path string, data io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", c.Config.BaseURL, path)
	r, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, fmt.Errorf("Create request: %v", err)
	}
	c.prepareRequest(r)

	return r, nil
}

func (c *Client) doRequest(r *http.Request, result interface{}) (*impart.Envelope, error) {
	resp, err := c.Config.Client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("Request: %v", err)
	}
	defer resp.Body.Close()

	env := &impart.Envelope{
		Code: resp.StatusCode,
	}
	if result != nil {
		env.Data = result

		err = json.NewDecoder(resp.Body).Decode(&env)
		if err != nil {
			return nil, err
		}
	}

	return env, nil
}

func (c *Client) prepareRequest(r *http.Request) {
	r.Header.Add("User-Agent", "go-writeas v1")
	r.Header.Add("Content-Type", "application/json")
	if c.Token != "" {
		r.Header.Add("Authorization", "Token "+c.Token)
	}
}
