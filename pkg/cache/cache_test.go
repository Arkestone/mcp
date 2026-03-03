package cache_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Arkestone/mcp/pkg/cache"
)

func TestList_CachesResult(t *testing.T) {
	calls := 0
	var c cache.List[int]
	c.TTL = time.Second

	fill := func() []int { calls++; return []int{1, 2, 3} }

	c.Get(fill)
	c.Get(fill)
	c.Get(fill)

	if calls != 1 {
		t.Errorf("fill called %d times, want 1", calls)
	}
}

func TestList_RefreshesAfterTTL(t *testing.T) {
	calls := 0
	var c cache.List[int]
	c.TTL = 10 * time.Millisecond

	fill := func() []int { calls++; return []int{1} }

	c.Get(fill)
	time.Sleep(20 * time.Millisecond)
	c.Get(fill)

	if calls != 2 {
		t.Errorf("fill called %d times, want 2", calls)
	}
}

func TestList_InvalidateForcesRefill(t *testing.T) {
	calls := 0
	var c cache.List[int]
	c.TTL = time.Hour // won't expire naturally

	fill := func() []int { calls++; return []int{1} }

	c.Get(fill)
	c.Invalidate()
	c.Get(fill)

	if calls != 2 {
		t.Errorf("fill called %d times, want 2", calls)
	}
}

func TestList_ThreadSafe(t *testing.T) {
	var calls atomic.Int64
	var c cache.List[int]
	c.TTL = time.Second

	fill := func() []int {
		calls.Add(1)
		time.Sleep(time.Millisecond) // simulate work
		return []int{1, 2, 3}
	}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Get(fill)
		}()
	}
	wg.Wait()

	// May be called more than once due to concurrent first access, but
	// must be far fewer than 50 (the whole point of caching).
	if calls.Load() > 10 {
		t.Errorf("fill called %d times concurrently, expected far fewer", calls.Load())
	}
}
