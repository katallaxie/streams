package codec

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNilDecoder(t *testing.T) {
	v, err := NilDecoder.Decode([]byte("foo"))
	require.NoError(t, err)
	assert.Nil(t, v)
}

func TestStringDecoder(t *testing.T) {
	v, err := StringDecoder.Decode([]byte("foo"))
	require.NoError(t, err)
	assert.Equal(t, "foo", v)
}

func TestByteDecoder(t *testing.T) {
	v, err := ByteDecoder.Decode([]byte("foo"))
	require.NoError(t, err)
	assert.Equal(t, []byte("foo"), v)
}

func TestInt64Decoder(t *testing.T) {
	v, err := Int64Decoder.Decode([]byte("42"))
	require.NoError(t, err)
	assert.Equal(t, int64(42), v)
}

func TestStringEncoder(t *testing.T) {
	v, err := StringEncoder.Encode("foo")
	require.NoError(t, err)
	assert.Equal(t, []byte("foo"), v)
}

func TestByteEncoder(t *testing.T) {
	v, err := ByteEncoder.Encode([]byte("foo"))
	require.NoError(t, err)
	assert.Equal(t, []byte("foo"), v)
}

func TestInt64Encoder(t *testing.T) {
	v, err := Int64Encoder.Encode(int64(42))
	require.NoError(t, err)
	assert.Equal(t, []byte("42"), v)
}

func TestByteEncrypter(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	assert.NotNil(t, privateKey)

	v, err := ByteEncrypter.Encrypt([]byte("foo"), &privateKey.PublicKey)
	require.NoError(t, err)
	assert.NotEqual(t, []byte("foo"), v)

	v, err = ByteDecrypter.Decrypt(v, privateKey)
	require.NoError(t, err)
	assert.Equal(t, []byte("foo"), v)
}

func TestStringEncypter(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	assert.NotNil(t, privateKey)

	v, err := StringEncrypter.Encrypt("foo", &privateKey.PublicKey)
	require.NoError(t, err)
	assert.NotEqual(t, []byte("foo"), v)

	s, err := StringDecrypter.Decrypt(v, privateKey)
	require.NoError(t, err)
	assert.Equal(t, "foo", s)
}
