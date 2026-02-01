# ratelimit

A tiny, stdlib-only in-process rate limiter (burst + periodic refill).

## Usage

```go
l := ratelimit.New(10, time.Second, 20) // 10 tokens/sec, burst 20

defer l.Stop()

if err := l.Acquire(ctx); err != nil {
    // ctx deadline/cancel
}
```

## Example

```bash
go test ./snippets/net/ratelimit
```
