package fns

import (
	"code.feelbus.cn/Polaris/polaris-sdk-go/polaris"
)

const (
	ServiceName = "fns"
)

type FNS struct {
	*polaris.Client
}

// New returns a handle to the FNS apis
func New(p polaris.ConfigProvider, cfgs ...*polaris.Config) *FNS {
	c := p.ClientConfig(ServiceName, cfgs...)
	svc := &FNS{
		Client: polaris.NewClient(c),
	}
	return svc
}
