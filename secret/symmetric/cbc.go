package symmetric

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/aid297/aid/str"
)

type Cbc struct{}

var CbcApp Cbc

func (Cbc) padPKCS7(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func (Cbc) unPadPKCS7(src []byte, blockSize int) ([]byte, error) {
	length := len(src)
	if blockSize <= 0 {
		return nil, fmt.Errorf("invalid blockSize: %d", blockSize)
	}

	if length%blockSize != 0 || length == 0 {
		return nil, errors.New("invalid data len")
	}

	unpadding := int(src[length-1])
	if unpadding > blockSize || unpadding == 0 {
		return nil, errors.New("invalid unpadding")
	}

	padding := src[length-unpadding:]
	for i := 0; i < unpadding; i++ {
		if padding[i] != byte(unpadding) {
			return nil, errors.New("invalid padding")
		}
	}

	return src[:(length - unpadding)], nil
}

func (Cbc) Encrypt(plainText, key, iv []byte, ivs ...[]byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	plainText = Cbc{}.padPKCS7(plainText, blockSize)
	ivValue := ([]byte)(nil)
	if len(ivs) > 0 {
		ivValue = ivs[0]
	} else {
		ivValue = iv
	}
	blockMode := cipher.NewCBCEncrypter(block, ivValue)
	cipherText := make([]byte, len(plainText))
	blockMode.CryptBlocks(cipherText, plainText)

	return cipherText, nil
}

func (Cbc) Decrypt(cipherText, key, iv []byte, ivs ...[]byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(cipherText) < blockSize {
		return nil, errors.New("cipherText too short")
	}
	ivValue := ([]byte)(nil)
	if len(ivs) > 0 {
		ivValue = ivs[0]
	} else {
		ivValue = iv
	}
	if len(cipherText)%blockSize != 0 {
		return nil, errors.New("cipherText is not a multiple of the block size")
	}
	blockModel := cipher.NewCBCDecrypter(block, ivValue)
	plainText := make([]byte, len(cipherText))
	blockModel.CryptBlocks(plainText, cipherText)
	plainText, e := Cbc{}.unPadPKCS7(plainText, blockSize)
	if e != nil {
		return nil, e
	}
	return plainText, nil
}

func (Cbc) Demo() {
	key := "tjp5OPIU1ETF5s33fsLWdA=="
	iv := "0987654321098765"

	encrypted, err := Cbc{}.Encrypt([]byte("abcdefghijklmnopqrstuvwxyz"), []byte(key), []byte(iv))
	if err != nil {
		str.TerminalLogApp.New("[CBC] encrypt: %v").Error(err)
	}

	base64Encoded := base64.StdEncoding.EncodeToString(encrypted)
	str.TerminalLogApp.New("[CBC] base64 encoded: %s").Success(base64Encoded)

	base64Decoded, base64DecodeErr := base64.StdEncoding.DecodeString(base64Encoded)
	if base64DecodeErr != nil {
		str.TerminalLogApp.New("[CBC] base64 decode %v").Error(base64DecodeErr)
	}

	decryptCBC, err := Cbc{}.Decrypt(base64Decoded, []byte(key), []byte(iv))
	if err != nil {
		str.TerminalLogApp.New("[CBC] decrypt: %v").Error(err)
	}

	str.TerminalLogApp.New("[CBC] decrypted: %s").Success(string(decryptCBC))
}
