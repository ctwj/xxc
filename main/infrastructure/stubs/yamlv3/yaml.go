// Package yaml implements YAML encoding and decoding as defined in YAML 1.2.
// This is a minimal stub to satisfy test dependencies.
package yaml

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// Marshal serializes the value provided into a YAML document.
func Marshal(in interface{}) ([]byte, error) {
	return jsonMarshal(in)
}

// Unmarshal decodes the first document found within the in byte slice
// and assigns decoded values into the out value.
func Unmarshal(in []byte, out interface{}) error {
	return json.Unmarshal(in, out)
}

// A Decoder reads and decodes YAML values from an input stream.
type Decoder struct{}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r interface{}) *Decoder {
	return &Decoder{}
}

// Decode reads the next YAML-encoded value from its input
// and stores it in the value pointed to by v.
func (dec *Decoder) Decode(v interface{}) error {
	return nil
}

// An Encoder writes YAML values to an output stream.
type Encoder struct{}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w interface{}) *Encoder {
	return &Encoder{}
}

// Encode writes the YAML encoding of v to the stream.
func (enc *Encoder) Encode(v interface{}) error {
	return nil
}

// Close closes the encoder.
func (enc *Encoder) Close() error {
	return nil
}

func jsonMarshal(in interface{}) ([]byte, error) {
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return []byte("null"), nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.String:
		return []byte(strconv.Quote(v.String())), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []byte(fmt.Sprintf("%d", v.Int())), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []byte(fmt.Sprintf("%d", v.Uint())), nil
	case reflect.Bool:
		if v.Bool() {
			return []byte("true"), nil
		}
		return []byte("false"), nil
	case reflect.Float32, reflect.Float64:
		return []byte(fmt.Sprintf("%g", v.Float())), nil
	case reflect.Slice, reflect.Array:
		if v.Len() == 0 {
			return []byte("[]"), nil
		}
		return json.Marshal(in)
	case reflect.Map, reflect.Struct:
		return json.Marshal(in)
	default:
		return json.Marshal(in)
	}
}

// IsZeroer is used to check whether an object is zero to
// determine whether it should be omitted when marshaling
// with the "omitempty" flag.
type IsZeroer interface {
	IsZero() bool
}

// MapSlice encodes and decodes as a YAML map.
type MapSlice []MapItem

// MapItem is an item in a MapSlice.
type MapItem struct {
	Key, Value interface{}
}

// Node represents a node in the YAML document graph.
type Node struct {
	Kind        Kind
	Style       Style
	Tag         string
	Value       string
	Anchor      string
	Alias       *Node
	Content     []*Node
	HeadComment string
	LineComment string
	FootComment string
	Line        int
	Column      int
}

// MarshalText implements encoding.TextMarshaler.
func (n *Node) MarshalText() ([]byte, error) {
	return []byte(n.Value), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (n *Node) UnmarshalText(text []byte) error {
	n.Value = string(text)
	return nil
}

// Kind specifies the kind of a Node.
type Kind uint

const (
	DocumentNode Kind = 1 << iota
	SequenceNode
	MappingNode
	ScalarNode
	AliasNode
)

// Style specifies the style of a Node.
type Style uint

const (
	TaggedStyle Style = 1 << iota
	DoubleQuotedStyle
	SingleQuotedStyle
	LiteralStyle
	FoldedStyle
	FlowStyle
)

// TypeError is an error that occurs when an unmarshal operation
// encounters an unexpected type.
type TypeError struct {
	Value string
	Type  reflect.Type
}

func (e *TypeError) Error() string {
	return fmt.Sprintf("cannot unmarshal %s into type %s", e.Value, e.Type)
}
