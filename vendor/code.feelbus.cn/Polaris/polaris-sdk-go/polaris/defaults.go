package polaris

import (
	"net/http"

	"code.feelbus.cn/Polaris/polaris-sdk-go/credentials"
)

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return NewConfig().
		WithDisableSSL(true).
		WithHTTPClient(http.DefaultClient).
		WithCredentials(credentials.AnonymousCredentials)
}
