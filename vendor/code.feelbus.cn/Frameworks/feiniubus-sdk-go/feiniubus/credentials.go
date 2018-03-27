package feiniubus

import (
	"sync"
)

// Value is a credentials value
type Value struct {
	AccessKeyID     string
	SecretAccessKey string
	ProviderName    string
}

// Provider is the interface for any component which will provider credentials
// Value
type Provider interface {
	Retrieve() (Value, error)
	IsExpired() bool
}

// Credentials providers synchronous safe retrieval of credentials.
type Credentials struct {
	creds        Value
	m            sync.Mutex
	provider     Provider
	forceRefresh bool
}

// NewCredentials returns a pointer to a new credentials
func NewCredentials(provider Provider) *Credentials {
	return &Credentials{
		provider:     provider,
		forceRefresh: true,
	}
}

// Get returns the credentials value
func (c *Credentials) Get() (Value, error) {
	c.m.Lock()
	defer c.m.Unlock()

	if c.isExpired() {
		creds, err := c.provider.Retrieve()
		if err != nil {
			return Value{}, err
		}
		c.creds = creds
		c.forceRefresh = false
	}

	return c.creds, nil
}

// Expire expires the credentials and forces them to be retrieved on the next
// call o Get
func (c *Credentials) Expire() {
	c.m.Lock()
	defer c.m.Unlock()

	c.forceRefresh = true
}

// IsExpired returns if the credentials are no longer valid
func (c *Credentials) IsExpired() bool {
	c.m.Lock()
	defer c.m.Unlock()

	return c.isExpired()
}

func (c *Credentials) isExpired() bool {
	return c.forceRefresh || c.provider.IsExpired()
}
