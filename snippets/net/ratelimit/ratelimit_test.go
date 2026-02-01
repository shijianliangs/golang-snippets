package ratelimit

import (
	"context"
	"testing"
	"time"
)

func TestLimiter_BurstThenRefill(t *testing.T) {
	l := New(2, 50*time.Millisecond, 3) // 2 tokens per 50ms, burst 3
	defer l.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// First 3 should be immediate (burst).
	start := time.Now()
	for i := 0; i < 3; i++ {
		if err := l.Acquire(ctx); err != nil {
			t.Fatal(err)
		}
	}
	if d := time.Since(start); d > 10*time.Millisecond {
		t.Fatalf("burst acquisitions too slow: %v", d)
	}

	// Next token should wait until next refill tick.
	start = time.Now()
	if err := l.Acquire(ctx); err != nil {
		t.Fatal(err)
	}
	if time.Since(start) < 40*time.Millisecond {
		t.Fatalf("expected to block for refill")
	}
}

func TestLimiter_ContextCancel(t *testing.T) {
	l := New(1, time.Second, 1)
	defer l.Stop()

	// Consume burst.
	_ = l.Acquire(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	if err := l.Acquire(ctx); err == nil {
		t.Fatalf("expected ctx error")
	}
}
