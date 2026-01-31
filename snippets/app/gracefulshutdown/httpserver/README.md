# net/http graceful shutdown (with OS signals)

A minimal pattern for running an `http.Server` with:
- `context.Context` cancellation
- OS signal handling (`SIGINT`, `SIGTERM`)
- graceful shutdown with timeout
- proper error handling (distinguish `http.ErrServerClosed`)

## Usage

```go
srv := &http.Server{Addr: ":8080", Handler: mux}

ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()

if err := gracefulshutdown.RunHTTPServer(ctx, srv, 10*time.Second); err != nil {
  log.Fatalf("server error: %v", err)
}
```

## Why

Graceful shutdown prevents:
- dropping in-flight requests
- leaving keep-alive connections hanging
- abrupt termination during deployments

## Notes

- Use a shutdown timeout that matches your infrastructure (load balancer drain time, etc.).
- This snippet is stdlib-only.
