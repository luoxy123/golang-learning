package feiniubus

import (
	"errors"
)

const (
	// StaticProviderName is
	StaticProviderName = "StaticProvider"
)

// StaticProvider is a set of credentials
type StaticProvider struct {
	Value
}

// NewStaticCredentials returns a pointer to a new credentials object
func NewStaticCredentials(id, secret string) *Credentials {
	return NewCredentials(&StaticProvider{
		Value: Value{
			AccessKeyID:     id,
			SecretAccessKey: secret,
		},
	})
}

// NewStaticCredentialsFromCreds returns
func NewStaticCredentialsFromCreds(creds Value) *Credentials {
	return NewCredentials(&StaticProvider{Value: creds})
}

// Retrieve returns
func (s *StaticProvider) Retrieve() (Value, error) {
	if s.AccessKeyID == "" || s.SecretAccessKey == "" {
		return Value{ProviderName: StaticProviderName}, errors.New("static credentials are empty")
	}

	if len(s.Value.ProviderName) == 0 {
		s.Value.ProviderName = StaticProviderName
	}

	return s.Value, nil
}

// IsExpired returns
func (s *StaticProvider) IsExpired() bool {
	return false
}
