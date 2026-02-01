package jwtverify

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"testing"
)

func TestVerify_EdDSA(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatal(err)
	}

	header := map[string]any{"alg": "EdDSA", "typ": "JWT"}
	payload := map[string]any{"sub": "123", "admin": true}

	p1 := b64url(mustJSON(t, header))
	p2 := b64url(mustJSON(t, payload))
	input := p1 + "." + p2
	sig := ed25519.Sign(priv, []byte(input))
	tok := input + "." + b64url(sig)

	h, p, err := Verify(tok, pub)
	if err != nil {
		t.Fatalf("Verify error: %v", err)
	}
	if h["alg"] != "EdDSA" {
		t.Fatalf("header alg = %v", h["alg"])
	}
	if p["sub"] != "123" {
		t.Fatalf("payload sub = %v", p["sub"])
	}
}

func TestVerify_BadSignature(t *testing.T) {
	pub, _, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatal(err)
	}

	tok := "e30.e30." + b64url([]byte("not-a-valid-sig"))
	_, _, err = Verify(tok, pub)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func b64url(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}
