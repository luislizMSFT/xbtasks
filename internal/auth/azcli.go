package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
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

// resolveAzPath finds the az CLI binary, expanding PATH for macOS .app bundles
// which launch with a minimal PATH that excludes Homebrew/user paths.
func resolveAzPath() string {
	if p, err := exec.LookPath("az"); err == nil {
		return p
	}
	if runtime.GOOS == "darwin" {
		for _, candidate := range []string{
			"/opt/homebrew/bin/az",
			"/usr/local/bin/az",
		} {
			if _, err := os.Stat(candidate); err == nil {
				return candidate
			}
		}
	}
	return "az"
}

// ensurePATH adds common tool directories so child processes can find
// binaries like az, git, etc. even when launched from a .app bundle.
func ensurePATH() {
	if runtime.GOOS != "darwin" {
		return
	}
	extra := []string{"/opt/homebrew/bin", "/usr/local/bin"}
	cur := os.Getenv("PATH")
	var missing []string
	for _, d := range extra {
		if !strings.Contains(cur, d) {
			missing = append(missing, d)
		}
	}
	if len(missing) > 0 {
		os.Setenv("PATH", cur+":"+strings.Join(missing, ":"))
	}
}

func init() {
	ensurePATH()
}

// GetToken runs `az account get-access-token` with the ADO resource ID and
// returns the access token string.
func (p *AzCliTokenProvider) GetToken() (string, error) {
	log.Printf("[auth] using %s token provider", p.Name())

	azBin := resolveAzPath()
	cmd := exec.Command(azBin, "account", "get-access-token",
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
