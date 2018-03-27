package feiniubus

import (
	"net/http"

	"github.com/FeiniuBus/log"
)

// Config provides service configuration for service client
type Config struct {
	Endpoint               *string
	DisableSSL             *bool
	HTTPClient             *http.Client
	Logger                 *log.Logger
	DisableParamValidation *bool
}

// NewConfig returns a new Config pointer
func NewConfig() *Config {
	return &Config{}
}

// WithEndpoint sets a config Endpoint
func (c *Config) WithEndpoint(endpoint string) *Config {
	c.Endpoint = &endpoint
	return c
}

// WithDisableSSL sets a config DisableSSL
func (c *Config) WithDisableSSL(disable bool) *Config {
	c.DisableSSL = &disable
	return c
}

// WithHTTPClient sets a config HTTPClient
func (c *Config) WithHTTPClient(client *http.Client) *Config {
	c.HTTPClient = client
	return c
}

// WithLogger sets a config logger
func (c *Config) WithLogger(logger *log.Logger) *Config {
	c.Logger = logger
	return c
}

// WithDisableParamValidation sets a config DisableParamValidation
func (c *Config) WithDisableParamValidation(disable bool) *Config {
	c.DisableParamValidation = &disable
	return c
}

// MergeIn merges the passed in configs into the existing config object.
func (c *Config) MergeIn(cfgs ...*Config) {
	for _, other := range cfgs {
		mergeInConfig(c, other)
	}
}

func mergeInConfig(dst *Config, other *Config) {
	if other == nil {
		return
	}

	if other.Endpoint != nil {
		dst.Endpoint = other.Endpoint
	}

	if other.DisableSSL != nil {
		dst.DisableSSL = other.DisableSSL
	}

	if other.HTTPClient != nil {
		dst.HTTPClient = other.HTTPClient
	}

	if other.Logger != nil {
		dst.Logger = other.Logger
	}

	if other.DisableParamValidation != nil {
		dst.DisableParamValidation = other.DisableParamValidation
	}
}

// Copy will return a shallow copy of the Config object
func (c *Config) Copy(cfgs ...*Config) *Config {
	dst := &Config{}
	dst.MergeIn(c)

	for _, cfg := range cfgs {
		dst.MergeIn(cfg)
	}

	return dst
}
