package xmlrpc

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
)

// Client implements a basic XMLRPC client
type Client struct {
	addr       string
	httpClient *http.Client
}

// NewClient returns a new instance of Client
// Pass in a true value for `insecure` to turn off certificate verification
func NewClient(addr string, insecure bool) *Client {
	transport := &http.Transport{}
	if insecure {
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	httpClient := &http.Client{Transport: transport}

	return &Client{
		addr:       addr,
		httpClient: httpClient,
	}
}

// Call calls the method with "name" with the given args
// Returns the result, and an error for communication errors
func (c *Client) Call(name string, args ...interface{}) (interface{}, error) {
	req := bytes.NewBuffer(nil)
	e := Marshal(req, name, args...)
	if e != nil {
		return nil, e
	}
	r, e := c.httpClient.Post(c.addr, "text/xml", req)
	if e != nil {
		return nil, e
	}
	defer r.Body.Close()

	_, v, f, e := Unmarshal(r.Body)
	if f != nil {
		e = fmt.Errorf("Error: %v: %v", e, f)
	}
	return v, e
}
