// Package rest provides a simple HTTP client
package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	constants "github.com/toolboxconstants"
)

// Client to call 3rd party APIs
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new Client with the supplied Base URL and timeoutSec
// if timeoutSec is less than 1, it will be set to constants.DefaultHTTPClientTimeoutSec
func NewClient(baseURL string, timeoutSec int) *Client {
	if timeoutSec < 1 {
		timeoutSec = constants.DefaultHTTPClientTimeoutSec
	}
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: time.Second * time.Duration(timeoutSec),
		},
	}
}

// MakeNewRequest to make new request for calling. It gives response & error.
func (c *Client) MakeNewRequest(method string, body interface{}, queryParams map[string]string, headerSetParams map[string]string, headerAddParams map[string]string) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, c.BaseURL, buf)
	if err != nil {
		return nil, err
	}

	if headerSetParams != nil {
		h := req.Header
		for key, value := range headerSetParams {
			h.Set(key, value)
		}
	}

	if headerAddParams != nil {
		h := req.Header
		for key, value := range headerAddParams {
			h.Add(key, value)
		}
	}

	if queryParams != nil {
		q := req.URL.Query()
		for key, value := range queryParams {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	return req, nil
}

/*Do sends an HTTP request and returns an HTTP response, following policy (such as redirects, cookies, auth) as configured on the client.

It also decode the incoming HTTP response according to our response structure which address is provided in "v".

An error is returned if caused by client policy (such as CheckRedirect), or failure to speak HTTP (such as a network connectivity problem). A non-2xx status code doesn't cause an error.
*/
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
		return resp, err
	}
	return resp, nil
}
