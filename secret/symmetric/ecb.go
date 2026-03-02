package symmetric

import (
	"bytes"
	"crypto/aes"
	"errors"
	"fmt"
)

type ECB struct{}

func NewECB() *ECB { return &ECB{} }

func (*ECB) New() *ECB { return NewECB() }

// padPKCS7 pads the plaintext to be a multiple of the block size
func (*ECB) padPKCS7(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(plaintext, padText...)
}

// unPadPKCS7 removes the padding from the decrypted text
func (*ECB) unPadPKCS7(plaintext []byte) []byte {
	length := len(plaintext)
	unPadding := int(plaintext[length-1])

	return plaintext[:(length - unPadding)]
}

func (*ECB) unPadPKCS72(src []byte, blockSize int) ([]byte, error) {
	length := len(src)
	if blockSize <= 0 {
		return nil, fmt.Errorf("invalid blockSize: %d", blockSize)
	}

	if length%blockSize != 0 || length == 0 {
		return nil, errors.New("invalid data len")
	}

	unPadding := int(src[length-1])
	if unPadding > blockSize {
		return nil, fmt.Errorf("invalid unPadding: %d", unPadding)
	}

	if unPadding == 0 {
		return nil, errors.New("invalid unPadding: 0")
	}

	padding := src[length-unPadding:]
	for i := 0; i < unPadding; i++ {
		if padding[i] != byte(unPadding) {
			return nil, errors.New("invalid padding")
		}
	}

	return src[:(length - unPadding)], nil
}

// Encrypt encrypts plaintext using AES in ECB mode
func (*ECB) Encrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	plaintext = NewECB().padPKCS7(plaintext, blockSize)
	cipherText := make([]byte, len(plaintext))

	for start := 0; start < len(plaintext); start += blockSize {
		block.Encrypt(cipherText[start:start+blockSize], plaintext[start:start+blockSize])
	}

	return cipherText, nil
}

// Decrypt decrypts cipherText using AES in ECB mode
func (*ECB) Decrypt(key, cipherText []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	if len(cipherText)%blockSize != 0 {
		return nil, fmt.Errorf("cipherText is not a multiple of the block size")
	}

	plaintext := make([]byte, len(cipherText))

	for start := 0; start < len(cipherText); start += blockSize {
		block.Decrypt(plaintext[start:start+blockSize], cipherText[start:start+blockSize])
	}

	return (&ECB{}).unPadPKCS72(plaintext, blockSize)
}
