package auth

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// invalidationCache tracks JWTs that have been explicitly logged out, so
// they're rejected even though they haven't expired yet.
type invalidationCache struct {
	c *cache.Cache
}

func newInvalidationCache() *invalidationCache {
	return &invalidationCache{c: cache.New(24*time.Hour, time.Hour)}
}

func (i *invalidationCache) invalidate(token string, ttl time.Duration) {
	i.c.Set(token, true, ttl)
}

func (i *invalidationCache) isInvalidated(token string) bool {
	_, found := i.c.Get(token)
	return found
}
