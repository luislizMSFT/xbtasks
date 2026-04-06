package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"
)

// adoResourceID is the Azure DevOps resource identifier for token acquisition.
const adoResourceID = "499b84ac-1321-427f-aa17-267ca6975798"

// azTokenResponse represents the JSON response from `az account get-access-token`.
type azTokenResponse struct {
	AccessToken string `json:"accessToken"`
	ExpiresOn   string `json:"expiresOn"`
	TokenType   string `json:"tokenType"`
}

// AzCliTokenProvider fetches ADO tokens via the az CLI.
type AzCliTokenProvider struct {
	lastExpiry time.Time
}

// NewAzCliTokenProvider creates a new az CLI-based token provider.
func NewAzCliTokenProvider() *AzCliTokenProvider {
	return &AzCliTokenProvider{}
}

// GetToken runs `az account get-access-token` with the ADO resource ID and
// returns the access token string.
func (p *AzCliTokenProvider) GetToken() (string, error) {
	log.Printf("[auth] using %s token provider", p.Name())

	cmd := exec.Command("az", "account", "get-access-token",
		"--resource", adoResourceID,
		"--output", "json",
	)

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("az CLI failed — run 'az login' first: %s", string(exitErr.Stderr))
		}
		return "", fmt.Errorf("az CLI failed — run 'az login' first: %w", err)
	}

	var resp azTokenResponse
	if err := json.Unmarshal(output, &resp); err != nil {
		return "", fmt.Errorf("parse az token response: %w", err)
	}

	if resp.AccessToken == "" {
		return "", fmt.Errorf("az CLI returned empty access token")
	}

	// Parse expiry — az cli uses format "2006-01-02 15:04:05.000000"
	if resp.ExpiresOn != "" {
		if t, err := time.Parse("2006-01-02 15:04:05.000000", resp.ExpiresOn); err == nil {
			p.lastExpiry = t
		}
	}

	return resp.AccessToken, nil
}

// Name returns the provider identifier.
func (p *AzCliTokenProvider) Name() string {
	return "az-cli"
}

// Expiry returns the expiry time of the last fetched token.
func (p *AzCliTokenProvider) Expiry() time.Time {
	return p.lastExpiry
}
