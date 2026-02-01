# atomicfile

Atomic file writes (write temp file → fsync → rename).

Useful for config files, small state files, caches, etc.

## Usage

```go
err := atomicfile.WriteFile("/path/to/app.conf", []byte("enabled=true\n"), 0o600)
```

## Streaming usage

```go
err := atomicfile.WriteFileFunc("/path/to/data.json", 0o644, func(w io.Writer) error {
    _, err := io.WriteString(w, "...")
    return err
})
```

## Example

```bash
go test ./snippets/io/atomicfile
```

## Notes

- Uses `os.Rename`, which is atomic on POSIX when source and destination are on the same filesystem.
- On Windows, `os.Rename` may fail if the destination exists.
