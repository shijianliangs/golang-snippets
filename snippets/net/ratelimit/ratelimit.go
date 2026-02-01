package ratelimit

import (
	"context"
	"time"
)

// Limiter is a small, dependency-free rate limiter.
//
// It implements a token-bucket-like interface using time.Ticker.
// The limiter allows up to burst immediate acquisitions, then refills at rate.
//
// Typical usage:
//
// 	l := ratelimit.New(10, time.Second, 20) // 10/sec with burst 20
// 	if err := l.Acquire(ctx); err != nil { ... }
// 	defer l.Stop()
//
// Note: this is best suited for in-process client-side throttling.
type Limiter struct {
	ch     chan struct{}
	ticker *time.Ticker
	done   chan struct{}
}

// New creates a limiter that refills rate tokens every per duration.
//
// For example, New(5, time.Second, 10) refills 5 tokens each second and allows
// bursts up to 10.
func New(rate int, per time.Duration, burst int) *Limiter {
	if rate <= 0 {
		rate = 1
	}
	if per <= 0 {
		per = time.Second
	}
	if burst <= 0 {
		burst = 1
	}

	l := &Limiter{
		ch:     make(chan struct{}, burst),
		ticker: time.NewTicker(per),
		done:   make(chan struct{}),
	}

	// Pre-fill burst.
	for i := 0; i < burst; i++ {
		l.ch <- struct{}{}
	}

	go func() {
		defer func() {
			// Drain ticker on stop.
			l.ticker.Stop()
		}()
		for {
			select {
			case <-l.done:
				return
			case <-l.ticker.C:
				// Refill up to rate tokens, without exceeding burst.
				for i := 0; i < rate; i++ {
					select {
					case l.ch <- struct{}{}:
					default:
						// full
						i = rate
					}
				}
			}
		}
	}()

	return l
}

// Acquire blocks until a token is available or ctx is done.
func (l *Limiter) Acquire(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-l.ch:
		return nil
	}
}

// Stop stops the internal refill goroutine.
func (l *Limiter) Stop() {
	select {
	case <-l.done:
		return
	default:
		close(l.done)
	}
}
