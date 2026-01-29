package httpclient

import (
	"context"
	"errors"
	"io"
	"math/rand"
	"net"
	"net/http"
	"time"
)

type Options struct {
	Timeout     time.Duration
	MaxRetries  int
	BaseBackoff time.Duration
	// RetryStatuses: if empty, defaults to 429 and 5xx.
	RetryStatuses map[int]bool
}

type Client struct {
	hc   *http.Client
	opt  Options
	rand *rand.Rand
}

func New(opt Options) *Client {
	if opt.Timeout == 0 {
		opt.Timeout = 10 * time.Second
	}
	if opt.BaseBackoff == 0 {
		opt.BaseBackoff = 200 * time.Millisecond
	}
	if opt.RetryStatuses == nil {
		opt.RetryStatuses = map[int]bool{429: true}
		for s := 500; s <= 599; s++ {
			opt.RetryStatuses[s] = true
		}
	}
	return &Client{
		hc: &http.Client{Timeout: opt.Timeout},
		opt: opt,
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	var lastErr error
	for attempt := 0; attempt <= c.opt.MaxRetries; attempt++ {
		resp, err := c.hc.Do(req)
		if err == nil && resp != nil && !c.shouldRetryStatus(resp.StatusCode) {
			return resp, nil
		}

		// If we got a response that we plan to retry, drain body to reuse TCP conn.
		if resp != nil && resp.Body != nil {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}

		if err != nil {
			lastErr = err
			if !isTransientNetErr(err) {
				break
			}
		} else {
			lastErr = errors.New("retryable http status")
		}

		if attempt == c.opt.MaxRetries {
			break
		}

		if err := sleepWithBackoff(req.Context(), c.opt.BaseBackoff, attempt, c.rand); err != nil {
			return nil, err
		}
	}
	return nil, lastErr
}

func (c *Client) shouldRetryStatus(code int) bool {
	return c.opt.RetryStatuses[code]
}

func sleepWithBackoff(ctx context.Context, base time.Duration, attempt int, r *rand.Rand) error {
	// expo backoff: base * 2^attempt, add jitter [0, base)
	sleep := base * time.Duration(1<<attempt)
	jitter := time.Duration(r.Int63n(int64(base)))
	wait := sleep + jitter
	select {
	case <-time.After(wait):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func isTransientNetErr(err error) bool {
	var ne net.Error
	if errors.As(err, &ne) {
		return ne.Timeout() || ne.Temporary()
	}
	// common transient: connection reset, EOF
	return true
}
