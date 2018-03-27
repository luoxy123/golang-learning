package feiniubus

import (
	"fmt"
	"regexp"
)

// Options provides the configuration needed to direct how the
// endpoints will be resolved
type Options struct {
	DisableSSL bool
}

var schemeRE = regexp.MustCompile("^([^:]+)://")

// AddScheme adds the HTTP or HTTPS schemes to a endpoint URL
func AddScheme(endpoint string, disableSSL bool) string {
	if !schemeRE.MatchString(endpoint) {
		scheme := "https"
		if disableSSL {
			scheme = "http"
		}
		endpoint = fmt.Sprintf("%s://%s", scheme, endpoint)
	}
	return endpoint
}

// ResolvedEndpoint is an endpoint that has been resolved
type ResolvedEndpoint struct {
	URL string
}
