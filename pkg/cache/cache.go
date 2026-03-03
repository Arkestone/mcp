// Package cache provides a generic in-memory list cache with TTL-based
// invalidation, designed to protect loaders from repeated expensive disk scans
// within a single agent turn.
//
// Usage:
//
//	var c cache.List[MyItem]
//
//	// In List():
//	return c.Get(func() []MyItem { return expensiveScan() })
//
//	// In ForceSync():
//	c.Invalidate()
package cache

import (
	"sync"
	"time"
)

// DefaultTTL is the default duration cached results remain valid.
// Chosen to be short enough for interactive responsiveness (a single agent
// turn completes well under 5 s) while eliminating redundant disk walks
// when multiple tools are called in quick succession.
const DefaultTTL = 5 * time.Second

// List is a thread-safe in-memory cache for a slice of items.
// The zero value is ready to use with DefaultTTL.
type List[T any] struct {
	mu      sync.RWMutex
	items   []T
	fetched time.Time
	valid   bool
	TTL     time.Duration // 0 means DefaultTTL
}

// Get returns cached items if still fresh, otherwise calls fill to populate
// the cache. Thread-safe; concurrent callers collapse to a single fill call.
func (c *List[T]) Get(fill func() []T) []T {
	ttl := c.TTL
	if ttl == 0 {
		ttl = DefaultTTL
	}

	// Fast path: valid cache under read lock.
	c.mu.RLock()
	if c.valid && time.Since(c.fetched) < ttl {
		items := c.items
		c.mu.RUnlock()
		return items
	}
	c.mu.RUnlock()

	// Slow path: acquire write lock and re-check (double-checked locking).
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.valid && time.Since(c.fetched) < ttl {
		return c.items
	}
	c.items = fill()
	c.fetched = time.Now()
	c.valid = true
	return c.items
}

// Invalidate marks the cache as stale. The next Get call will re-populate.
func (c *List[T]) Invalidate() {
	c.mu.Lock()
	c.valid = false
	c.mu.Unlock()
}
