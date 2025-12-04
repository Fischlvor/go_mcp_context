package cache

import (
	"errors"
	"time"
)

// ErrCacheMiss is returned when the key is not found in cache
var ErrCacheMiss = errors.New("cache miss")

// Cache defines the interface for caching
type Cache interface {
	// Get retrieves a value from cache
	Get(key string, dest interface{}) error

	// Set stores a value in cache with TTL
	Set(key string, value interface{}, ttl time.Duration) error

	// Delete removes a key from cache
	Delete(key string) error

	// Exists checks if a key exists in cache
	Exists(key string) (bool, error)

	// Clear removes all keys with the given prefix
	Clear(prefix string) error

	// Close closes the cache connection
	Close() error
}
