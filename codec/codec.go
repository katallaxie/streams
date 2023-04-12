package codec

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"strconv"
)

// Encrypter is a function that encrypts a value.
type Encrypter[T any] func(T, *rsa.PublicKey) ([]byte, error)

// Decrypter is a function that decrypts a value.
type Decrypter[T any] func([]byte, *rsa.PrivateKey) (T, error)

// ByteEncrypter represents a byte encrypter.
var ByteEncrypter Encrypter[[]byte] = func(b []byte, key *rsa.PublicKey) ([]byte, error) {
	hash := sha256.New()
	ciperText, err := rsa.EncryptOAEP(hash, rand.Reader, key, b, nil)
	if err != nil {
		return nil, err
	}

	return ciperText, nil
}

// Encrypt is a function that encryptes a value.
func (f Encrypter[T]) Encrypt(value T, key *rsa.PublicKey) ([]byte, error) {
	return f(value, key)
}

// Decrypt is a function that decryptes a value.
func (f Decrypter[T]) Decrypt(value []byte, key *rsa.PrivateKey) (T, error) {
	return f(value, key)
}

// ByteDecrypter represents a byte decrypter.
var ByteDecrypter Decrypter[[]byte] = func(b []byte, key *rsa.PrivateKey) ([]byte, error) {
	hash := sha256.New()
	plainText, err := rsa.DecryptOAEP(hash, rand.Reader, key, b, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}

// StringDecrypter represents a string decrypter.
var StringDecrypter Decrypter[string] = func(value []byte, key *rsa.PrivateKey) (string, error) {
	hash := sha256.New()
	plainText, err := rsa.DecryptOAEP(hash, rand.Reader, key, value, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

// StringEncrypter represents a string encrypter.
var StringEncrypter Encrypter[string] = func(b string, key *rsa.PublicKey) ([]byte, error) {
	hash := sha256.New()
	ciperText, err := rsa.EncryptOAEP(hash, rand.Reader, key, []byte(b), nil)
	if err != nil {
		return nil, err
	}

	return ciperText, nil
}

// NilDecoder is a decoder that always returns a nil, no matter the input.
var NilDecoder Decoder[any] = func([]byte) (any, error) { return nil, nil }

// Decoder is an adapter allowing to use a function as a decoder.
type Decoder[T any] func(value []byte) (T, error)

// Decode transforms byte data to the desired type.
func (f Decoder[T]) Decode(value []byte) (T, error) {
	return f(value)
}

// StringDecoder represents a string decoder.
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

// ByteEncryptedDecoder represents a byte encrypted decoder.
func ByteEncryptedDecoder() Decoder[[]byte] {
	return func(b []byte) ([]byte, error) {
		return b, nil
	}
}

// EncryptedEncoder represents an encrypted encoder.
var EncryptedEncoder Encoder[string] = func(v string) ([]byte, error) {
	return []byte(v), nil
}
