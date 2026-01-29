# HTTP client (retry + timeout + backoff)

A small, reusable `net/http` client wrapper for Go:
- request-level timeout
- retry on transient errors (network errors, 5xx, 429)
- exponential backoff with jitter
- context support

## Usage

```go
c := httpclient.New(httpclient.Options{
  Timeout: 5 * time.Second,
  MaxRetries: 3,
  BaseBackoff: 200 * time.Millisecond,
})

req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.com", nil)
resp, err := c.Do(req)
```

## Notes
- Retries are only safe for idempotent requests by default. If you retry POST, ensure your API is idempotent.
