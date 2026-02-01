# jwtverify

Verify a JWS-signed JWT (compact serialization) using only the Go standard library.

This snippet focuses on **signature verification**. It deliberately does **not** validate claims like `exp`, `nbf`, `aud`, etc.

## Supported algorithms

- `RS256` / `RS384` / `RS512` (RSA PKCS#1 v1.5)
- `PS256` / `PS384` / `PS512` (RSA-PSS)
- `ES256` / `ES384` / `ES512` (ECDSA, raw `r||s` signature per JWS)
- `EdDSA` (Ed25519)

## Usage

```go
header, payload, err := jwtverify.Verify(token, publicKey)
if err != nil {
    // invalid token/signature
}
_ = header
_ = payload
```

## Example

Run:

```bash
go test ./snippets/crypto/jwtverify -run TestVerify_EdDSA
```

## Notes

- The JWS ECDSA signature format is **raw** `r||s` (fixed-size) and not ASN.1 DER.
- Rejects `alg=none`.

## References

- RFC 7515: JSON Web Signature (JWS)
- RFC 7519: JSON Web Token (JWT)
