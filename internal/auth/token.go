package auth

import (
	"sync"
	"time"
)

// TokenProvider abstracts ADO token acquisition.
// Implementations: AzCliTokenProvider, PATTokenProvider.
type TokenProvider interface {
	// GetToken returns a valid ADO access token, refreshing if needed.
	GetToken() (string, error)
	// Name returns the provider name for logging (e.g., "az-cli", "pat").
	Name() string
}

// CachedTokenProvider wraps any TokenProvider with TTL-based caching and mutex.
type CachedTokenProvider struct {
	inner          TokenProvider
	token          string
	expiry         time.Time
	mu             sync.Mutex
	bufferDuration time.Duration // re-fetch buffer before actual expiry
}

// NewCachedTokenProvider creates a cached wrapper around any TokenProvider.
// bufferDuration is how long before actual expiry to refresh (e.g., 5 minutes).
func NewCachedTokenProvider(inner TokenProvider, bufferDuration time.Duration) *CachedTokenProvider {
	return &CachedTokenProvider{
		inner:          inner,
		bufferDuration: bufferDuration,
	}
}

// GetToken returns a cached token if still valid, otherwise fetches a new one.
func (c *CachedTokenProvider) GetToken() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.token != "" && time.Now().Before(c.expiry.Add(-c.bufferDuration)) {
		return c.token, nil
	}

	token, err := c.inner.GetToken()
	if err != nil {
		return "", err
	}

	c.token = token

	// Check if inner provider exposes expiry (e.g., AzCliTokenProvider)
	if ep, ok := c.inner.(interface{ Expiry() time.Time }); ok {
		c.expiry = ep.Expiry()
	} else {
		// Default: assume ~60 minute token lifetime (az cli default)
		c.expiry = time.Now().Add(50 * time.Minute)
	}

	return token, nil
}

// Name delegates to the inner provider.
func (c *CachedTokenProvider) Name() string {
	return c.inner.Name()
}
