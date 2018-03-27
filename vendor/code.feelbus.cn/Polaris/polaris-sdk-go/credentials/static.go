package credentials

import (
	"errors"
)

// consts
const (
	StaticProviderName = "StaticProvider"
)

// vars
var (
	ErrStaticCredentialsEmpty = errors.New("static credentials are empty")
)

// StaticProvider is a set of credentials which are set programmatically,
// and will never expire.
type StaticProvider struct {
	Value
}

// NewStaticCredentials returns a pointer to a new Credentials object
func NewStaticCredentials(key, secret string) *Credentials {
	return NewCredentials(&StaticProvider{
		Value: Value{
			AccessKey: key,
			Secret:    secret,
		},
	})
}

// Retrieve returns the credentials or error if the credentials are invalid
func (s *StaticProvider) Retrieve() (Value, error) {
	if s.AccessKey == "" || s.Secret == "" {
		return Value{ProviderName: StaticProviderName}, ErrStaticCredentialsEmpty
	}

	if len(s.Value.ProviderName) == 0 {
		s.Value.ProviderName = StaticProviderName
	}

	return s.Value, nil
}

// IsExpired returns if the credentials are expired.
func (s *StaticProvider) IsExpired() bool {
	return false
}
