package sm4

import (
	"bytes"
	"testing"
)

var (
	testKey   = []byte("1234567890abcdef")
	testIV    = []byte("abcdef1234567890")
	testPlain = []byte("hello, SM4 encrypt test!")
)

func TestECB(t *testing.T) {
	sm4Helper, err := New(KeyBytes(testKey))
	if err != nil {
		t.Fatalf("创建 ECB 对象失败：%v", err)
	}
	cipherText, err := sm4Helper.EncryptECB(testPlain)
	if err != nil {
		t.Fatalf("EncryptECB failed: %v", err)
	}
	plain, err := sm4Helper.DecryptECB(cipherText)
	if err != nil {
		t.Fatalf("DecryptECB failed: %v", err)
	}
	if !bytes.Equal(plain, testPlain) {
		t.Fatalf("ECB mismatch: got %s, want %s", plain, testPlain)
	}
	t.Logf("ECB OK: %s", plain)
}

func TestECBBase64(t *testing.T) {
	sm4Helper, err := New(KeyBytes(testKey))
	if err != nil {
		t.Fatalf("创建 ECB 对象失败：%v", err)
	}
	b64, err := sm4Helper.EncryptECBBase64(testPlain)
	if err != nil {
		t.Fatalf("EncryptECBBase64 failed: %v", err)
	}
	t.Logf("ECB base64: %s", b64)
	plain, err := sm4Helper.DecryptECBBase64(b64)
	if err != nil {
		t.Fatalf("DecryptECBBase64 failed: %v", err)
	}
	if !bytes.Equal(plain, testPlain) {
		t.Fatalf("ECB base64 mismatch: got %s, want %s", plain, testPlain)
	}
	t.Logf("ECB base64 OK: %s", plain)
}

func TestCBC(t *testing.T) {
	sm4Helper, err := New(KeyBytes(testKey), IVBytes(testIV))
	if err != nil {
		t.Fatalf("创建 CBC 对象失败：%v", err)
	}
	cipherText, err := sm4Helper.EncryptCBC(testPlain)
	if err != nil {
		t.Fatalf("EncryptCBC failed: %v", err)
	}
	plain, err := sm4Helper.DecryptCBC(cipherText)
	if err != nil {
		t.Fatalf("DecryptCBC failed: %v", err)
	}
	if !bytes.Equal(plain, testPlain) {
		t.Fatalf("CBC mismatch: got %s, want %s", plain, testPlain)
	}
	t.Logf("CBC OK: %s", plain)
}

func TestCBCBase64(t *testing.T) {
	sm4Helper, err := New(KeyBytes(testKey), IVBytes(testIV))
	if err != nil {
		t.Fatalf("创建 CBC 对象失败：%v", err)
	}
	b64, err := sm4Helper.EncryptCBCBase64(testPlain)
	if err != nil {
		t.Fatalf("EncryptCBCBase64 failed: %v", err)
	}
	t.Logf("CBC base64: %s", b64)
	plain, err := sm4Helper.DecryptCBCBase64(b64)
	if err != nil {
		t.Fatalf("DecryptCBCBase64 failed: %v", err)
	}
	if !bytes.Equal(plain, testPlain) {
		t.Fatalf("CBC base64 mismatch: got %s, want %s", plain, testPlain)
	}
	t.Logf("CBC base64 OK: %s", plain)
}

func TestInvalidKey(t *testing.T) {
	sm4Helper, err := New(KeyString("shortkey"))
	if err != nil {
		t.Fatalf("创建对象失败：%v", err)
	}
	_, err = sm4Helper.EncryptECB(testPlain)
	if err == nil {
		t.Fatal("expected error for invalid key length")
	}
	t.Logf("Invalid key error: %v", err)
}
