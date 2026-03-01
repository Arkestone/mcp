package syncer

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	s := New(time.Second, func() {})
	if s == nil {
		t.Fatal("expected non-nil Syncer")
	}
	if s.interval != time.Second {
		t.Fatalf("expected interval %v, got %v", time.Second, s.interval)
	}
	if s.fn == nil {
		t.Fatal("expected non-nil fn")
	}
}

func TestStartCallsFnImmediately(t *testing.T) {
	var count atomic.Int64
	s := New(time.Hour, func() { count.Add(1) })
	s.Start(context.Background())
	defer s.Stop()

	if got := count.Load(); got != 1 {
		t.Fatalf("expected fn called once immediately, got %d", got)
	}
}

func TestStartRunsFnPeriodically(t *testing.T) {
	var count atomic.Int64
	s := New(50*time.Millisecond, func() { count.Add(1) })
	s.Start(context.Background())

	time.Sleep(200 * time.Millisecond)
	s.Stop()

	if got := count.Load(); got <= 1 {
		t.Fatalf("expected fn called more than once, got %d", got)
	}
}

func TestStopHaltsPeriodicLoop(t *testing.T) {
	var count atomic.Int64
	s := New(50*time.Millisecond, func() { count.Add(1) })
	s.Start(context.Background())

	time.Sleep(150 * time.Millisecond)
	s.Stop()

	after := count.Load()
	time.Sleep(150 * time.Millisecond)

	if got := count.Load(); got != after {
		t.Fatalf("expected counter to stop at %d after Stop, got %d", after, got)
	}
}

func TestForceSync(t *testing.T) {
	var count atomic.Int64
	s := New(time.Hour, func() { count.Add(1) })
	s.Start(context.Background())
	defer s.Stop()

	before := count.Load()
	s.ForceSync()

	if got := count.Load(); got != before+1 {
		t.Fatalf("expected counter %d after ForceSync, got %d", before+1, got)
	}
}

func TestStopIsIdempotent(t *testing.T) {
	var count atomic.Int64
	s := New(50*time.Millisecond, func() { count.Add(1) })
	s.Start(context.Background())

	s.Stop()
	s.Stop() // should not panic
}

func TestContextCancellationStopsLoop(t *testing.T) {
	var count atomic.Int64
	ctx, cancel := context.WithCancel(context.Background())
	s := New(50*time.Millisecond, func() { count.Add(1) })
	s.Start(ctx)

	time.Sleep(150 * time.Millisecond)
	cancel()
	// Wait for goroutine to observe cancellation.
	s.wg.Wait()

	after := count.Load()
	time.Sleep(150 * time.Millisecond)

	if got := count.Load(); got != after {
		t.Fatalf("expected counter to stop at %d after context cancel, got %d", after, got)
	}
}

func TestConcurrentForceSync(t *testing.T) {
	var count atomic.Int64
	s := New(time.Hour, func() { count.Add(1) })
	s.Start(context.Background())
	defer s.Stop()

	const goroutines = 10
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for range goroutines {
		go func() {
			defer wg.Done()
			s.ForceSync()
		}()
	}
	wg.Wait()

	// 1 from Start + 10 from ForceSync
	if got := count.Load(); got != 1+goroutines {
		t.Fatalf("expected %d calls, got %d", 1+goroutines, got)
	}
}

// ---------------------------------------------------------------------------
// Additional nominal / error / limit tests
// ---------------------------------------------------------------------------

func TestStopBeforeStart(t *testing.T) {
	s := New(time.Hour, func() {})
	// Stop without Start should not panic
	s.Stop()
}

func TestVeryShortInterval(t *testing.T) {
	var count atomic.Int64
	s := New(10*time.Millisecond, func() { count.Add(1) })
	s.Start(context.Background())

	// Wait for several ticks
	time.Sleep(100 * time.Millisecond)
	s.Stop()

	if got := count.Load(); got < 3 {
		t.Errorf("expected at least 3 calls with 10ms interval over 100ms, got %d", got)
	}
}

func TestForceSyncBlocksUntilComplete(t *testing.T) {
	var running atomic.Int64
	s := New(time.Hour, func() {
		running.Add(1)
		time.Sleep(20 * time.Millisecond)
		running.Add(-1)
	})
	s.Start(context.Background())
	defer s.Stop()

	s.ForceSync()
	// After ForceSync returns, fn should no longer be running
	if running.Load() != 0 {
		t.Error("ForceSync should wait for fn to complete")
	}
}
