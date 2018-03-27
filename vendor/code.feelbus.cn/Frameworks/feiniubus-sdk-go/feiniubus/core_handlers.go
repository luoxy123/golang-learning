package feiniubus

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type lener interface {
	Len() int
}

// RFC822 returns an RFC822 formatted timestamp
const RFC822 = "Mon, 2 Jan 2006 15:04:05 GMT"

var errValueNotSet = fmt.Errorf("value not set")

// ValidateEndpointHandler is a reqeust handler to validate a request had the
// appropriate Endpoint set.
var ValidateEndpointHandler = NamedHandler{
	Name: "core.ValidateEndpointHandler",
	Fn: func(r *Request) {
		if r.ClientInfo.Endpoint == "" {
			r.Error = ErrMissingEndpoint
		}
	},
}

// SDKVersionUserAgentHandler is a request handler for adding the SDK Version to the user agent.
var SDKVersionUserAgentHandler = NamedHandler{
	Name: "core.SDKVersionUserAgentHandler",
	Fn:   MakeAddToUserAgentHandler(SDKName, SDKVersion, runtime.Version(), runtime.GOOS, runtime.GOARCH),
}

// BuildContentLengthHandler builds the content length of a request
var BuildContentLengthHandler = NamedHandler{
	Name: "core.BuildContentLengthHandler",
	Fn: func(r *Request) {
		var length int64

		if slength := r.HTTPRequest.Header.Get("Content-Length"); slength != "" {
			length, _ = strconv.ParseInt(slength, 10, 64)
		} else {
			switch body := r.Body.(type) {
			case nil:
				length = 0
			case lener:
				length = int64(body.Len())
			case io.Seeker:
				r.BodyStart, _ = body.Seek(0, 1)
				end, _ := body.Seek(0, 2)
				length = end - r.BodyStart
			default:
				panic("Cannot get length of body, must provide `ContentLength`")
			}
		}

		if length > 0 {
			r.HTTPRequest.ContentLength = length
			r.HTTPRequest.Header.Set("Content-Length", fmt.Sprintf("&d", length))
		} else {
			r.HTTPRequest.ContentLength = 0
			r.HTTPRequest.Header.Del("Content-Length")
		}
	},
}

var reStatusCode = regexp.MustCompile(`^(\d{3})`)

// SendHandler is a request handler to send service request using HTTP client
var SendHandler = NamedHandler{
	Name: "core.SendHandler",
	Fn: func(r *Request) {
		var err error
		r.HTTPResponse, err = r.Config.HTTPClient.Do(r.HTTPRequest)
		if err != nil {
			if r.HTTPResponse != nil {
				r.HTTPResponse.Body.Close()
			}
			if e, ok := err.(*url.Error); ok && e.Err != nil {
				if s := reStatusCode.FindStringSubmatch(e.Err.Error()); s != nil {
					code, _ := strconv.ParseInt(s[1], 10, 64)
					r.HTTPResponse = &http.Response{
						StatusCode: int(code),
						Status:     http.StatusText(int(code)),
						Body:       ioutil.NopCloser(bytes.NewReader([]byte{})),
					}
					return
				}
			}
			if r.HTTPResponse == nil {
				r.HTTPResponse = &http.Response{
					StatusCode: int(0),
					Status:     http.StatusText(int(0)),
					Body:       ioutil.NopCloser(bytes.NewReader([]byte{})),
				}
			}
			r.Error = err
		}
	},
}

// ValidateResponseHandler is a request handler to validate service response.
var ValidateResponseHandler = NamedHandler{
	Name: "core.ValidateResponseHandler",
	Fn: func(r *Request) {
		if r.HTTPResponse.StatusCode == 0 || r.HTTPResponse.StatusCode >= 300 {
			r.Error = errors.New("unknown error")
		}
	},
}

// BuildHandler is a named request handler for building request params
var BuildHandler = NamedHandler{
	Name: "feiniubussdk.rest.Build",
	Fn: func(r *Request) {
		r.HTTPRequest.Method = r.Operation.HTTPMethod

		if r.UseQueryString() {
			if r.ParamsFilled() {
				query := r.HTTPRequest.URL.Query()
				for k, v := range r.Params {
					query.Set(k, v)
				}
				r.HTTPRequest.URL.RawQuery = query.Encode()
			}
		}

		if r.ContentFilled() {
			var err error
			if strings.ToUpper(r.Operation.ContentType) == "FORM" {
				k := reflect.TypeOf(r.Content).Kind()
				if k == reflect.Array || k == reflect.Slice {
					r.Error = fmt.Errorf("The request content can't be array or slice when Content-Type is %s", strconv.Quote("FORM"))
					return
				}

				v := reflect.ValueOf(r.Content).Elem()
				err = buildFormBody(r, v)
				if err != nil {
					r.Error = err
					return
				}
				r.HTTPRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
			} else {
				err = buildJSONBody(r, r.Content)
				if err != nil {
					r.Error = err
					return
				}
				r.HTTPRequest.Header.Set("Content-Type", "application/json")
			}
		}

		r.HTTPRequest.Header.Set("Accept-Encoding", "gzip")
		r.HTTPRequest.Header.Set("Accept", "application/json")
		r.HTTPRequest.Header.Set("Connection", "keep-alive")
	},
}

func buildJSONBody(r *Request, v interface{}) error {
	c, err := json.Marshal(v)
	if err != nil {
		return err
	}
	r.SetBufferBody(c)
	return nil
}

func buildFormBody(r *Request, v reflect.Value) error {
	u := &url.URL{}
	query := u.Query()

	for i := 0; i < v.NumField(); i++ {
		m := v.Field(i)
		if n := v.Type().Field(i).Name; n[0:1] == strings.ToLower(n[0:1]) {
			continue
		}

		if m.IsValid() {
			field := v.Type().Field(i)
			name := field.Tag.Get("locationName")
			if name == "" {
				name = field.Name
			}
			if m.Kind() == reflect.Ptr {
				m = m.Elem()
			}
			if !m.IsValid() {
				continue
			}
			if field.Tag.Get("ignore") != "" {
				continue
			}

			str, err := convertType(v)
			if err == errValueNotSet {
				continue
			} else if err != nil {
				return err
			}

			query.Set(name, str)
		}
	}

	c := query.Encode()
	r.SetStringBody(c)
	return nil
}

func convertType(v reflect.Value) (string, error) {
	v = reflect.Indirect(v)
	if !v.IsValid() {
		return "", errValueNotSet
	}

	var str string
	switch value := v.Interface().(type) {
	case string:
		str = value
	case []byte:
		str = base64.StdEncoding.EncodeToString(value)
	case bool:
		str = strconv.FormatBool(value)
	case int64:
		str = strconv.FormatInt(value, 10)
	case float64:
		str = strconv.FormatFloat(value, 'f', -1, 64)
	case time.Time:
		str = value.UTC().Format(RFC822)
	default:
		err := fmt.Errorf("Unsupported value for param %v (%s)", v.Interface(), v.Type())
		return "", err
	}

	return str, nil
}

// ErrorResponse is
type ErrorResponse struct {
	Code    string `json:"error"`
	Message string `json:"error_description"`
}

// UnmarshalErrorHandler is a named request handler for error unmarshal
var UnmarshalErrorHandler = NamedHandler{
	Name: "feiniubus.rest.ErrorUnmarshal",
	Fn: func(r *Request) {
		defer r.HTTPResponse.Body.Close()

		bodyBytes, err := ioutil.ReadAll(r.HTTPResponse.Body)
		if err != nil {
			r.Error = NewError("SerializationError", "failed to read from HTTP response body", err)
			return
		}

		if r.HTTPResponse.StatusCode != 400 {
			r.Error = NewError(r.HTTPResponse.Status, fmt.Sprintf("%d: %s(%s)", r.HTTPResponse.StatusCode, r.HTTPResponse.Status, string(bodyBytes)), nil)
			return
		}

		var respErr ErrorResponse
		decodeErr := json.Unmarshal(bodyBytes, &respErr)
		if decodeErr == nil {
			r.Error = NewError(respErr.Code, respErr.Message, nil)
			return
		}

		var respErrs []ErrorResponse
		decodeErr = json.Unmarshal(bodyBytes, &respErrs)
		if decodeErr == nil {
			if len(respErrs) > 0 {
				r.Error = NewError(respErrs[0].Code, respErrs[0].Message, nil)
				return
			}

			r.Error = NewError(r.HTTPResponse.Status, fmt.Sprintf("%d: %s(%s)", r.HTTPResponse.StatusCode, r.HTTPResponse.Status, string(bodyBytes)), nil)
			return
		}

		r.Error = NewError("SerializationError", "failed to decode json error response", decodeErr)
	},
}

// UnmarshalHandler is a named request handler for unmarshaling response body
var UnmarshalHandler = NamedHandler{
	Name: "feiniubus.rest.Unmarshal",
	Fn: func(r *Request) {
		defer r.HTTPResponse.Body.Close()
		if r.DataFilled() && r.UnmarshalerFilled() {
			// TODO(xqlun): gzip decompression
			encoding := r.HTTPResponse.Header.Get("Content-Encoding")
			if encoding != "" {
				if encoding != "gzip" {
					r.Error = fmt.Errorf("unsupported content-encoding type: %s", encoding)
					return
				}

				reader, err := gzip.NewReader(r.HTTPResponse.Body)
				defer reader.Close()
				if err != nil {
					r.Error = err
					return
				}

				err = r.Unmarshaler.Unmarshal(reader, r.Data)
				if err != nil {
					r.Error = err
				}
			} else {
				err := r.Unmarshaler.Unmarshal(r.HTTPResponse.Body, r.Data)
				if err != nil {
					r.Error = errors.New("failed decoding http response")
				}
			}
		}
	},
}
