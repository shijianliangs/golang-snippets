# Worker pool (generic, context-aware)

A small worker pool for running jobs concurrently with:
- fixed number of workers
- generic input/output type (`T`)
- `context.Context` cancellation
- error propagation

This is intentionally minimal and stdlib-only.

## Usage

```go
p := workerpool.New[string](4)

ctx := context.Background()
outputs, err := p.Run(ctx, []string{"a", "b", "c"}, func(ctx context.Context, s string) (string, error) {
  // do work
  return strings.ToUpper(s), nil
})

// outputs has same length and ordering as inputs.
```

## Notes

- If any job returns an error, the pool cancels remaining work.
- If you need different output type from input, either:
  - map to an intermediate type first, or
  - split into separate pools per output type.
