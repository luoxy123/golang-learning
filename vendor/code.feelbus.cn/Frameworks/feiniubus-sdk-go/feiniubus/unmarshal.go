package feiniubus

import (
	"encoding/json"
	"io"
)

// Unmarshaler interface
type Unmarshaler interface {
	Unmarshal(r io.Reader, v interface{}) error
}

// JSONUnmarshaler is
type JSONUnmarshaler struct {
}

// Unmarshal returns
func (u *JSONUnmarshaler) Unmarshal(r io.Reader, v interface{}) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(v)
}
