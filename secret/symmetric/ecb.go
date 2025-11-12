package symmetric

import (
	"bytes"
	"crypto/aes"
	"errors"
	"fmt"
)

type Ecb struct{}

var EcbApp Ecb

// padPKCS7 pads the plaintext to be a multiple of the block size
func (Ecb) padPKCS7(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(plaintext, padtext...)
}

// unPadPKCS7 removes the padding from the decrypted text
func (Ecb) unPadPKCS7(plaintext []byte) []byte {
	length := len(plaintext)
	unpadding := int(plaintext[length-1])

	return plaintext[:(length - unpadding)]
}

func (Ecb) unPadPKCS72(src []byte, blockSize int) ([]byte, error) {
	length := len(src)
	if blockSize <= 0 {
		return nil, fmt.Errorf("invalid blockSize: %d", blockSize)
	}

	if length%blockSize != 0 || length == 0 {
		return nil, errors.New("invalid data len")
	}

	unpadding := int(src[length-1])
	if unpadding > blockSize {
		return nil, fmt.Errorf("invalid unpadding: %d", unpadding)
	}

	if unpadding == 0 {
		return nil, errors.New("invalid unpadding: 0")
	}

	padding := src[length-unpadding:]
	for i := 0; i < unpadding; i++ {
		if padding[i] != byte(unpadding) {
			return nil, errors.New("invalid padding")
		}
	}

	return src[:(length - unpadding)], nil
}

// Encrypt encrypts plaintext using AES in ECB mode
func (Ecb) Encrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	plaintext = Ecb{}.padPKCS7(plaintext, blockSize)
	cipherText := make([]byte, len(plaintext))

	for start := 0; start < len(plaintext); start += blockSize {
		block.Encrypt(cipherText[start:start+blockSize], plaintext[start:start+blockSize])
	}

	return cipherText, nil
}

// Decrypt decrypts cipherText using AES in ECB mode
func (Ecb) Decrypt(key, cipherText []byte) ([]byte, error) {
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

	return Ecb{}.unPadPKCS72(plaintext, blockSize)
}
