package feiniubus

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strings"
	"time"
)

// Request is the service request to be made
type Request struct {
	Config       Config
	ClientInfo   ClientInfo
	Handlers     Handlers
	Time         time.Time
	ExpireTime   time.Duration
	Operation    *Operation
	HTTPRequest  *http.Request
	HTTPResponse *http.Response
	Body         io.ReadSeeker
	BodyStart    int64
	Params       map[string]string
	Content      interface{}
	Error        error
	Data         interface{}
	Unmarshaler  Unmarshaler

	built    bool
	safeBody *offsetReader
}

// Operation is the service API operation to be made
type Operation struct {
	Name           string
	HTTPMethod     string
	HTTPPath       string
	UseQueryString bool
	ContentType    string
	Content        interface{}
	Params         map[string]string
	Data           interface{}
	Unmarshaler    Unmarshaler
}

// New returns a new Request pointer for the service API
func New(cfg Config, clientInfo ClientInfo, handlers Handlers, operation *Operation) *Request {
	method := operation.HTTPMethod
	if method == "" {
		method = "GET"
	}

	httpReq, _ := http.NewRequest(method, "", nil)

	var err error
	u, _ := url.Parse(clientInfo.Endpoint)
	u.Path = path.Join(u.Path, operation.HTTPPath)
	httpReq.URL, err = url.Parse(u.String())
	if err != nil {
		httpReq.URL = &url.URL{}
	}

	r := &Request{
		Config:      cfg,
		ClientInfo:  clientInfo,
		Handlers:    handlers.Copy(),
		Time:        time.Now(),
		ExpireTime:  0,
		Operation:   operation,
		HTTPRequest: httpReq,
		Body:        nil,
		Params:      operation.Params,
		Content:     operation.Content,
		Error:       err,
		Data:        operation.Data,
		Unmarshaler: operation.Unmarshaler,
	}

	r.SetBufferBody([]byte{})

	return r
}

// UseQueryString returns if the request has no body
func (r *Request) UseQueryString() bool {
	if r.Operation.HTTPMethod != "" {
		if strings.ToUpper(r.Operation.HTTPMethod) == "GET" || strings.ToUpper(r.Operation.HTTPMethod) == "DELETE" {
			return true
		}
	}
	return r.Operation.UseQueryString
}

// ParamsFilled returns if the request's parameters have been populated
func (r *Request) ParamsFilled() bool {
	return r.Params != nil && len(r.Params) > 0
}

// ContentFilled returns if the request's content hava been populated
func (r *Request) ContentFilled() bool {
	if r.Content == nil {
		return false
	}

	k := reflect.TypeOf(r.Content).Kind()
	if k != reflect.Array && k != reflect.Slice {
		return reflect.ValueOf(r.Content).Elem().IsValid()
	}

	return true
}

// DataFilled returns true if the request's data for response has been set
func (r *Request) DataFilled() bool {
	return r.Data != nil && reflect.ValueOf(r.Data).Elem().IsValid()
}

// UnmarshalerFilled returns true if the request's Unmarshaler for response has been set
func (r *Request) UnmarshalerFilled() bool {
	return r.Unmarshaler != nil && reflect.ValueOf(r.Unmarshaler).Elem().IsValid()
}

// SetBufferBody will set the request's body bytes
func (r *Request) SetBufferBody(buf []byte) {
	r.SetReaderBody(bytes.NewReader(buf))
}

// SetStringBody sets the body of the request to be backed by a string
func (r *Request) SetStringBody(s string) {
	r.SetReaderBody(strings.NewReader(s))
}

// SetReaderBody will set the request's body reader
func (r *Request) SetReaderBody(reader io.ReadSeeker) {
	r.Body = reader
	r.ResetBody()
}

// Send will send the request returning error if errors are encountered.
func (r *Request) Send() error {
	r.Build()
	if r.Error != nil {
		return r.Error
	}
	r.Handlers.Send.Run(r)
	if r.Error != nil {
		return r.Error
	}

	r.Handlers.ValidateResponse.Run(r)
	if r.Error != nil {
		r.Handlers.UnmarshalError.Run(r)
		if r.Error != nil {
			return r.Error
		}

		return nil
	}

	r.Handlers.Unmarshal.Run(r)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

// Build will build the request's object
func (r *Request) Build() error {
	if !r.built {
		r.Handlers.Validate.Run(r)
		if r.Error != nil {
			return r.Error
		}
		r.Handlers.Build.Run(r)
		if r.Error != nil {
			return r.Error
		}
		r.built = true
	}
	return r.Error
}

// ResetBody rewinds the request body backto its starting position
func (r *Request) ResetBody() {
	if r.safeBody != nil {
		r.safeBody.Close()
	}

	r.safeBody = newOffsetReader(r.Body, r.BodyStart)

	l, err := computeBodyLength(r.Body)
	if err != nil {
		r.Error = err
		return
	}

	if l == 0 {
		r.HTTPRequest.Body = noBodyReader
	} else if l > 0 {
		r.HTTPRequest.Body = r.safeBody
	} else {
		switch r.Operation.HTTPMethod {
		case "GET", "HEAD", "DELETE":
			r.HTTPRequest.Body = noBodyReader
		default:
			r.HTTPRequest.Body = r.safeBody
		}
	}
}

func computeBodyLength(r io.ReadSeeker) (int64, error) {
	seekable := true
	switch v := r.(type) {
	case ReaderSeekerCloser:
		seekable = v.IsSeeker()
	case *ReaderSeekerCloser:
		seekable = v.IsSeeker()
	}
	if !seekable {
		return -1, nil
	}

	curOffset, err := r.Seek(0, 1)
	if err != nil {
		return 0, err
	}

	endOffset, err := r.Seek(0, 2)
	if err != nil {
		return 0, err
	}

	_, err = r.Seek(curOffset, 0)
	if err != nil {
		return 0, err
	}

	return endOffset - curOffset, nil
}

// GetBody will return a io.ReadSeeker of the request's underlying
// input Body
func (r *Request) GetBody() io.ReadSeeker {
	return r.safeBody
}

func (r *Request) copy() *Request {
	req := &Request{}
	*req = *r
	req.Handlers = r.Handlers.Copy()
	op := *r.Operation
	req.Operation = &op
	return req
}

// AddToUserAgent adds the string to the end of the request's current user agent.
func AddToUserAgent(r *Request, s string) {
	curUA := r.HTTPRequest.Header.Get("User-Agent")
	if len(curUA) > 0 {
		s = curUA + " " + s
	}
	r.HTTPRequest.Header.Set("User-Agent", s)
}
