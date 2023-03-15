package codec

import (
	"strconv"
)

// NilDecoder is a decoder that always returns a nil, no matter the input.
var NilDecoder Decoder[any] = func([]byte) (any, error) { return nil, nil }

// Decoder is an adapter allowing to use a function as a decoder.
type Decoder[T any] func(value []byte) (T, error)

// Decode transforms byte data to the desired type.
func (f Decoder[T]) Decode(value []byte) (T, error) {
	return f(value)
}

// StringDecoder ...
var StringDecoder Decoder[string] = func(b []byte) (string, error) {
	return string(b), nil
}

// ByteDecoder represents a byte decoder.
var ByteDecoder Decoder[[]byte] = func(b []byte) ([]byte, error) {
	return b, nil
}

// Int64Decoder represents an int64 decoder.
var Int64Decoder Decoder[int64] = func(b []byte) (int64, error) {
	return strconv.ParseInt(string(b), 10, 64)
}

// Encoder is an adapter allowing to use a function as an encoder.
type Encoder[T any] func(T) ([]byte, error)

// Encode transforms the typed data to bytes.
func (f Encoder[T]) Encode(value T) ([]byte, error) {
	return f(value)
}

// StringEncoder is a string encoder.
var StringEncoder Encoder[string] = func(v string) ([]byte, error) {
	return []byte(v), nil
}

// ByteEncoder represents a byte encoder.
var ByteEncoder Encoder[[]byte] = func(v []byte) ([]byte, error) {
	return v, nil
}

// Int64Encoder represents an int64 encoder.
var Int64Encoder Encoder[int64] = func(v int64) ([]byte, error) {
	return []byte(strconv.FormatInt(v, 10)), nil
}
