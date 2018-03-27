package feiniubus

import (
	"net/url"
	"path"
)

// Session satisfies the service client's client.ClientConfigProvider
type Session struct {
	Config   *Config
	Handlers Handlers
}

const (
	// DefaultEndpoint is a default endpoint
	DefaultEndpoint = "dev.feelbus.cn"
)

// NewSession returns a new Session created from SDK defaults
func NewSession(cfgs ...*Config) (*Session, error) {
	return newSession(cfgs...)
}

func newSession(cfgs ...*Config) (*Session, error) {
	cfg := DefaultConfig()
	handlers := DefaultHandlers()

	cfg.MergeIn(cfgs...)

	s := &Session{
		Config:   cfg,
		Handlers: handlers,
	}

	return s, nil
}

// ClientConfig satisfies the client.ClientConfigProvider interface
func (s *Session) ClientConfig(serviceName string, cfgs ...*Config) ClientConfig {
	cfg, _ := s.clientConfigWithErr(serviceName, cfgs...)
	return cfg
}

func (s *Session) clientConfigWithErr(serviceName string, cfgs ...*Config) (ClientConfig, error) {
	s = s.Copy(cfgs...)

	var resolved ResolvedEndpoint
	var err error

	var endpoint string
	if endpoint = StringValue(s.Config.Endpoint); len(endpoint) == 0 {
		endpoint = DefaultEndpoint
	}

	host := AddScheme(endpoint, BoolValue(s.Config.DisableSSL))
	u, _ := url.Parse(host)
	u.Path = path.Join(u.Path, serviceName)
	resolved.URL = u.String()

	return ClientConfig{
		Config:   s.Config,
		Handlers: s.Handlers,
		Endpoint: resolved.URL,
	}, err
}

// Copy creates and returns a copy of the current Session
func (s *Session) Copy(cfgs ...*Config) *Session {
	newSession := &Session{
		Config:   s.Config.Copy(cfgs...),
		Handlers: s.Handlers.Copy(),
	}

	return newSession
}
