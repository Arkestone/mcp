// Package syncer provides a reusable background sync goroutine.
// It runs a sync function immediately on start, then periodically.
package syncer

import (
	"context"
	"sync"
	"time"
)

// Syncer runs a function on a periodic interval in the background.
type Syncer struct {
	interval time.Duration
	fn       func()
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

// New creates a Syncer that calls fn every interval.
func New(interval time.Duration, fn func()) *Syncer {
	return &Syncer{interval: interval, fn: fn}
}

// Start runs fn immediately, then launches a background goroutine
// that calls fn every interval until ctx is canceled or Stop is called.
func (s *Syncer) Start(ctx context.Context) {
	ctx, s.cancel = context.WithCancel(ctx)
	s.fn()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.fn()
			}
		}
	}()
}

// Stop cancels the background goroutine and waits for it to exit.
func (s *Syncer) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
	s.wg.Wait()
}

// ForceSync runs the sync function immediately (blocks until complete).
func (s *Syncer) ForceSync() {
	s.fn()
}
