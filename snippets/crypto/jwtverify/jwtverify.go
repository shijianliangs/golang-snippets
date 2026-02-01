package jwtverify

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

// Verify verifies a JWS-signed JWT (three-segment compact serialization) using a
// key.
//
// Supported algorithms:
//   - RS256/RS384/RS512 (RSA PKCS#1 v1.5)
//   - PS256/PS384/PS512 (RSA-PSS)
//   - ES256/ES384/ES512 (ECDSA)
//   - EdDSA (Ed25519)
//
// It returns the decoded header and payload JSON objects if the signature is valid.
//
// This snippet intentionally does NOT validate claims (exp/nbf/aud/etc.).
func Verify(token string, key any) (header map[string]any, payload map[string]any, err error) {
	p1, rest, ok := strings.Cut(token, ".")
	if !ok {
		return nil, nil, errors.New("jwt: expected 3 segments")
	}
	p2, p3, ok := strings.Cut(rest, ".")
	if !ok {
		return nil, nil, errors.New("jwt: expected 3 segments")
	}

	rawHdr, err := b64urlDecode(p1)
	if err != nil {
		return nil, nil, fmt.Errorf("jwt: decode header: %w", err)
	}
	rawPayload, err := b64urlDecode(p2)
	if err != nil {
		return nil, nil, fmt.Errorf("jwt: decode payload: %w", err)
	}
	sig, err := b64urlDecode(p3)
	if err != nil {
		return nil, nil, fmt.Errorf("jwt: decode signature: %w", err)
	}

	if err := json.Unmarshal(rawHdr, &header); err != nil {
		return nil, nil, fmt.Errorf("jwt: parse header json: %w", err)
	}
	if err := json.Unmarshal(rawPayload, &payload); err != nil {
		return nil, nil, fmt.Errorf("jwt: parse payload json: %w", err)
	}

	alg, _ := header["alg"].(string)
	if alg == "" {
		return nil, nil, errors.New("jwt: missing alg")
	}
	if alg == "none" {
		return nil, nil, errors.New("jwt: alg none not allowed")
	}

	signingInput := p1 + "." + p2
	if err := verifyJWS(alg, signingInput, sig, key); err != nil {
		return nil, nil, err
	}
	return header, payload, nil
}

func verifyJWS(alg, signingInput string, sig []byte, key any) error {
	switch alg {
	case "RS256", "RS384", "RS512":
		pub, ok := key.(*rsa.PublicKey)
		if !ok {
			return errors.New("jwt: RS* requires *rsa.PublicKey")
		}
		h, err := hashForAlg(alg)
		if err != nil {
			return err
		}
		return rsa.VerifyPKCS1v15(pub, h, digest(h, []byte(signingInput)), sig)
	case "PS256", "PS384", "PS512":
		pub, ok := key.(*rsa.PublicKey)
		if !ok {
			return errors.New("jwt: PS* requires *rsa.PublicKey")
		}
		h, err := hashForAlg(alg)
		if err != nil {
			return err
		}
		opts := &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash, Hash: h}
		return rsa.VerifyPSS(pub, h, digest(h, []byte(signingInput)), sig, opts)
	case "ES256", "ES384", "ES512":
		pub, ok := key.(*ecdsa.PublicKey)
		if !ok {
			return errors.New("jwt: ES* requires *ecdsa.PublicKey")
		}
		h, err := hashForAlg(alg)
		if err != nil {
			return err
		}
		// JWS uses raw (r||s) with fixed size based on curve.
		sz := (pub.Curve.Params().BitSize + 7) / 8
		if len(sig) != 2*sz {
			return errors.New("jwt: invalid ECDSA signature size")
		}
		r := new(big.Int).SetBytes(sig[:sz])
		s := new(big.Int).SetBytes(sig[sz:])
		if !ecdsa.Verify(pub, digest(h, []byte(signingInput)), r, s) {
			return errors.New("jwt: invalid signature")
		}
		return nil
	case "EdDSA":
		pub, ok := key.(ed25519.PublicKey)
		if !ok {
			return errors.New("jwt: EdDSA requires ed25519.PublicKey")
		}
		if !ed25519.Verify(pub, []byte(signingInput), sig) {
			return errors.New("jwt: invalid signature")
		}
		return nil
	default:
		return fmt.Errorf("jwt: unsupported alg %q", alg)
	}
}

func hashForAlg(alg string) (crypto.Hash, error) {
	switch alg {
	case "RS256", "PS256", "ES256":
		return crypto.SHA256, nil
	case "RS384", "PS384", "ES384":
		return crypto.SHA384, nil
	case "RS512", "PS512", "ES512":
		return crypto.SHA512, nil
	default:
		return 0, fmt.Errorf("jwt: no hash for alg %q", alg)
	}
}

func digest(h crypto.Hash, msg []byte) []byte {
	switch h {
	case crypto.SHA256:
		sum := sha256.Sum256(msg)
		return sum[:]
	case crypto.SHA384:
		sum := sha512.Sum384(msg)
		return sum[:]
	case crypto.SHA512:
		sum := sha512.Sum512(msg)
		return sum[:]
	default:
		panic("unsupported hash")
	}
}

func b64urlDecode(s string) ([]byte, error) {
	dec := base64.RawURLEncoding
	b, err := dec.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return b, nil
}
