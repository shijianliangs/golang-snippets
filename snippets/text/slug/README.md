# slug

A tiny, dependency-free ASCII slug generator.

## Usage

```go
s := slug.Make("Hello, World!")
// s == "hello-world"
```

## Example

```bash
go test ./snippets/text/slug
```

## Notes

- Keeps only ASCII `[a-z0-9]`.
- Collapses separators/punctuation into single `-`.
- Drops non-ASCII characters (no transliteration).
