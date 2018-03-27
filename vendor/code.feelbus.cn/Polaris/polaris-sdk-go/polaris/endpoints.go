package polaris

import (
	"fmt"
	"regexp"
)

var schemeRE = regexp.MustCompile("^([^:]+)://")

func addScheme(endpoint string, disableSSL bool) string {
	if !schemeRE.MatchString(endpoint) {
		scheme := "https"
		if disableSSL {
			scheme = "http"
		}
		endpoint = fmt.Sprintf("%s://%s", scheme, endpoint)
	}

	return endpoint
}
