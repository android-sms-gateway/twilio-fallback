package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

type Encryptor interface {
	Encrypt(text string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

type encryptor struct {
	key []byte
}

// NewEncryptor returns a new Encryptor instance using the given key.
// The key must be at least 32 bytes, if it is not, an error is returned.
func NewEncryptor(key []byte) (Encryptor, error) {
	if len(key) < 32 {
		return nil, errors.New("key must be at least 32 bytes")
	}
	return &encryptor{key: key[:32]}, nil
}

// Encrypt encrypts the given text using AES GCM. The output is base64
// encoded and includes the nonce and ciphertext.
func (e *encryptor) Encrypt(text string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(text), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts the given ciphertext using AES GCM. The input is
// base64 encoded and should include the nonce and ciphertext. If the
// input is malformed, an error is returned.
func (e *encryptor) Decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("malformed ciphertext")
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
