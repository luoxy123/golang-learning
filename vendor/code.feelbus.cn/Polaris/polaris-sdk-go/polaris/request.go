package polaris

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/FeiniuBus/signer"
)

// Operation is the service api operation to be made
type Operation struct {
	HTTPMethod string
	HTTPPath   string
}

// Request is
type Request struct {
	config *Config
	method string
	path   string
	url    *url.URL
	params url.Values
	body   io.ReadSeeker
	header http.Header
	obj    interface{}
	ctx    context.Context
}

// Send is
func (r *Request) Send() ([]byte, error) {
	req, err := r.toHTTP()
	if err != nil {
		return nil, err
	}

	resp, err := r.config.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return make([]byte, 0), nil
	}

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(resp.Body)
		defer reader.Close()
	default:
		reader = resp.Body
	}

	var buf bytes.Buffer
	io.Copy(&buf, reader)

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Unexpected response code: %d (%s)", resp.StatusCode, buf.Bytes())
	}

	return buf.Bytes(), nil
}

func (r *Request) buildURI() {
	u, _ := url.Parse(StringValue(r.config.Address))
	u.Path = path.Join(u.Path, r.config.serviceName, r.path)

	r.url, _ = url.Parse(u.String())
}

func (r *Request) toHTTP() (*http.Request, error) {
	r.buildURI()
	if r.params != nil && len(r.params) > 0 {
		r.url.RawQuery = r.params.Encode()
	}
	if r.body == nil && r.obj != nil {
		b, err := json.Marshal(r.obj)
		if err != nil {
			return nil, err
		}
		r.header.Set("Content-Type", "application/json")
		rs := bytes.NewReader(b)
		r.body = rs
	}

	sr := &signer.Request{
		Body:   r.body,
		Method: r.method,
		Header: r.header,
		URL:    r.url,
	}

	creds, err := r.config.Credentials.Get()
	if err != nil {
		return nil, err
	}
	s := signer.NewHMACSignerV1(creds.AccessKey, creds.Secret)
	res := s.Sign(sr, 5*time.Second)

	for k, v := range res.Header {
		if len(v) == 0 {
			continue
		}
		r.header.Set(k, v[0])
	}

	req, err := http.NewRequest(r.method, r.url.String(), r.body)
	if err != nil {
		return nil, err
	}

	req.URL.Host = r.url.Host
	req.URL.Scheme = r.url.Scheme
	req.Host = r.url.Host
	req.Header = r.header

	if r.ctx != nil {
		return req.WithContext(r.ctx), nil
	}

	return req, nil
}
