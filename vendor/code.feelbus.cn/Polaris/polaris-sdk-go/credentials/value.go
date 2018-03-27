package credentials

import (
	"sync"
)

// AnonymousCredentials is an empty Credential object
var AnonymousCredentials = NewStaticCredentials("", "")

// Value is
type Value struct {
	AccessKey    string
	Secret       string
	ProviderName string
}

// Provider is
type Provider interface {
	Retrieve() (Value, error)
	IsExpired() bool
}

// Credentials provides
type Credentials struct {
	creds        Value
	forceRefresh bool
	m            sync.Mutex

	provider Provider
}

// NewCredentials returns a pointer to a new NewCredentials with the provider set.
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

// Expire expires the credentials and forces them to be retrieved on the
// next call to Get()
func (c *Credentials) Expire() {
	c.m.Lock()
	defer c.m.Unlock()

	c.forceRefresh = true
}

// IsExpired returns if the credentials are no longer valid.
func (c *Credentials) IsExpired() bool {
	c.m.Lock()
	defer c.m.Unlock()

	return c.isExpired()
}

func (c *Credentials) isExpired() bool {
	return c.forceRefresh || c.provider.IsExpired()
}
