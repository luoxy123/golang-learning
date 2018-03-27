package polaris

import (
	"net/http"

	"code.feelbus.cn/Polaris/polaris-sdk-go/credentials"
)

// ConfigProvider provides a generic way for a service clients
type ConfigProvider interface {
	ClientConfig(serviceName string, cfgs ...*Config) *Config
}

// Config is used to configure the creation of a client
type Config struct {
	Address     *string
	DisableSSL  *bool
	HTTPClient  *http.Client
	serviceName string
	Credentials *credentials.Credentials
}

// WithDisableSSL sets
func (c *Config) WithDisableSSL(disable bool) *Config {
	c.DisableSSL = &disable
	return c
}

// WithHTTPClient sets
func (c *Config) WithHTTPClient(client *http.Client) *Config {
	c.HTTPClient = client
	return c
}

// WithCredentials is
func (c *Config) WithCredentials(creds *credentials.Credentials) *Config {
	c.Credentials = creds
	return c
}

// MergeIn is
func (c *Config) MergeIn(cfgs ...*Config) {
	for _, other := range cfgs {
		mergeConfig(c, other)
	}
}

// Copy returns
func (c *Config) Copy(cfgs ...*Config) *Config {
	dst := &Config{}
	dst.MergeIn(c)

	for _, cfg := range cfgs {
		dst.MergeIn(cfg)
	}

	return dst
}

func mergeConfig(dst *Config, other *Config) {
	if other == nil {
		return
	}

	if other.Address != nil {
		dst.Address = other.Address
	}
	if other.DisableSSL != nil {
		dst.DisableSSL = other.DisableSSL
	}
	if other.HTTPClient != nil {
		dst.HTTPClient = other.HTTPClient
	}
	if other.serviceName != "" {
		dst.serviceName = other.serviceName
	}
	if other.Credentials != nil {
		dst.Credentials = other.Credentials
	}
}

// NewConfig returns a new Config pointer
func NewConfig() *Config {
	return &Config{}
}
