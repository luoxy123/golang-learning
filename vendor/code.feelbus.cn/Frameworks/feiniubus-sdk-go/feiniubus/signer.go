package feiniubus

// import (
// 	"crypto/hmac"
// 	"crypto/sha256"
// 	"encoding/hex"
// 	"io"
// 	"net/http"
// 	"net/url"
// 	"strings"
// 	"time"

// 	"bytes"
// 	"fmt"
// 	"sort"
// )

// const (
// 	timeFormat      = "20060102T150405Z"
// 	shortTimeFormat = "20060102"

// 	// emptyStringSHA256 is a SHA256 of an empty string
// 	emptyStringSHA256 = `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`
// 	scheme            = "FNSIGN"
// 	terminator        = "feiniubus_request"
// )

// var noEscape [256]bool

// func init() {
// 	for i := 0; i < len(noEscape); i++ {
// 		noEscape[i] = (i >= 'A' && i <= 'Z') ||
// 			(i >= 'a' && i <= 'z') ||
// 			(i >= '0' && i <= '9') ||
// 			i == '-' ||
// 			i == '.' ||
// 			i == '_' ||
// 			i == '~'
// 	}
// }

// // Signer applies signing to given request.
// type Signer struct {
// 	Credentials            *Credentials
// 	DisableURIPathEscaping bool
// 	currentTimeFn          func() time.Time
// }

// // NewSigner returns a Signer pointer
// func NewSigner(credentials *Credentials, options ...func(*Signer)) *Signer {
// 	signer := &Signer{
// 		Credentials: credentials,
// 	}

// 	for _, option := range options {
// 		option(signer)
// 	}

// 	return signer
// }

// type signingCtx struct {
// 	Request *http.Request
// 	Body    io.ReadSeeker
// 	Query   url.Values
// 	Time    time.Time

// 	DisableURIPathEscaping bool

// 	credValues    Value
// 	formattedTime string

// 	bodyDigest       string
// 	canonicalString  string
// 	credentialString string
// 	stringToSign     string
// 	signature        string
// 	authorization    string
// }

// func (sign Signer) signWithBody(r *http.Request, body io.ReadSeeker, signTime time.Time) error {
// 	currentTimeFn := sign.currentTimeFn
// 	if currentTimeFn == nil {
// 		currentTimeFn = time.Now
// 	}

// 	ctx := &signingCtx{
// 		Request: r,
// 		Body:    body,
// 		Query:   r.URL.Query(),
// 		Time:    signTime,
// 		DisableURIPathEscaping: sign.DisableURIPathEscaping,
// 	}

// 	for key := range ctx.Query {
// 		sort.Strings(ctx.Query[key])
// 	}

// 	if ctx.isRequestSigned() {
// 		ctx.Time = currentTimeFn()
// 	}

// 	var err error
// 	ctx.credValues, err = sign.Credentials.Get()
// 	if err != nil {
// 		return err
// 	}

// 	ctx.build()
// 	return nil
// }

// // SignRequestHandler is a named request handler
// var SignRequestHandler = NamedHandler{
// 	Name: "SignRequestHandler",
// 	Fn:   SignSDKRequest,
// }

// // SignSDKRequest signs a request
// func SignSDKRequest(req *Request) {
// 	signSDKRequestWithCurrTime(req, time.Now)
// }

// func signSDKRequestWithCurrTime(req *Request, curTimeFn func() time.Time) {
// 	sign := NewSigner(req.Config.Credentials, func(sign *Signer) {
// 		sign.currentTimeFn = curTimeFn
// 	})

// 	signingTime := req.Time
// 	err := sign.signWithBody(req.HTTPRequest, req.GetBody(), signingTime)
// 	if err != nil {
// 		req.Error = err
// 		return
// 	}
// }

// func (ctx *signingCtx) build() {
// 	ctx.buildTime()
// 	ctx.buildBodyDigest()
// 	ctx.buildCanonicalString()
// 	ctx.buildStringToSign()
// 	ctx.buildSignature()

// 	parts := []string{
// 		"Credential=" + ctx.credValues.AccessKeyID,
// 		"Signature=" + ctx.signature,
// 	}

// 	value := strings.Join(parts, ",")
// 	ctx.Request.Header.Set("Authorization", fmt.Sprintf("%s %s", scheme, value))
// }

// func (ctx *signingCtx) buildStringToSign() {
// 	var prefix = strings.Join([]string{
// 		"HMAC-SHA256",
// 		ctx.formattedTime,
// 	}, "-")

// 	ctx.stringToSign = strings.Join([]string{
// 		prefix,
// 		ctx.credValues.AccessKeyID,
// 		hex.EncodeToString(makeSha256([]byte(ctx.canonicalString))),
// 	}, "\n")
// }

// func (ctx *signingCtx) buildCanonicalString() {
// 	ctx.Request.URL.RawQuery = strings.Replace(ctx.Query.Encode(), "+", "%20", -1)

// 	uri := getURIPath(ctx.Request.URL)

// 	if !ctx.DisableURIPathEscaping {
// 		uri = escapePath(uri, false)
// 	}

// 	buffer := bytes.NewBufferString("")
// 	buffer.WriteString(fmt.Sprintf("%s\n", ctx.Request.Method))
// 	buffer.WriteString(fmt.Sprintf("%s\n", uri))
// 	buffer.WriteString(fmt.Sprintf("%s\n", ctx.Request.URL.RawQuery))
// 	buffer.WriteString(ctx.bodyDigest)

// 	ctx.canonicalString = buffer.String()
// }

// func (ctx *signingCtx) buildBodyDigest() {
// 	var hash string
// 	if ctx.Body == nil {
// 		hash = emptyStringSHA256
// 	} else {
// 		hash = hex.EncodeToString(makeSha256Reader(ctx.Body))
// 	}

// 	ctx.bodyDigest = hash
// }

// func (ctx *signingCtx) buildSignature() {
// 	secret := ctx.credValues.SecretAccessKey
// 	date := makeHmac([]byte(scheme+secret), []byte(ctx.formattedTime))
// 	credentials := makeHmac(date, []byte(terminator))
// 	signature := makeHmac(credentials, []byte(ctx.stringToSign))
// 	ctx.signature = hex.EncodeToString(signature)
// }

// func (ctx *signingCtx) buildTime() {
// 	ctx.formattedTime = ctx.Time.UTC().Format(timeFormat)

// 	ctx.Request.Header.Add("x-feiniubus-date", ctx.formattedTime)
// }

// func (ctx *signingCtx) isRequestSigned() bool {
// 	if ctx.Request.Header.Get("Authorization") != "" {
// 		return true
// 	}

// 	return false
// }

// func makeHmac(key []byte, data []byte) []byte {
// 	hash := hmac.New(sha256.New, key)
// 	hash.Write(data)
// 	return hash.Sum(nil)
// }

// func makeSha256(data []byte) []byte {
// 	hash := sha256.New()
// 	hash.Write(data)
// 	return hash.Sum(nil)
// }

// func makeSha256Reader(reader io.ReadSeeker) []byte {
// 	hash := sha256.New()
// 	start, _ := reader.Seek(0, 1)
// 	defer reader.Seek(start, 0)

// 	io.Copy(hash, reader)
// 	return hash.Sum(nil)
// }

// func escapePath(path string, encodeSep bool) string {
// 	var buf bytes.Buffer
// 	for i := 0; i < len(path); i++ {
// 		c := path[i]
// 		if noEscape[c] || (c == '/' && !encodeSep) {
// 			buf.WriteByte(c)
// 		} else {
// 			fmt.Fprintf(&buf, "%%%02X", c)
// 		}
// 	}
// 	return buf.String()
// }

// func getURIPath(u *url.URL) string {
// 	var uri string

// 	if len(u.Opaque) > 0 {
// 		uri = "/" + strings.Join(strings.Split(u.Opaque, "/")[3:], "/")
// 	} else {
// 		uri = u.EscapedPath()
// 	}

// 	if len(uri) == 0 {
// 		uri = "/"
// 	} else {
// 		if strings.HasPrefix(uri, "/") {
// 			uri = uri[1:]
// 		}
// 		parts := strings.Split(uri, "/")[1:]
// 		uri = "/" + strings.Join(parts, "/")
// 	}

// 	return uri
// }
