package visitor

import (
	"code.feelbus.cn/Frameworks/feiniubus-sdk-go/feiniubus"
)

// Visitor is a service client for visitor service
type Visitor struct {
	*feiniubus.Client
}

// prefix API calls
const (
	ServiceName = "visitor"
)

// New creates a new instance of the order client
func New(p feiniubus.ClientConfigProvider, cfgs ...*feiniubus.Config) *Visitor {
	c := p.ClientConfig(ServiceName, cfgs...)
	return newClient(*c.Config, c.Handlers, c.Endpoint)
}

func newClient(cfg feiniubus.Config, handlers feiniubus.Handlers, endpoint string) *Visitor {
	svc := &Visitor{
		Client: feiniubus.NewClient(cfg, feiniubus.ClientInfo{
			ServiceName: ServiceName,
			Endpoint:    endpoint,
		}, handlers),
	}

	return svc
}
