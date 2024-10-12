package utils

import "bytes"

// Custom unmarshaler that rejects extra fields
func StrictUnmarshal(data []byte, structure interface{}) error {
	// Unmarshal into the actual struct
	decoder := JSON.NewDecoder(bytes.NewBuffer(data))
	decoder.DisallowUnknownFields() // Reject unknown fields
	return decoder.Decode(structure)
}
