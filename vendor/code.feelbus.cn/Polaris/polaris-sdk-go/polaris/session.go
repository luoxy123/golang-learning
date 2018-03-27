package polaris

// Session provides a central location to create service client
type Session struct {
	Config *Config
}

const (
	defaultAddress = "dc.feiniubus.com"
)

// NewSession is
func NewSession(cfgs ...*Config) *Session {
	cfg := DefaultConfig()
	userCfg := &Config{}
	userCfg.MergeIn(cfgs...)

	cfg.MergeIn(userCfg)

	return &Session{
		Config: cfg,
	}
}

// Copy is
func (s *Session) Copy(cfgs ...*Config) *Session {
	newSession := &Session{
		Config: s.Config.Copy(cfgs...),
	}

	return newSession
}

// ClientConfig implement the ConfigProvider interface
func (s *Session) ClientConfig(servcieName string, cfgs ...*Config) *Config {
	s = s.Copy(cfgs...)
	result := &Config{}
	result.MergeIn(s.Config)

	var endpoint string
	if endpoint = StringValue(s.Config.Address); len(endpoint) == 0 {
		endpoint = defaultAddress
	}

	host := addScheme(endpoint, BoolValue(s.Config.DisableSSL))

	result.Address = String(host)
	result.serviceName = servcieName

	return result
}
