package jsonx

import "testing"

type cfg struct {
	A string `json:"a"`
	B int    `json:"b"`
}

func TestDecodeStrict_OK(t *testing.T) {
	var c cfg
	if err := DecodeStrict([]byte(`{"a":"x","b":2}`), &c); err != nil {
		t.Fatal(err)
	}
	if c.A != "x" || c.B != 2 {
		t.Fatalf("got %+v", c)
	}
}

func TestDecodeStrict_UnknownField(t *testing.T) {
	var c cfg
	if err := DecodeStrict([]byte(`{"a":"x","b":2,"c":3}`), &c); err == nil {
		t.Fatalf("expected error")
	}
}

func TestDecodeStrict_Trailing(t *testing.T) {
	var c cfg
	if err := DecodeStrict([]byte(`{"a":"x","b":2} trailing`), &c); err == nil {
		t.Fatalf("expected error")
	}
}
