package polaris

import (
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"strings"
)

// Client provides a client to the service api
type Client struct {
	Config *Config
}

// NewRequest returns a new Request pointer for the service api
func (c *Client) NewRequest(operation *Operation, params url.Values, data interface{}) (*Request, error) {
	method := operation.HTTPMethod
	if method == "" {
		method = "GET"
	}

	r := &Request{
		config: c.Config.Copy(),
		method: method,
		params: make(map[string][]string),
		obj:    data,
		header: make(http.Header),
		path:   operation.HTTPPath,
	}

	if params != nil && len(params) > 0 {
		r.params = params
	}

	r.header.Set("Accept", "application/json")
	r.header.Set("Accept-Encoding", "gzip")

	ua := buildUserAgent("polaris-sdk-go", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	r.header.Set("User-Agent", ua)

	return r, nil
}

// NewClient returns a pointer to a new service client.
func NewClient(cfg *Config) *Client {
	svc := &Client{
		Config: cfg,
	}

	return svc
}

func buildUserAgent(sdkname string, extra ...string) string {
	ua := sdkname
	if len(extra) > 0 {
		ua += fmt.Sprintf(" (%s)", strings.Join(extra, "; "))
	}

	return ua
}
