package main

import (
	"bytes"
	"testing"

	"github.com/aid297/aid/v2/secret/symmetric/aes"
)

func TestECB(t *testing.T) {
	var (
		testKey   = []byte("1234567890abcdef")
		testPlain = []byte("hello, AES encrypt test!")
	)

	aesHelper, err := aes.New(aes.KeyBytes(testKey))
	if err != nil {
		t.Fatalf("创建 ECB 对象失败：%v", err)
	}

	cipherText, err := aesHelper.EncryptECB(testPlain)
	if err != nil {
		t.Fatalf("EncryptECB failed: %v", err)
	}

	plain, err := aesHelper.DecryptECB(cipherText)
	if err != nil {
		t.Fatalf("DecryptECB failed: %v", err)
	}

	if !bytes.Equal(plain, testPlain) {
		t.Fatalf("ECB mismatch: got %s, want %s", plain, testPlain)
	}
}

func TestECBBase64(t *testing.T) {
	var (
		testKey   = []byte("1234567890abcdef")
		testPlain = []byte("hello, AES encrypt test!")
	)

	aesHelper, err := aes.New(aes.KeyBytes(testKey))
	if err != nil {
		t.Fatalf("创建 ECB 对象失败：%v", err)
	}

	b64, err := aesHelper.EncryptECBBase64(testPlain)
	if err != nil {
		t.Fatalf("EncryptECBBase64 failed: %v", err)
	}

	plain, err := aesHelper.DecryptECBBase64(b64)
	if err != nil {
		t.Fatalf("DecryptECBBase64 failed: %v", err)
	}

	if !bytes.Equal(plain, testPlain) {
		t.Fatalf("ECB base64 mismatch: got %s, want %s", plain, testPlain)
	}
}

func TestCBC(t *testing.T) {
	var (
		testKey   = []byte("1234567890abcdef")
		testIV    = []byte("abcdef1234567890")
		testPlain = []byte("hello, AES encrypt test!")
	)

	aesHelper, err := aes.New(aes.KeyBytes(testKey), aes.IVBytes(testIV))
	if err != nil {
		t.Fatalf("创建 CBC 对象失败：%v", err)
	}

	cipherText, err := aesHelper.EncryptCBC(testPlain)
	if err != nil {
		t.Fatalf("EncryptCBC failed: %v", err)
	}

	plain, err := aesHelper.DecryptCBC(cipherText)
	if err != nil {
		t.Fatalf("DecryptCBC failed: %v", err)
	}

	if !bytes.Equal(plain, testPlain) {
		t.Fatalf("CBC mismatch: got %s, want %s", plain, testPlain)
	}
}

func TestCBCBase64(t *testing.T) {
	var (
		testKey   = []byte("1234567890abcdef")
		testIV    = []byte("abcdef1234567890")
		testPlain = []byte("hello, AES encrypt test!")
	)

	aesHelper, err := aes.New(aes.KeyBytes(testKey), aes.IVBytes(testIV))
	if err != nil {
		t.Fatalf("创建 CBC 对象失败：%v", err)
	}

	b64, err := aesHelper.EncryptCBCBase64(testPlain)
	if err != nil {
		t.Fatalf("EncryptCBCBase64 failed: %v", err)
	}

	plain, err := aesHelper.DecryptCBCBase64(b64)
	if err != nil {
		t.Fatalf("DecryptCBCBase64 failed: %v", err)
	}

	if !bytes.Equal(plain, testPlain) {
		t.Fatalf("CBC base64 mismatch: got %s, want %s", plain, testPlain)
	}
}

func TestInvalidKey(t *testing.T) {
	var testPlain = []byte("hello, AES encrypt test!")

	aesHelper, err := aes.New(aes.KeyString("shortkey"))
	if err != nil {
		t.Fatalf("创建对象失败：%v", err)
	}

	_, err = aesHelper.EncryptECB(testPlain)
	if err == nil {
		t.Fatal("expected error for invalid key length")
	}
}

func TestCBC192And256(t *testing.T) {
	var (
		testKey24 = []byte("1234567890abcdefghijklmn")
		testKey32 = []byte("1234567890abcdefghijklmnopqrstuv")
		testIV    = []byte("abcdef1234567890")
		testPlain = []byte("hello, AES encrypt test!")
		tests     = []struct {
			name string
			key  []byte
		}{
			{name: "AES-192", key: testKey24},
			{name: "AES-256", key: testKey32},
		}
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aesHelper, err := aes.New(aes.KeyBytes(tt.key), aes.IVBytes(testIV))
			if err != nil {
				t.Fatalf("创建 CBC 对象失败：%v", err)
			}
			cipherText, err := aesHelper.EncryptCBC(testPlain)
			if err != nil {
				t.Fatalf("EncryptCBC failed: %v", err)
			}
			plain, err := aesHelper.DecryptCBC(cipherText)
			if err != nil {
				t.Fatalf("DecryptCBC failed: %v", err)
			}
			if !bytes.Equal(plain, testPlain) {
				t.Fatalf("CBC mismatch: got %s, want %s", plain, testPlain)
			}
		})
	}
}

func TestRandKeyWithBits(t *testing.T) {
	var (
		key       []byte
		testPlain = []byte("hello, AES encrypt test!")
	)

	aesHelper, err := aes.New(aes.RandKeyWithBits(aes.AESKey256, &key), aes.RandIV())
	if err != nil {
		t.Fatalf("创建对象失败：%v", err)
	}

	if len(key) != 32 {
		t.Fatalf("随机 key 长度错误: got %d, want 32", len(key))
	}

	if _, err = aesHelper.EncryptCBC(testPlain); err != nil {
		t.Fatalf("EncryptCBC failed: %v", err)
	}
}

func TestRandKeyWithKeySize(t *testing.T) {
	var (
		key       []byte
		testPlain = []byte("hello, AES encrypt test!")
	)

	aesHelper, err := aes.New(aes.KeySize(aes.AESKey192), aes.RandKey(&key), aes.RandIV())
	if err != nil {
		t.Fatalf("创建对象失败：%v", err)
	}

	if len(key) != 24 {
		t.Fatalf("KeySize+RandKey 长度错误: got %d, want 24", len(key))
	}

	if _, err = aesHelper.EncryptCBC(testPlain); err != nil {
		t.Fatalf("EncryptCBC failed: %v", err)
	}
}

func TestRandKeyWithBitsInvalid(t *testing.T) {
	_, err := aes.New(aes.RandKeyWithBits(111), aes.RandIV())
	if err == nil {
		t.Fatal("expected error for invalid key bits")
	}
}
