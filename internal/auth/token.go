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

// CachedTokenProvider wraps any TokenProvider with TTL-based caching.
// Uses RWMutex for non-blocking cache reads and a separate refreshMu to
// serialize refresh operations so only one az CLI call happens at a time
// while other goroutines wait briefly then get the cached result.
type CachedTokenProvider struct {
	inner          TokenProvider
	token          string
	expiry         time.Time
	mu             sync.RWMutex // protects token/expiry reads and writes
	refreshMu      sync.Mutex   // serializes refresh operations
	bufferDuration time.Duration
}

// NewCachedTokenProvider creates a cached wrapper around any TokenProvider.
// bufferDuration is how long before actual expiry to refresh (e.g., 5 minutes).
func NewCachedTokenProvider(inner TokenProvider, bufferDuration time.Duration) *CachedTokenProvider {
	return &CachedTokenProvider{
		inner:          inner,
		bufferDuration: bufferDuration,
	}
}

// isValid reports whether the cached token is still usable (not expired minus buffer).
func (c *CachedTokenProvider) isValid() bool {
	return c.token != "" && time.Now().Before(c.expiry.Add(-c.bufferDuration))
}

// GetToken returns a cached token if still valid, otherwise fetches a new one.
// Cache hits are non-blocking (RLock). Refresh is serialized so only one
// goroutine calls the inner provider; others wait and then see the fresh cache.
func (c *CachedTokenProvider) GetToken() (string, error) {
	// Fast path: concurrent read-lock check
	c.mu.RLock()
	if c.isValid() {
		token := c.token
		c.mu.RUnlock()
		return token, nil
	}
	c.mu.RUnlock()

	// Slow path: serialize refresh attempts
	c.refreshMu.Lock()
	defer c.refreshMu.Unlock()

	// Double-check: another goroutine may have refreshed while we waited
	c.mu.RLock()
	if c.isValid() {
		token := c.token
		c.mu.RUnlock()
		return token, nil
	}
	c.mu.RUnlock()

	// Actually refresh (no locks held — won't block other readers)
	token, err := c.inner.GetToken()
	if err != nil {
		return "", err
	}

	// Update cache under write lock
	c.mu.Lock()
	c.token = token
	if ep, ok := c.inner.(interface{ Expiry() time.Time }); ok {
		c.expiry = ep.Expiry()
	} else {
		c.expiry = time.Now().Add(50 * time.Minute)
	}
	c.mu.Unlock()

	return token, nil
}

// Name delegates to the inner provider.
func (c *CachedTokenProvider) Name() string {
	return c.inner.Name()
}
