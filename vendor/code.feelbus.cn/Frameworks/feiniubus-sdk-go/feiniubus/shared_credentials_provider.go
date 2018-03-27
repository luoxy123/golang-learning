package feiniubus

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-ini/ini"
)

// SharedCredsProviderName provides a name of SharedCreds provider
const SharedCredsProviderName = "SharedCredentialsProvider"

// SharedCredentialsProvider retrieves credentials from the current user's home
// directory.
type SharedCredentialsProvider struct {
	Filename  string
	Profile   string
	retrieved bool
}

// NewSharedCredentials returns a pointer to a new Credentials object
func NewSharedCredentials(filename, profile string) *Credentials {
	return NewCredentials(&SharedCredentialsProvider{
		Filename: filename,
		Profile:  profile,
	})
}

// Retrieve reads and extracts the shared credentials from the current
// user's home directory
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

// IsExpired returns if the shared credentials have expired
func (p *SharedCredentialsProvider) IsExpired() bool {
	return !p.retrieved
}

func loadProfile(filename, profile string) (Value, error) {
	config, err := ini.Load(filename)
	if err != nil {
		return Value{ProviderName: SharedCredsProviderName}, errors.New("failed to load shared credentials file")
	}
	iniProfile, err := config.GetSection(profile)
	if err != nil {
		return Value{ProviderName: SharedCredsProviderName}, errors.New("failed to get profile")
	}

	id, err := iniProfile.GetKey("access_key_id")
	if err != nil {
		return Value{ProviderName: SharedCredsProviderName}, fmt.Errorf("shared credentials %s in %s did not contain access_key_id", profile, filename)
	}

	secret, err := iniProfile.GetKey("secret_access_key")
	if err != nil {
		return Value{ProviderName: SharedCredsProviderName}, fmt.Errorf("shared credentials %s in %s did not contain secret_access_key", profile, filename)
	}

	return Value{
		AccessKeyID:     id.String(),
		SecretAccessKey: secret.String(),
		ProviderName:    SharedCredsProviderName,
	}, nil
}

func (p *SharedCredentialsProvider) filename() (string, error) {
	if p.Filename == "" {
		if p.Filename = os.Getenv("FNBUS_SHARED_CREDENTIALS_FILE"); p.Filename != "" {
			return p.Filename, nil
		}

		homeDir := os.Getenv("HOME")
		if homeDir == "" {
			homeDir = os.Getenv("USERPROFILE")
		}
		if homeDir == "" {
			return "", errors.New("user home directory not found")
		}

		p.Filename = filepath.Join(homeDir, ".feiniubus", "credentials")
	}

	return p.Filename, nil
}

func (p *SharedCredentialsProvider) profile() string {
	if p.Profile == "" {
		p.Profile = os.Getenv("FNBUS_PROFILE")
	}
	if p.Profile == "" {
		p.Profile = "default"
	}

	return p.Profile
}
