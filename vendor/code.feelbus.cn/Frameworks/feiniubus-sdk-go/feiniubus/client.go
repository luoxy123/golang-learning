package feiniubus

// ClientInfo wraps immutable data from the client.Client structure.
type ClientInfo struct {
	ServiceName string
	Endpoint    string
}

// ClientConfig provides configuration to a service client instance
type ClientConfig struct {
	Config   *Config
	Handlers Handlers
	Endpoint string
}

// ClientConfigProvider provides a generic way for a service client
type ClientConfigProvider interface {
	ClientConfig(serviceName string, cfgs ...*Config) ClientConfig
}

// Client implements the base client request and response handling
type Client struct {
	Config Config
	ClientInfo
	Handlers Handlers
}

// NewClient will return a pointer to a new initialized service client
func NewClient(cfg Config, info ClientInfo, handlers Handlers, options ...func(*Client)) *Client {
	svc := &Client{
		Config:     cfg,
		ClientInfo: info,
		Handlers:   handlers,
	}

	for _, option := range options {
		option(svc)
	}

	return svc
}

// NewRequest returns a new request pointer
func (c *Client) NewRequest(operation *Operation) *Request {
	return New(c.Config, c.ClientInfo, c.Handlers, operation)
}
