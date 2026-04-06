package auth

import (
	"fmt"
	"log"

	"github.com/zalando/go-keyring"
)

const (
	defaultServiceName = "team-ado-tool"
	defaultKeyName     = "pat"
)

// PATTokenProvider retrieves a Personal Access Token from the OS keyring.
type PATTokenProvider struct {
	serviceName string
	keyName     string
}

// NewPATTokenProvider creates a new PAT-based token provider using the OS keyring.
func NewPATTokenProvider() *PATTokenProvider {
	return &PATTokenProvider{
		serviceName: defaultServiceName,
		keyName:     defaultKeyName,
	}
}

// GetToken retrieves the PAT from the OS keyring.
func (p *PATTokenProvider) GetToken() (string, error) {
	log.Printf("[auth] using %s token provider", p.Name())

	pat, err := keyring.Get(p.serviceName, p.keyName)
	if err != nil {
		return "", fmt.Errorf("no PAT stored — use Settings to configure: %w", err)
	}

	if pat == "" {
		return "", fmt.Errorf("no PAT stored — use Settings to configure")
	}

	return pat, nil
}

// Name returns the provider identifier.
func (p *PATTokenProvider) Name() string {
	return "pat"
}
