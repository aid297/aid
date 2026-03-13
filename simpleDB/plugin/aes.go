package plugin

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

func init() {
	RegisterEncryptor("aes", func(key string) Encryptor {
		// Use SHA-256 to derive a 32-byte key from the input string
		hash := sha256.Sum256([]byte(key))
		return &AESEncryptor{key: hash[:]}
	})
}

type AESEncryptor struct {
	key []byte
}

func (a *AESEncryptor) Encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Encrypt and authenticate the data
	// The nonce is prepended to the ciphertext
	return gcm.Seal(nonce, nonce, data, nil), nil
}

func (a *AESEncryptor) Decrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(data) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
