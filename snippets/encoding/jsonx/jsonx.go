package jsonx

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// DecodeStrict decodes JSON from data into v.
//
// It rejects:
//   - unknown fields (DisallowUnknownFields)
//   - trailing non-whitespace after the first JSON value
//
// This is handy for config files and API payload validation.
func DecodeStrict(data []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()

	if err := dec.Decode(v); err != nil {
		return err
	}
	// Ensure there's no second JSON value.
	if dec.More() {
		return fmt.Errorf("jsonx: multiple JSON values")
	}
	// Decoder doesn't expose a direct "only whitespace remains" check, but Decode
	// on an empty interface will return io.EOF only if no non-space remains.
	var extra any
	if err := dec.Decode(&extra); err == nil {
		return fmt.Errorf("jsonx: trailing data")
	}
	return nil
}
