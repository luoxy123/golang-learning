package credentials

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/go-ini/ini"
)

// consts
const (
	SharedCredsProviderName = "SharedCredentialsProvider"
)

// vars
var (
	ErrSharedCredentialsHomeNotFound = errors.New("user home directory not found")
)

// SharedCredentialsProvider retrieves credentials from the currnet user's home directory
type SharedCredentialsProvider struct {
	Filename  string
	Profile   string
	retrieved bool
}

// NewSharedCredentials returns a pointer to a new Credentials object
// wrapping the Profile file provider.
func NewSharedCredentials(filename, profile string) *Credentials {
	return NewCredentials(&SharedCredentialsProvider{
		Filename: filename,
		Profile:  profile,
	})
}

// Retrieve reads and extracts the shared credentials from the current
// users home directory.
func (p *SharedCredentialsProvider) Retrieve() (Value, error) {
	p.retrieved = false

	filename, err := p.filename()
	if err != nil {
		return Value{ProviderName: SharedCredsProviderName}, err
	}

	creds, err := loadProfile(filename, p.profile())
	if err != nil {
		return Value{ProviderName: SharedCredsProviderName}, err
	}

	p.retrieved = true
	return creds, nil
}

// IsExpired returns if the shared credentials have expired.
func (p *SharedCredentialsProvider) IsExpired() bool {
	return !p.retrieved
}

func loadProfile(filename, profile string) (Value, error) {
	config, err := ini.Load(filename)
	if err != nil {
		return Value{ProviderName: SharedCredsProviderName}, err
	}
	iniProfile, err := config.GetSection(profile)
	if err != nil {
		return Value{ProviderName: SharedCredsProviderName}, err
	}

	key, err := iniProfile.GetKey("access_key")
	if err != nil {
		return Value{ProviderName: SharedCredsProviderName}, err
	}

	secret, err := iniProfile.GetKey("secret")
	if err != nil {
		return Value{ProviderName: SharedCredsProviderName}, err
	}

	return Value{
		AccessKey:    key.String(),
		Secret:       secret.String(),
		ProviderName: SharedCredsProviderName,
	}, nil
}

func (p *SharedCredentialsProvider) filename() (string, error) {
	if p.Filename == "" {
		homeDir := os.Getenv("HOME")
		if homeDir == "" {
			homeDir = os.Getenv("USERPROFILE")
		}

		if homeDir == "" {
			return "", ErrSharedCredentialsHomeNotFound
		}

		p.Filename = filepath.Join(homeDir, ".polaris", "credentials")
	}

	return p.Filename, nil
}

func (p *SharedCredentialsProvider) profile() string {
	if p.Profile == "" {
		p.Profile = "default"
	}

	return p.Profile
}
