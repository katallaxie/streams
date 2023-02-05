package streams

// NilDecoder is a decoder that always returns a nil, no matter the input.
var NilDecoder DecoderFunc[any] = func([]byte) (any, error) { return nil, nil }

// Decoder represents a Kafka data decoder.
type Decoder[T any] interface {
	// Decode transforms byte data to the desired type.
	Decode([]byte) (T, error)
}

// DecoderFunc is an adapter allowing to use a function as a decoder.
type DecoderFunc[T any] func(value []byte) (T, error)

// Decode transforms byte data to the desired type.
func (f DecoderFunc[T]) Decode(value []byte) (T, error) {
	return f(value)
}

// Encoder represents a Kafka data encoder.
type Encoder[T any] interface {
	// Encode transforms the typed data to bytes.
	Encode(T) ([]byte, error)
}

// EncoderFunc is an adapter allowing to use a function as an encoder.
type EncoderFunc[T any] func(T) ([]byte, error)

// Encode transforms the typed data to bytes.
func (f EncoderFunc[T]) Encode(value T) ([]byte, error) {
	return f(value)
}

// ByteDecoder represents a byte decoder.
type ByteDecoder struct{}

// Decode transforms byte data to the desired type.
func (d ByteDecoder) Decode(b []byte) (interface{}, error) {
	return b, nil
}

// ByteEncoder represents a byte encoder.
type ByteEncoder struct{}

// Encode transforms the typed data to bytes.
func (e ByteEncoder) Encode(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, nil
	}

	return v.([]byte), nil
}

// StringDecoder represents a string decoder.
type StringDecoder struct{}

// Decode transforms byte data to a string.
func (d StringDecoder) Decode(b []byte) (string, error) {
	return string(b), nil
}

// StringEncoder represents a string encoder.
type StringEncoder struct{}

// Encode transforms the string data to bytes.
func (e StringEncoder) Encode(v string) ([]byte, error) {
	return []byte(v), nil
}
