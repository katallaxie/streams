package codec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNilDecoder(t *testing.T) {
	v, err := NilDecoder.Decode([]byte("foo"))
	assert.NoError(t, err)
	assert.Nil(t, v)
}

func TestStringDecoder(t *testing.T) {
	v, err := StringDecoder.Decode([]byte("foo"))
	assert.NoError(t, err)
	assert.Equal(t, "foo", v)
}

func TestByteDecoder(t *testing.T) {
	v, err := ByteDecoder.Decode([]byte("foo"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("foo"), v)
}

func TestInt64Decoder(t *testing.T) {
	v, err := Int64Decoder.Decode([]byte("42"))
	assert.NoError(t, err)
	assert.Equal(t, int64(42), v)
}

func TestStringEncoder(t *testing.T) {
	v, err := StringEncoder.Encode("foo")
	assert.NoError(t, err)
	assert.Equal(t, []byte("foo"), v)
}

func TestByteEncoder(t *testing.T) {
	v, err := ByteEncoder.Encode([]byte("foo"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("foo"), v)
}

func TestInt64Encoder(t *testing.T) {
	v, err := Int64Encoder.Encode(int64(42))
	assert.NoError(t, err)
	assert.Equal(t, []byte("42"), v)
}
