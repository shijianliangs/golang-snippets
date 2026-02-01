# jsonx

Strict JSON decoding helpers (stdlib-only).

## DecodeStrict

- rejects unknown fields
- rejects trailing data after the JSON value

### Usage

```go
var c Config
if err := jsonx.DecodeStrict(b, &c); err != nil {
    // handle invalid JSON
}
```

### Example

```bash
go test ./snippets/encoding/jsonx
```
