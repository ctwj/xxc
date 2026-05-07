// Package yaml implements YAML encoding and decoding as defined in YAML 1.1.
// This is a minimal stub to satisfy test dependencies.
package yaml

import (
	"encoding/json"
)

// Marshal serializes the value provided into a YAML document.
func Marshal(in interface{}) ([]byte, error) {
	return json.Marshal(in)
}

// Unmarshal decodes the first document found within the in byte slice
// and assigns decoded values into the out value.
func Unmarshal(in []byte, out interface{}) error {
	return json.Unmarshal(in, out)
}
