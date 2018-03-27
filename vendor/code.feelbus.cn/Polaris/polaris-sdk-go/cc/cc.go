package cc

import (
	"code.feelbus.cn/Polaris/polaris-sdk-go/polaris"
)

// consts
const (
	ServiceName = "cc"
)

// CC is used to manipulate the configuration center api
type CC struct {
	*polaris.Client
}

// New returns a handle to the CC apis
func New(p polaris.ConfigProvider, cfgs ...*polaris.Config) *CC {
	c := p.ClientConfig(ServiceName, cfgs...)
	svc := &CC{
		Client: polaris.NewClient(c),
	}
	return svc
}
