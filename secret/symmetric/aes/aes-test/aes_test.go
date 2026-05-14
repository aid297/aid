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

	aesHelper, err := aes.New(aes.KeyBytes(testKey), aes.AlgorithmECB())
	if err != nil {
		t.Fatalf("创建 ECB 对象失败：%v", err)
	}

	cipherText, err := aesHelper.Encrypt(testPlain)
	if err != nil {
		t.Fatalf("加密失败：%v", err)
	}

	plain, err := aesHelper.Decrypt(cipherText)
	if err != nil {
		t.Fatalf("解密失败：%v", err)
	}

	if !bytes.Equal(plain, testPlain) {
		t.Fatalf("比对失败：结果 %s，期望 %s", plain, testPlain)
	}
}

func TestECBBase64(t *testing.T) {
	var (
		testKey   = []byte("1234567890abcdef")
		testPlain = []byte("hello, AES encrypt test!")
	)

	aesHelper, err := aes.New(aes.KeyBytes(testKey), aes.AlgorithmECB())
	if err != nil {
		t.Fatalf("创建 ECB 对象失败：%v", err)
	}

	b64, err := aesHelper.EncryptBase64(testPlain)
	if err != nil {
		t.Fatalf("加密失败：%v", err)
	}

	plain, err := aesHelper.DecryptBase64(b64)
	if err != nil {
		t.Fatalf("解密失败：%v", err)
	}

	if !bytes.Equal(plain, testPlain) {
		t.Fatalf("比对失败：结果 %s，期望 %s", plain, testPlain)
	}
}

func TestCBC(t *testing.T) {
	var (
		testKey   = []byte("1234567890abcdef")
		testIV    = []byte("abcdef1234567890")
		testPlain = []byte("hello, AES encrypt test!")
	)

	aesHelper, err := aes.New(aes.KeyBytes(testKey), aes.IVBytes(testIV), aes.AlgorithmCBC())
	if err != nil {
		t.Fatalf("创建 CBC 对象失败：%v", err)
	}

	cipherText, err := aesHelper.Encrypt(testPlain)
	if err != nil {
		t.Fatalf("加密失败：%v", err)
	}

	plain, err := aesHelper.Decrypt(cipherText)
	if err != nil {
		t.Fatalf("解密失败：%v", err)
	}

	if !bytes.Equal(plain, testPlain) {
		t.Fatalf("比对失败：结果 %s，期望 %s", plain, testPlain)
	}
}

func TestCBCBase64(t *testing.T) {
	var (
		testKey   = []byte("1234567890abcdef")
		testIV    = []byte("abcdef1234567890")
		testPlain = []byte("hello, AES encrypt test!")
	)

	aesHelper, err := aes.New(aes.KeyBytes(testKey), aes.IVBytes(testIV), aes.AlgorithmCBC())
	if err != nil {
		t.Fatalf("创建 CBC 对象失败：%v", err)
	}

	b64, err := aesHelper.EncryptBase64(testPlain)
	if err != nil {
		t.Fatalf("加密失败：%v", err)
	}

	plain, err := aesHelper.DecryptBase64(b64)
	if err != nil {
		t.Fatalf("解密失败：%v", err)
	}

	if !bytes.Equal(plain, testPlain) {
		t.Fatalf("比对失败：结果 %s，期望 %s", plain, testPlain)
	}
}

func TestInvalidKey(t *testing.T) {
	var testPlain = []byte("hello, AES encrypt test!")

	aesHelper, err := aes.New(aes.KeyString("shortkey"), aes.AlgorithmECB())
	if err != nil {
		t.Fatalf("创建对象失败：%v", err)
	}

	_, err = aesHelper.Encrypt(testPlain)
	if err == nil {
		t.Fatal("验证key长度错误")
	}
	t.Logf("验证不通过：%v", err)
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
			aesHelper, err := aes.New(aes.KeyBytes(tt.key), aes.IVBytes(testIV), aes.AlgorithmCBC())
			if err != nil {
				t.Fatalf("创建 CBC 对象失败：%v", err)
			}
			cipherText, err := aesHelper.Encrypt(testPlain)
			if err != nil {
				t.Fatalf("加密失败：%v", err)
			}
			plain, err := aesHelper.Decrypt(cipherText)
			if err != nil {
				t.Fatalf("解密失败：%v", err)
			}
			if !bytes.Equal(plain, testPlain) {
				t.Fatalf("比对失败：结果 %s，期望 %s", plain, testPlain)
			}
		})
	}
}

func TestRandKeyWithBits(t *testing.T) {
	var (
		key       []byte
		testPlain = []byte("hello, AES encrypt test!")
	)

	aesHelper, err := aes.New(aes.RandKeyWithBits(aes.AES256, &key), aes.RandIV())
	if err != nil {
		t.Fatalf("创建对象失败：%v", err)
	}

	if len(key) != 32 {
		t.Fatalf("随机 key 长度错误: got %d, want 32", len(key))
	}

	if _, err = aesHelper.Encrypt(testPlain); err != nil {
		t.Fatalf("加密失败：%v", err)
	}
}

func TestRandKeyWithKeySize(t *testing.T) {
	var (
		key       []byte
		testPlain = []byte("hello, AES encrypt test!")
	)

	aesHelper, err := aes.New(aes.KeySize(aes.AES192), aes.RandKey(&key), aes.RandIV(), aes.AlgorithmCBC())
	if err != nil {
		t.Fatalf("创建对象失败：%v", err)
	}

	if len(key) != 24 {
		t.Fatalf("KeySize+RandKey 长度错误: got %d, want 24", len(key))
	}

	if _, err = aesHelper.Encrypt(testPlain); err != nil {
		t.Fatalf("加密失败：%v", err)
	}
}

func TestRandKeyWithBitsInvalid(t *testing.T) {
	_, err := aes.New(aes.RandKeyWithBits(111), aes.RandIV())
	if err == nil {
		t.Fatal("验证key位数错误")
	}
}
